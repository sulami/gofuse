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

// Some error we want to return for testing that is just not nil.
type someError string

func (s someError) Error() string {
	return string(s)
}

// Faux action function that returns depending on the length of the byte
// array we pass.
// 0    - return success
// 1    - return failure
// else - do not return, trigger a timeout
func FauxAction(in *[]byte) (*[]byte, error) {
	if len(*in) == 0 {
		return nil, nil
	} else if len(*in) == 1 {
		var err someError = "Fail."
		return nil, err
	} else {
		for {
			// Trigger a timeout
		}
	}
}

// Faux log function that should probably do something.
func FauxLog(s string) {
}

// Test initialization of new fuses.
func TestNewFuse(t *testing.T) {
	f := NewFuse(FauxAction, FauxLog, time.Second, 3, 2 * time.Second, 5)

	if f == nil {
		t.Error("Allocation failure.")
	} else if !f.good {
		t.Error("Fuse is bad.")
	} else if f.action == nil {
		t.Error("Action is nil.")
	} else if f.log == nil {
		t.Error("Log is nil.")
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
	f := NewFuse(FauxAction, FauxLog, time.Second, 3, 2 * time.Second, 5)

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
	f := NewFuse(FauxAction, FauxLog, time.Second, 3, 2 * time.Second, 5)

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

