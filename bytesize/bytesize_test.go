// Public Domain (-) 2018-present, The Peerbase Authors.
// See the Peerbase UNLICENSE file for details.

package bytesize

import (
	"testing"
)

func TestInt(t *testing.T) {
	type elem struct {
		v      Value
		expect int
		err    bool
	}
	for _, elem := range []elem{
		{100 * KB, 102400, false},
		{2048 * KB, 2097152, false},
		{32 * MB, 33554432, false},
		{32 * GB, 34359738368, false},
		{32 * TB, 35184372088832, false},
		{32 * PB, 36028797018963968, false},
		{MB + 1234, 1049810, false},
		{15360 * PB, 0, false},
	} {
		v, err := elem.v.Int()
		if err != nil {
			if elem.err {
				continue
			}
			t.Errorf("unexpected error when getting int value from %v: %s", elem.v, err)
			continue
		}
		if elem.expect != v {
			t.Errorf("mismatching int value: expected %d, got %d", elem.expect, v)
		}
	}
}

func TestParse(t *testing.T) {
	type elem struct {
		v      string
		expect Value
		err    bool
	}
	for _, elem := range []elem{
		{"100KB", 100 * KB, false},
		{"2MB", 2048 * KB, false},
		{"2048KB", 2048 * KB, false},
		{"2048 KB", 0, true},
		{"2048 kilobytes", 0, true},
		{"2048", 2 * KB, false},
		{"32MB", 32 * MB, false},
		{"32GB", 32 * GB, false},
		{"32TB", 32 * TB, false},
		{"8096PB", 8096 * PB, false},
		{"8096000PB", 0, true},
		{"8096 MiB", 0, true},
		{"0x1fa0 MiB", 0, true},
	} {
		v, err := Parse(elem.v)
		if err != nil {
			if elem.err {
				continue
			}
			t.Errorf("unexpected error when parsing %q: %s", elem.v, err)
			continue
		}
		if elem.expect != v {
			t.Errorf("mismatching parsed value: expected %q, got %q", elem.expect, v)
		}
	}
}

func TestString(t *testing.T) {
	type elem struct {
		v      Value
		expect string
	}
	for _, elem := range []elem{
		{100 * KB, "100KB"},
		{2048 * KB, "2MB"},
		{32 * MB, "32MB"},
		{32 * GB, "32GB"},
		{32 * TB, "32TB"},
		{32 * PB, "32PB"},
		{MB + 1234, "1049810B"},
	} {
		repr := elem.v.String()
		if elem.expect != repr {
			t.Errorf("mismatching string representation: expected %q, got %q", elem.expect, repr)
		}
	}
}
