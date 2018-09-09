// Public Domain (-) 2018-present, The Peerbase Authors.
// See the Peerbase UNLICENSE file for details.

// +build !linux !windows

package process

func subreaper() bool {
	return false
}
