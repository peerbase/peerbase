// Public Domain (-) 2018-present, The Peerbase Authors.
// See the Peerbase UNLICENSE file for details.

package overflow

import (
	"testing"
)

func TestMulU64(t *testing.T) {
	type elem struct {
		a      uint64
		b      uint64
		expect bool
	}
	for _, elem := range []elem{
		{4294967291, 4294967271, true},
		{4294967291, 4294967321, false},
	} {
		v, ok := MulU64(elem.a, elem.b)
		if elem.expect != ok {
			t.Errorf(
				"unexpected result for MulU64(%d, %d): expected %v, got %v (%d)",
				elem.a, elem.b, elem.expect, ok, v)
		}
	}
}
