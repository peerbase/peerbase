// Public Domain (-) 2018-present, The Peerbase Authors.
// See the Peerbase UNLICENSE file for details.

package osexit_test

import (
	"os"
	"testing"

	"peerbase.net/go/mock/osexit"
)

var osExit = os.Exit

func TestOsExit(t *testing.T) {
	osExit = osexit.Set()
	osExit(2)
	if !osexit.Called() {
		t.Fatalf("Mock os.Exit was not called")
	}
	status := osexit.Status()
	if status != 2 {
		t.Fatalf("Mock os.Exit did not set the right status code: expected 2, got %d", status)
	}
	osExit(3)
	status = osexit.Status()
	if status != 2 {
		t.Fatalf("Mock os.Exit overrode the status set by a previous call: expected 2, got %d", status)
	}
	osexit.Reset()
	if osexit.Called() {
		t.Fatalf("The reset mock os.Exit claims to have been called")
	}
	status = osexit.Status()
	if status != 0 {
		t.Fatalf("The reset mock os.Exit returned a non-zero status code: expected 0, got %d", status)
	}
}
