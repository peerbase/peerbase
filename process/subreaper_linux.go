// Public Domain (-) 2018-present, The Peerbase Authors.
// See the Peerbase UNLICENSE file for details.

package process

import (
	"golang.org/x/sys/unix"
)

var subreaper = func() bool {
	unix.Prctl(unix.PR_SET_CHILD_SUBREAPER, uintptr(1), 0, 0, 0)
	return true
}
