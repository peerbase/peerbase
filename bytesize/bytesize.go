// Public Domain (-) 2018-present, The Peerbase Authors.
// See the Peerbase UNLICENSE file for details.

// Package bytesize provides support for dealing with byte size values.
package bytesize

import (
	"fmt"
	"strconv"
	"strings"

	"peerbase.net/go/overflow"
)

// Constants representing common multiples of byte sizes.
const (
	KB = 1024
	MB = 1024 * KB
	GB = 1024 * MB
	TB = 1024 * GB
	PB = 1024 * TB
)

const maxInt = Value(^uint(0) >> 1)

// Value represents a byte size value.
type Value uint64

// Int checks if the byte size value would overflow the platform int and, if
// not, returns the value as an int.
func (v Value) Int() (int, error) {
	if v > maxInt {
		return 0, fmt.Errorf("bytesize: value %d overflows platform int", v)
	}
	return int(v), nil
}

// MarshalYAML implements the YAML encoding interface.
func (v Value) MarshalYAML() (interface{}, error) {
	return v.String(), nil
}

func (v Value) String() string {
	switch {
	case v%PB == 0:
		return strconv.FormatUint(uint64(v)/PB, 10) + "PB"
	case v%TB == 0:
		return strconv.FormatUint(uint64(v)/TB, 10) + "TB"
	case v%GB == 0:
		return strconv.FormatUint(uint64(v)/GB, 10) + "GB"
	case v%MB == 0:
		return strconv.FormatUint(uint64(v)/MB, 10) + "MB"
	case v%KB == 0:
		return strconv.FormatUint(uint64(v)/KB, 10) + "KB"
	default:
		return strconv.FormatUint(uint64(v), 10) + "B"
	}
}

// UnmarshalYAML implements the YAML decoding interface.
func (v *Value) UnmarshalYAML(unmarshal func(interface{}) error) error {
	raw := ""
	err := unmarshal(&raw)
	if err != nil {
		return err
	}
	*v, err = Parse(raw)
	return err
}

// Parse tries to parse a byte size Value from the given string. A byte size
// string is a sequence of decimal numbers and a unit suffix, e.g. "20GB",
// "1024KB", "100M", etc. Valid units are "B", "K", "KB", "M", "MB", "G", "GB",
// "T", "TB", "P", and "PB".
func Parse(s string) (Value, error) {
	var (
		err  error
		ok   bool
		unit string
		v    uint64
	)
	for i := len(s) - 1; i >= 0; i-- {
		char := s[i]
		if char >= '0' && char <= '9' {
			v, err = strconv.ParseUint(s[:i+1], 10, 64)
			if err != nil {
				return 0, fmt.Errorf("bytesize: unable to parse the decimal part of %q: %s", s, err)
			}
			unit = strings.TrimSpace(s[i+1:])
			break
		}
	}
	switch strings.ToLower(unit) {
	case "", "b", "byte", "bytes":
		ok = true
	case "k", "kb", "kilobyte", "kilobytes":
		v, ok = overflow.MulU64(v, KB)
	case "m", "mb", "megabyte", "megabytes":
		v, ok = overflow.MulU64(v, MB)
	case "g", "gb", "gigabyte", "gigabytes":
		v, ok = overflow.MulU64(v, GB)
	case "t", "tb", "terabyte", "terabytes":
		v, ok = overflow.MulU64(v, TB)
	case "p", "pb", "petabyte", "petabytes":
		v, ok = overflow.MulU64(v, PB)
	default:
		return 0, fmt.Errorf("bytesize: unsupported unit %q specified in %q", unit, s)
	}
	if !ok {
		return 0, fmt.Errorf("bytesize: string value %q overflows uint64 when parsed", s)
	}
	return Value(v), nil
}
