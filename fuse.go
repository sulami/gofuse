package fuse

import "time"

type Fuse struct {
	// Is this fuse currently blown?
	good bool

	// Function to execute to try to get a positive response.
	action func(*[]byte) ([]byte, error)

	// Function to execute in case bad (or good) things happen for
	// logging purposes.
	log func(string)

	// Timout after which to call a query.
	requestTimeout time.Duration

	// How often requests are allowed to time out before the fuse
	// blows.
	requestTries uint

	// The interval in which we try to contact an offline fuse.
	recoveryInterval time.Duration

	// How often requests have to come back successfully before the
	// fuse gets enabled again.
	recoveryTries uint

	// How often we successfully contacted a blown fuse in a row.
	recoverySuccesses uint
}

// Call the supplied action to determine the current status. If it
// returns failure or times out, return false.
func (f *Fuse) try() bool {
	// TODO
	// do the action
	// start the timeout
	// wait for a result or the timeout
	// return the result
	return false
}

// Recovery has succeeded, bring the fuse back online.
func (f *Fuse) unblow() {
	// f.recoverySuccesses = 0
	// f.good = true
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
	// f.good = false
	// go f.recovery()
}

// TODO
// Everything that needs to be publicly visible, like
// new fuse
// set fuse options
// query fuse status
// delete fuse

