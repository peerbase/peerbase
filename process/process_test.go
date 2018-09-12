// Public Domain (-) 2018-present, The Peerbase Authors.
// See the Peerbase UNLICENSE file for details.

package process

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
	"testing"

	"peerbase.net/go/mock/osexit"
)

func TestCreatePIDFile(t *testing.T) {
	tmp := setup(t)
	defer cleanup(tmp)
	fpath := filepath.Join(tmp, "test.pid")
	err := CreatePIDFile(fpath)
	if err != nil {
		t.Fatalf("Unexpected error creating PID file: %s", err)
	}
	written, err := ioutil.ReadFile(fpath)
	if err != nil {
		t.Fatalf("Unexpected error reading PID file: %s", err)
	}
	expected := os.Getpid()
	pid, err := strconv.ParseInt(string(written), 10, 64)
	if err != nil {
		t.Fatalf("Unexpected error parsing PID file contents as an int: %s", err)
	}
	if expected != int(pid) {
		t.Fatalf("Mismatching PID file contents: expected %d, got %d", expected, int(pid))
	}
	osExit = osexit.Set()
	Exit(2)
	if !osexit.Called() || osexit.Status() != 2 {
		t.Fatalf("Exit call did not behave as expected")
	}
	_, err = os.Stat(fpath)
	if err == nil {
		t.Fatalf("Calling Exit did not remove the created PID file as expected")
	}
	if !os.IsNotExist(err) {
		t.Fatalf("Calling Exit did not remove the created PID file as expected, got error: %s", err)
	}
	fpath = filepath.Join(tmp+"-nonexistent-directory", "test.pid")
	err = CreatePIDFile(fpath)
	if err == nil {
		t.Fatalf("Expected an error when creating PID file in a non-existent directory")
	}
}

func TestDisableDefaultExit(t *testing.T) {
	testMode = true
	ResetHandlers()
	called := false
	SetExitHandler(func() {
		called = true
	})
	osExit = osexit.Set()
	resetExit()
	send(syscall.SIGTERM)
	if !osexit.Called() {
		t.Fatalf("os.Exit was not called on SIGTERM")
	}
	if !called {
		t.Fatalf("Exit handler not called on SIGTERM")
	}
	DisableDefaultExit()
	called = false
	osexit.Reset()
	resetExit()
	send(syscall.SIGTERM)
	if osexit.Called() {
		t.Fatalf("os.Exit was called on SIGTERM even when DisableDefaultExit had been called")
	}
	if !called {
		t.Fatalf("Exit handler not called on the second SIGTERM")
	}
}

func TestExit(t *testing.T) {
	testMode = true
	ResetHandlers()
	called := false
	SetExitHandler(func() {
		called = true
	})
	osExit = osexit.Set()
	resetExit()
	Exit(7)
	if !osexit.Called() {
		t.Fatalf("Exit did not call os.Exit")
	}
	status := osexit.Status()
	if status != 7 {
		t.Fatalf("Exit did not set the right status code: expected 7; got %d", status)
	}
	if !called {
		t.Fatalf("Exit handler was not called on calling Exit")
	}
	osexit.Reset()
	go func() {
		Exit(8)
	}()
	<-testSig
	wait <- struct{}{}
	if osexit.Called() {
		t.Fatalf("Second call to Exit called os.Exit")
	}
}

func TestInit(t *testing.T) {
	tmp := setup(t)
	defer cleanup(tmp)
	locked, err := Init(tmp, "peerbase")
	if err != nil {
		t.Fatalf("Unexpected error initialising process: %s", err)
	}
	if !locked {
		t.Fatalf("Init failed to acquire process lock")
	}
	locked, err = Init(tmp+"-nonexistent-directory", "peerbase")
	if locked {
		t.Fatalf("Successfully acquired a process lock in a non-existing directory")
	}
	if err == nil {
		t.Fatalf("Expected an error when calling Init in a non-existing directory")
	}
}
func TestLock(t *testing.T) {
	tmp := setup(t)
	defer cleanup(tmp)
	err := Lock(tmp, "peerbase")
	if err != nil {
		t.Fatalf("Unexpected error acquiring Lock: %s", err)
	}
	err = Lock(tmp, "peerbase")
	if err == nil {
		t.Fatalf("Expected an error when calling Lock on an already locked path")
	}
	fpath := filepath.Join(tmp, fmt.Sprintf("peerbase-%d.lock", os.Getpid()))
	_, err = os.Stat(fpath)
	if err != nil {
		t.Fatalf("Unexpected error accessing the raw lock file: %s", err)
	}
	osExit = osexit.Set()
	Exit(2)
	_, err = os.Stat(fpath)
	if err == nil {
		t.Fatalf("Calling Exit did not remove the lock file as expected")
	}
	if !os.IsNotExist(err) {
		t.Fatalf("Calling Exit did not remove the lock file as expected, got error: %s", err)
	}
	err = Lock(tmp+"-nonexistent-directory", "peerbase")
	if err == nil {
		t.Fatalf("Expected an error when calling Lock in a non-existing directory")
	}
}

func TestSignalHandler(t *testing.T) {
	testMode = true
	ResetHandlers()
	called := false
	SetSignalHandler(syscall.SIGHUP, func() {
		called = true
	})
	send(syscall.SIGABRT)
	if called {
		t.Fatalf("Signal handler erroneously called on SIGABRT")
	}
	send(syscall.SIGHUP)
	if !called {
		t.Fatalf("Signal handler not called on SIGHUP")
	}
}

func cleanup(tmp string) {
	os.RemoveAll(tmp)
}

func resetExit() {
	mu.Lock()
	exiting = false
	mu.Unlock()
}

func send(sig syscall.Signal) {
	syscall.Kill(syscall.Getpid(), sig)
	<-testSig
}

func setup(t *testing.T) string {
	testMode = true
	resetExit()
	ResetHandlers()
	tmp, err := ioutil.TempDir("", "peerbase-process")
	if err != nil {
		t.Skipf("Unable to create temporary directory for tests: %s", err)
	}
	return tmp
}
