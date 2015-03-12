package fuse

import (
	"io"
	"log"
	"time"
)

type Fuse struct {
	// Is this fuse currently blown?
	good bool

	// Function to execute to try to get a positive response.
	action func(*[]byte, chan []byte)

	// Logger to use when logging anything on our own, eg. blown
	// fusesg
	logger *log.Logger

	// Timout after which to call a query.
	requestTimeout time.Duration

	// How often requests are allowed to time out before the fuse
	// blows.
	requestTries uint

	// How often we already failed in a row.
	requestFails uint

	// The interval in which we try to contact an offline fuse.
	recoveryInterval time.Duration

	// How often requests have to come back successfully before the
	// fuse gets enabled again.
	recoveryTries uint

	// How often we successfully contacted a blown fuse in a row.
	recoverySuccesses uint
}

// Call the supplied action to determine the current status. Returns a
// non-nil error if it times out.
func (f *Fuse) try(in *[]byte, out chan []byte) error {
	retval := make(chan []byte)
	timeout := make(chan bool, 1)

	go func () {
		time.Sleep(f.requestTimeout)
		timeout <- false
	} ()

	go f.action(in, retval)

	select {
	case <-timeout:
		f.requestFails++
		return nil // TODO return a non-nil error
	case r := <-retval:
		out <- r
	}

	return nil
}

// Recovery has succeeded, bring the fuse back online.
func (f *Fuse) unblow() {
	f.recoverySuccesses = 0
	f.good = true
}

// Try to get query successes until we hit the threshold and unblow the
// fuse officially. To be run as goroutine.
func (f *Fuse) recovery() {
	// TODO
	// sleep f.recoveryInterval
	// try recovery
	// if success:
	//	increment recoverySuccesses
	//	if recoverySuccesses >= recoveryTries:
	//		unblow(?)
	// if fail:
	//	recoverySuccesses = 0
}

// Blow the fuse and initiate recovery.
func (f *Fuse) blow() {
	f.good = false
	// go f.recovery()
}

func (f *Fuse) log(msg string) {
	f.logger.Print(msg)
}

// Create and initialize a new fuse and return it.
func NewFuse(action func(*[]byte, chan []byte),
             logwriter io.Writer,
             requestTimeout time.Duration,
             requestTries uint,
             recoveryInterval time.Duration,
             recoveryTries uint) *Fuse {
	f := new(Fuse)
	if f == nil {
		return nil
	}

	f.good = true
	f.action = action
	f.logger = log.New(logwriter, "gofuse: ",
	                   log.Lshortfile|log.Lmicroseconds)
	f.requestTimeout = requestTimeout
	f.requestTries = requestTries
	f.recoveryInterval = recoveryInterval
	f.recoveryTries = recoveryTries

	return f
}

func (f *Fuse) Query(in *[]byte, out chan []byte) {
	err := f.try(in, out)
	if (err != nil) {
		f.log("Timeout triggered.")
		// TODO perhabs blow
	}
}

// TODO
// set fuse options
// query fuse status
// delete fuse

