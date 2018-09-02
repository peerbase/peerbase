// Public Domain (-) 2018-present, The Peerbase Authors.
// See the Peerbase UNLICENSE file for details.

package bytesize

import (
	"testing"

	"gopkg.in/yaml.v2"
)

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
		{"2048 KB", 2048 * KB, false},
		{"2048 kilobytes", 2048 * KB, false},
		{"2048", 2 * KB, false},
		{"8096 PB", 8096 * PB, false},
		{"8096000 PB", 8096 * PB, true},
		{"8096 MiB", 8096 * MB, true},
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
	} {
		repr := elem.v.String()
		if elem.expect != repr {
			t.Errorf("mismatching string representation: expected %q, got %q", elem.expect, repr)
		}
	}
}

func TestYAMLMarshal(t *testing.T) {
	type elem struct {
		v      Value
		expect string
	}
	for _, elem := range []elem{
		{100 * KB, "100KB"},
		{2048 * KB, "2MB"},
	} {
		out, err := yaml.Marshal(elem.v)
		if err != nil {
			t.Errorf("unexpected error when marshalling as YAML: %s", err)
			continue
		}
		if out[len(out)-1] != '\n' {
			t.Errorf("expected newline as the final byte of the marshalled YAML output: %q", out)
			continue
		}
		out = out[:len(out)-1]
		if elem.expect != string(out) {
			t.Errorf("mismatching YAML encoding: expected %q, got %q", elem.expect, out)
		}
	}
}
