package eon

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"unicode/utf8"
)

const hex = "0123456789abcdef"

var (
	encoders sync.Map
	mstates  sync.Pool
)

var (
	marshalerType = reflect.TypeOf((*Marshaler)(nil)).Elem()
)

type encoder func(*mstate, reflect.Value) error

type mstate struct {
	bytes.Buffer
	commented map[string]bool
	comments  map[string]string
	scratch   [64]byte
}

type stringEncoder struct {
	multiline bool
}

func (s *stringEncoder) encode(m *mstate, rv reflect.Value) error {
	v := rv.String()
	start := m.Len()
	m.WriteByte('"')
	from := 0
	for i := 0; i < len(v); {
		c := v[i]
		if c < utf8.RuneSelf {
			if isPrintable(c) {
				i++
				continue
			}
			if from < i {
				m.WriteString(v[from:i])
			}
			m.WriteByte('\\')
			switch c {
			case '"', '\\':
				m.WriteByte(c)
			case '\n':
				m.Truncate(start)
				encodeMultiline(m, v)
				return nil
			case '\r':
				m.WriteByte('r')
			case '\t':
				m.WriteByte('t')
			default:
				m.WriteByte('x')
				m.WriteByte(hex[c>>4])
				m.WriteByte(hex[c&0xf])
			}
			i++
			from = i
			continue
		}
		r, size := utf8.DecodeRuneInString(v[i:])
		if r == utf8.RuneError {
			if from < i {
				m.WriteString(v[from:i])
			}
			m.WriteString(`\ufffd`)
			i += size
			from = i
			continue
		}
		i += size
	}
	if from < len(v) {
		m.WriteString(v[from:])
	}
	m.WriteByte('"')
	return nil
}

func encodeBool(m *mstate, rv reflect.Value) error {
	if rv.Bool() {
		m.WriteString("true")
	} else {
		m.WriteString("false")
	}
	return nil
}

func encodeInt(m *mstate, rv reflect.Value) error {
	m.Write(strconv.AppendInt(m.scratch[:0], rv.Int(), 10))
	return nil
}

func encodeMultiline(m *mstate, v string) {
}

func encodeUint(m *mstate, rv reflect.Value) error {
	m.Write(strconv.AppendUint(m.scratch[:0], rv.Uint(), 10))
	return nil
}

func getEncoder(rt reflect.Type) (encoder, error) {
	if enc, ok := encoders.Load(rt); ok {
		return enc.(encoder), nil
	}
	var (
		enc encoder
		err error
		wg  sync.WaitGroup
	)
	wg.Add(1)
	// Add a temporary handler to deal with recursive types.
	actual, loaded := encoders.LoadOrStore(rt, func(m *mstate, rv reflect.Value) error {
		wg.Wait()
		if err != nil {
			return err
		}
		return enc(m, rv)
	})
	if loaded {
		return actual.(encoder), nil
	}
	enc, err = typeEncoder(rt)
	if err == nil {
		encoders.Store(rt, enc)
	}
	wg.Done()
	return enc, err
}

func isPrintable(c byte) bool {
	if c > 34 && c < 127 {
		return true
	}
	return c == 32 || c == 33
}

func marshal(v interface{}, comments map[string]string) ([]byte, error) {
	rv := reflect.ValueOf(v)
	if !rv.IsValid() {
		return nil, ErrNilInterfaceValue
	}
	enc, err := getEncoder(rv.Type())
	if err != nil {
		return nil, err
	}
	m := newMstate(comments)
	if err = enc(m, rv); err != nil {
		mstates.Put(m)
		return nil, err
	}
	out := make([]byte, m.Len())
	copy(out, m.Bytes())
	mstates.Put(m)
	return out, nil
}

func newMstate(comments map[string]string) *mstate {
	if v := mstates.Get(); v != nil {
		e := v.(*mstate)
		e.Reset()
		if comments == nil {
			e.commented = nil
			e.comments = nil
		} else {
			e.commented = map[string]bool{}
			e.comments = comments
		}
		return e
	}
	if comments == nil {
		return &mstate{}
	}
	return &mstate{
		commented: map[string]bool{},
		comments:  comments,
	}
}

func newStringEncoder(multiline bool) *stringEncoder {
	return &stringEncoder{
		multiline: multiline,
	}
}

func typeEncoder(rt reflect.Type) (encoder, error) {
	kind := rt.Kind()
	switch kind {
	case reflect.Bool:
		return encodeBool, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return encodeInt, nil
	case reflect.String:
		return newStringEncoder(false).encode, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return encodeUint, nil
	}
	return nil, fmt.Errorf("eon: could not create encoder for %s", rt)
}
