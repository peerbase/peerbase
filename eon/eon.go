// Public Domain (-) 2018-present, The Peerbase Authors.
// See the Peerbase UNLICENSE file for details.

// Package eon implements encoding and decoding of EON (Extensible Object
// Notation) data. The mapping between EON and Go values are documented in the
// description of the Marshal and Unmarshal functions.
package eon

import (
	"errors"
)

// Error values.
var (
	ErrFloatInf          = errors.New("eon: cannot encode Inf float value")
	ErrFloatNaN          = errors.New("eon: cannot encode NaN float value")
	ErrNilInterfaceValue = errors.New("eon: cannot encode nil interface value")
)

// Marshaler is the interface implemented by types that can marshal themselves
// into valid EON.
type Marshaler interface {
	MarshalEON(scratch []byte, opts EncodeOpts) ([]byte, error)
}

// Unmarshaler is the interface implemented by types that can unmarshal an EON
// description of themselves. UnmarshalEON must copy any of the data in the
// given byte slice if it wishes to retain the data after returning.
type Unmarshaler interface {
	UnmarshalEON([]byte) error
}

// Marshal returns the EON encoding of v.
//
// EON cannot represent cyclic data structures and Marshal does not handle them.
// Passing cyclic structures to Marshal will result in an infinite recursion.
func Marshal(v interface{}) ([]byte, error) {
	return marshal(v, nil)
}

// MarshalWithComments is like Marshal but includes the given comment headers.
func MarshalWithComments(v interface{}, comments map[string]string) ([]byte, error) {
	return marshal(v, comments)
}

// Unmarshal parses the EON-encoded data and stores the result in the value
// pointed to by v.
func Unmarshal(data []byte, v interface{}) error {
	return nil
}
