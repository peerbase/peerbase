// Public Domain (-) 2010-present, The Peerbase Authors.
// See the Peerbase UNLICENSE file for details.

// Package process provides utilities to manage the current system process.
package process

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
)

var (
	exitDisabled bool
	exiting      bool
	mu           sync.RWMutex // protects exitDisabled, exiting, registry
)

var (
	osExit   = os.Exit
	registry = make(map[os.Signal][]func())
	testMode = false
	testSig  = make(chan struct{}, 10)
	wait     = make(chan struct{})
)

type lockFile struct {
	file string
	link string
}

func (l *lockFile) release() {
	os.Remove(l.file)
	os.Remove(l.link)
}

// CreatePIDFile writes the current process ID to a new file at the given path.
// The written file is removed when the process exits on receiving an
// os.Interrupt or SIGTERM signal — either directly or through the Exit call.
func CreatePIDFile(path string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0660)
	if err != nil {
		return err
	}
	fmt.Fprintf(f, "%d", os.Getpid())
	err = f.Close()
	if err == nil {
		SetExitHandler(func() {
			os.Remove(path)
		})
	}
	return err
}

// DisableDefaultExit will prevent the process from automatically exiting after
// processing os.Interrupt or SIGTERM signals. The disabled status will not be
// enforced if Exit is called.
func DisableDefaultExit() {
	mu.Lock()
	exitDisabled = true
	mu.Unlock()
}

// Exit runs the registered exit handlers, as if the os.Interrupt signal had
// been sent, and then terminates the process with the given status code. Exit
// blocks until the process terminates if it has already been called elsewhere.
func Exit(code int) {
	mu.Lock()
	if exiting {
		mu.Unlock()
		if testMode {
			testSig <- struct{}{}
		}
		<-wait
		return
	}
	exiting = true
	handlers := registry[os.Interrupt]
	mu.Unlock()
	for _, handler := range handlers {
		handler()
	}
	osExit(code)
}

// Init acquires a process lock and writes the PID file for the current process.
// It returns false for the locked status if there was an error acquiring the
// process lock.
func Init(directory string, name string) (locked bool, err error) {
	err = Lock(directory, name)
	if err != nil {
		return false, err
	}
	return true, CreatePIDFile(filepath.Join(directory, name+".pid"))
}

// Lock tries to acquire a process lock in the given directory. The acquired
// lock file is released when the process exits on receiving an os.Interrupt or
// SIGTERM signal — either directly or through the Exit call.
//
// This function has only been tested for correctness on Unix systems with
// filesystems where link is atomic. It may not work as expected on NFS mounts
// or on platforms like Windows.
func Lock(directory string, name string) error {
	file := filepath.Join(directory, fmt.Sprintf("%s-%d.lock", name, os.Getpid()))
	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0660)
	if err != nil {
		return err
	}
	f.Close()
	link := filepath.Join(directory, name+".lock")
	err = os.Link(file, link)
	if err != nil {
		// We don't remove the lock file here so that calling Lock multiple
		// times from the same process doesn't remove an existing lock.
		return err
	}
	l := &lockFile{
		file: file,
		link: link,
	}
	SetExitHandler(l.release)
	return nil
}

// ResetHandlers drops all currently registered handlers.
func ResetHandlers() {
	mu.Lock()
	registry = map[os.Signal][]func(){}
	mu.Unlock()
}

// SetExitHandler registers the given handler function to run when receiving
// os.Interrupt or SIGTERM signals.
func SetExitHandler(handler func()) {
	mu.Lock()
	registry[os.Interrupt] = prepend(registry[os.Interrupt], handler)
	registry[syscall.SIGTERM] = prepend(registry[syscall.SIGTERM], handler)
	mu.Unlock()
}

// SetSignalHandler registers the given handler function to run when receiving
// the specified signal.
func SetSignalHandler(signal os.Signal, handler func()) {
	mu.Lock()
	registry[signal] = prepend(registry[signal], handler)
	mu.Unlock()
}

func handleSignals() {
	notifier := make(chan os.Signal, 100)
	signal.Notify(notifier)
	go func() {
		for sig := range notifier {
			mu.RLock()
			handlers, found := registry[sig]
			mu.RUnlock()
			if found {
				for _, handler := range handlers {
					handler()
				}
			}
			mu.RLock()
			if !exitDisabled {
				if sig == syscall.SIGTERM || sig == os.Interrupt {
					osExit(1)
				}
			}
			mu.RUnlock()
			if testMode {
				testSig <- struct{}{}
			}
		}
	}()
}

func prepend(xs []func(), handler func()) []func() {
	return append([]func(){handler}, xs...)
}

func init() {
	handleSignals()
}
