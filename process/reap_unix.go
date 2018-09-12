// Public Domain (-) 2018-present, The Peerbase Authors.
// See the Peerbase UNLICENSE file for details.

// +build !windows

package process

import (
	"os"
	"os/signal"
	"syscall"
)

var (
	osGetpid = os.Getpid
)

// ReapOrphans listens out for SIGCHLD and reaps orphaned processes.
//
// If the current process is not PID 1, then on certain platforms (currently
// only Linux), it registers the process as a subreaper so as to be able to reap
// orphans. On other platforms, the function exits immediately and does nothing.
func ReapOrphans() {
	if osGetpid() != 1 && !SetAsSubreaper() {
		return
	}
	sigs := make(chan struct{}, 100)
	go notifySIGCHLD(sigs)
	status := syscall.WaitStatus(0)
	for range sigs {
		for {
			_, err := syscall.Wait4(-1, &status, 0, nil)
			if err == syscall.ECHILD {
				break
			}
		}
		if testMode {
			testSig <- struct{}{}
			break
		}
	}
}

// SetAsSubreaper tries to set the current process as a subreaper on the
// platforms that support it. It returns a boolean indicating whether it was
// successful or not.
func SetAsSubreaper() bool {
	return subreaper()
}

// We use an intermediary channel to handle SIGCHLD so that we don't block
// whilst reaping.
func notifySIGCHLD(sigs chan struct{}) {
	notifier := make(chan os.Signal, 100)
	signal.Notify(notifier, syscall.SIGCHLD)
	for range notifier {
		select {
		case sigs <- struct{}{}:
		default:
		}
		if testMode {
			break
		}
	}
	if testMode {
		signal.Stop(notifier)
	}
}
