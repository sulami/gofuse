package fuse

import (
	"testing"
	"time"
)

// Test the general sanity of the testing environment.
func TestSanity(t *testing.T) {
	if (1 != 1) {
		t.Fail()
	}
}

// Faux action function that returns depending on the length of the byte
// array we pass.
// 0    - return success
// 1    - return failure
// else - do not return, trigger a timeout
func FauxAction(in *[]byte, out chan []byte) {
	if len(*in) == 0 {
		out <- []byte("S")
	} else if len(*in) == 1 {
		out <- []byte("F")
	} else {
		for {
			// Trigger a timeout
		}
	}
}

// Test initialization of new fuses.
func TestNewFuse(t *testing.T) {
	// TODO use log.Logger
	f := NewFuse(FauxAction, nil, 1, time.Second, 3, 2 * time.Second, 5)

	if f == nil {
		t.Error("Allocation failure.")
	} else if !f.good {
		t.Error("Fuse is bad.")
	} else if f.action == nil {
		t.Error("Action is nil.")
	} else if f.logger == nil {
		t.Error("Logger is nil.")
	} else if f.requestTimeout != time.Second {
		t.Error("requestTimeout mismatch.")
	} else if f.requestTries != 3 {
		t.Error("requestTries mismatch.")
	} else if f.recoveryInterval != 2 * time.Second {
		t.Error("recoveryInterval mismatch.")
	} else if f.recoveryTries != 5 {
		t.Error("recoveryTries mismatch.")
	}
}

// Test blowing of fuses.
func TestBlowFuse(t *testing.T) {
	f := NewFuse(FauxAction, nil, 1, time.Second, 3, 2 * time.Second, 5)

	if !f.good {
		t.Error("Fuse is already blown.")
	}

	f.blow()

	if f.good {
		t.Error("Fuse has not been blown.")
	}

	// TODO check for recovery status
}

// Test unblowing of fuses.
func TestUnblowFuse(t *testing.T) {
	f := NewFuse(FauxAction, nil, 1, time.Second, 3, 2 * time.Second, 5)

	f.blow()

	if f.good {
		t.Error("Fuse has not been blown.")
	}

	// Emulate we already had a few successes for testing the reset.
	f.recoverySuccesses = 2

	f.unblow()

	if !f.good {
		t.Error("Fuse has not been unblown.")
	} else if f.recoverySuccesses != 0 {
		t.Error("recoverySuccesses have not been reset.")
	}

	// TODO check for recovery status
}

func TestTimeout(t *testing.T) {
	f := NewFuse(FauxAction, nil, 1, time.Second / 5, 3, 3 * time.Second, 5)

	arg := []byte("TIMEOUT TIME")
	retval := make(chan []byte)
	timeout := make(chan bool)

	// Timeout BEFORE the Fuse-internal timeout should trigger.
	go func() {
		time.Sleep(time.Second / 10)
		timeout <- true
	}()
	go f.Query(&arg, retval)

	select {
	case <-retval:
		t.Error("Action returned data where it was not supposed to.")
	case <-f.timeout:
		t.Error("Fuse timeout triggered too early.")
	case <-timeout:
		// That's actually a good thing™
	}

	// Timeout AFTER the Fuse-internal timeout should trigger.
	go func() {
		time.Sleep(time.Second / 2)
		timeout <- true
	}()

	go f.Query(&arg, retval)

	select {
	case <-retval:
		t.Error("Action returned data where it was not supposed to.")
	case <-f.timeout:
		// That's actually a good thing™
	case <-timeout:
		t.Error("Action did not return data in time.")
	}
}

func TestBlowingInUse(t *testing.T) {
	// This fuse will timeout really quickly and blow after just 2
	// failures. We use a larger queue size because we do not want
	// to check the results.
	f := NewFuse(FauxAction, nil, 10, time.Second / 10, 2, time.Second, 5)

	arg := []byte("TIMEOUT TIME")
	retval := make(chan []byte)

	// Do not use goroutines to keep track of where we are.
	if f.requestFails != 0 {
		t.Error("Wrong baseline.")
	}
	if !f.good {
		t.Error("Fuse is already blown.")
	}

	f.Query(&arg, retval)

	if f.requestFails != 1 {
		t.Error("requestFails did not count up to one.")
	}
	if !f.good {
		t.Error("Fuse has blown too early.")
	}

	f.Query(&arg, retval)

	if f.requestFails != 2 {
		t.Error("requestFails did not count up to two.")
	}
	if f.good {
		t.Error("Fuse did not blow when it was supposed to.")
	}
}

func TestBlownFuseReturnsFast(t *testing.T) {
	// We make sure a blown fuse returns failure immediately and
	// does not wait for the timeout.
	f := NewFuse(FauxAction, nil, 1, time.Second, 2, time.Second, 5)
	f.good = false

	// Does not matter at all.
	arg := []byte("TIMEOUT TIME")
	retval := make(chan []byte)

	timeout := time.After(time.Second / 10)
	go f.Query(&arg, retval)

	select {
	case <-f.timeout:
		// Desired behaviour.
	case <-retval:
		// There should be nothing coming out of here.
		t.Error("Blown fuse sent data.")
	case <-timeout:
		t.Error("Blown fuse did not return immediately.")
	}
}

