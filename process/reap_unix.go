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
	notifier := make(chan os.Signal, 4096)
	signal.Notify(notifier, syscall.SIGCHLD)
	status := syscall.WaitStatus(0)
	for range notifier {
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
	if testMode {
		signal.Stop(notifier)
	}
}

// SetAsSubreaper tries to set the current process as a subreaper on the
// platforms that support it. It returns a boolean indicating whether it was
// successful or not.
func SetAsSubreaper() bool {
	return subreaper()
}
