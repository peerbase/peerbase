package eon

import (
	"bytes"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"sync"
	"time"
	"unicode/utf8"

	"peerbase.net/go/bytesize"
)

// Encoding options.
const (
	OptInline EncodeOpts = 1 << iota
	OptMultiline
	OptToplevel
)

const hex = "0123456789abcdef"

var (
	encoders sync.Map
	mstates  sync.Pool
)

var (
	bytesizeType  = reflect.TypeOf(bytesize.Value(0))
	marshalerType = reflect.TypeOf((*Marshaler)(nil)).Elem()
	durationType  = reflect.TypeOf(time.Duration(0))
)

// EncodeOpts defines the various options for a value encoder.
type EncodeOpts int

func (e EncodeOpts) inline() bool {
	return e&OptInline != 0
}

func (e EncodeOpts) multiline() bool {
	return e&OptMultiline != 0
}

type arrayEncoder struct {
	elem encoder
	size int
}

func (e *arrayEncoder) encode(m *mstate, rv reflect.Value, opts EncodeOpts) error {
	return nil
}

type encoder func(*mstate, reflect.Value, EncodeOpts) error

type fieldEncoder struct {
	enc  encoder
	idx  int
	name string
}

type mapEncoder struct {
	key   encoder
	value encoder
}

func (e *mapEncoder) encode(m *mstate, rv reflect.Value, opts EncodeOpts) error {
	return nil
}

type mstate struct {
	bytes.Buffer
	commented map[string]bool
	comments  map[string]string
	indent    int
	scratch   [64]byte
}

type sliceEncoder struct {
	elem   encoder
	inline bool
}

func (e *sliceEncoder) encode(m *mstate, rv reflect.Value, opts EncodeOpts) error {
	if e.inline {
		m.WriteByte('[')
	}
	n := rv.Len()
	for i := 0; i < n; i++ {
		if e.inline && i != 0 {
			m.WriteByte(' ')
		}
		e.elem(m, rv.Index(i), 0)
	}
	if e.inline {
		m.WriteByte(']')
	}
	return nil
}

type structEncoder struct {
	fields []*fieldEncoder
}

func (e *structEncoder) encode(m *mstate, rv reflect.Value, opts EncodeOpts) error {
	return nil
}

func encodeBool(m *mstate, rv reflect.Value, opts EncodeOpts) error {
	if rv.Bool() {
		m.WriteString("true")
	} else {
		m.WriteString("false")
	}
	return nil
}

func encodeByteSlice(m *mstate, rv reflect.Value, opts EncodeOpts) error {
	v := rv.String()
	start := m.Len()
	m.WriteByte('"')
	_ = start
	_ = v
	m.WriteByte('"')
	return nil
}

func encodeByteSize(m *mstate, rv reflect.Value, opts EncodeOpts) error {
	m.WriteString(bytesize.Value(rv.Uint()).String())
	return nil
}

func encodeDuration(m *mstate, rv reflect.Value, opts EncodeOpts) error {
	m.WriteString(time.Duration(rv.Int()).String())
	return nil
}

func encodeFloat(m *mstate, v float64, bits int) error {
	if v > math.MaxFloat64 {
		return ErrFloatInf
	}
	if v < -math.MaxFloat64 {
		return ErrFloatInf
	}
	if v != v {
		return ErrFloatNaN
	}
	m.Write(strconv.AppendFloat(m.scratch[:0], v, 'f', -1, bits))
	return nil
}

func encodeFloat32(m *mstate, rv reflect.Value, opts EncodeOpts) error {
	return encodeFloat(m, rv.Float(), 32)
}

func encodeFloat64(m *mstate, rv reflect.Value, opts EncodeOpts) error {
	return encodeFloat(m, rv.Float(), 64)
}

func encodeInt(m *mstate, rv reflect.Value, opts EncodeOpts) error {
	m.Write(strconv.AppendInt(m.scratch[:0], rv.Int(), 10))
	return nil
}

func encodeMarshaler(m *mstate, rv reflect.Value, opts EncodeOpts) error {
	v := rv.Interface().(Marshaler)
	out, err := v.MarshalEON(m.scratch[:0], opts)
	if err != nil {
		return err
	}
	m.Write(out)
	return nil
}

func encodeString(m *mstate, rv reflect.Value, opts EncodeOpts) error {
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
				if opts.inline() {
					m.WriteByte('n')
				} else {
					m.Truncate(start)
					encodeMultiline(m, v)
					return nil
				}
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

func encodeMultiline(m *mstate, v string) {}

func encodeUint(m *mstate, rv reflect.Value, opts EncodeOpts) error {
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
	actual, loaded := encoders.LoadOrStore(rt, func(m *mstate, rv reflect.Value, opts EncodeOpts) error {
		wg.Wait()
		if err != nil {
			return err
		}
		return enc(m, rv, opts)
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
	if err = enc(m, rv, OptToplevel); err != nil {
		mstates.Put(m)
		return nil, err
	}
	out := make([]byte, m.Len())
	copy(out, m.Bytes())
	mstates.Put(m)
	return out, nil
}

func newArrayEncoder(rt reflect.Type) (encoder, error) {
	elem, err := getEncoder(rt.Elem())
	if err != nil {
		return nil, err
	}
	return (&arrayEncoder{
		elem: elem,
		size: rt.Len(),
	}).encode, nil
}

func newMapEncoder(rt reflect.Type) (encoder, error) {
	k, err := getEncoder(rt.Key())
	if err != nil {
		return nil, err
	}
	v, err := getEncoder(rt.Elem())
	if err != nil {
		return nil, err
	}
	return (&mapEncoder{
		key:   k,
		value: v,
	}).encode, nil
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

func newSliceEncoder(rt reflect.Type) (encoder, error) {
	elem, err := getEncoder(rt.Elem())
	if err != nil {
		return nil, err
	}
	return (&sliceEncoder{
		elem:   elem,
		inline: true,
	}).encode, nil
}

func newStructEncoder(rt reflect.Type) (encoder, error) {
	var fields []*fieldEncoder
	n := rt.NumField()
	for i := 0; i < n; i++ {
		f := rt.Field(i)
		if f.Anonymous {
			continue
		}
		tag := f.Tag.Get("eon")
		_ = tag
		enc, err := getEncoder(f.Type)
		if err != nil {
			return nil, err
		}
		fields = append(fields, &fieldEncoder{
			enc: enc,
			idx: i,
		})
	}
	return (&structEncoder{
		fields: fields,
	}).encode, nil
}

func typeEncoder(rt reflect.Type) (encoder, error) {
	if rt.Implements(marshalerType) {
		return encodeMarshaler, nil
	}
	kind := rt.Kind()
	switch kind {
	case reflect.Array:
		return newArrayEncoder(rt)
	case reflect.Bool:
		return encodeBool, nil
	case reflect.Float32:
		return encodeFloat32, nil
	case reflect.Float64:
		return encodeFloat64, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		return encodeInt, nil
	case reflect.Int64:
		if rt == durationType {
			return encodeDuration, nil
		}
		return encodeInt, nil
	case reflect.Map:
		return newMapEncoder(rt)
	case reflect.Ptr:
	case reflect.Slice:
		if rt.Elem().Kind() == reflect.Uint8 {
			return encodeByteSlice, nil
		}
		return newSliceEncoder(rt)
	case reflect.String:
		return encodeString, nil
	case reflect.Struct:
		return newStructEncoder(rt)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return encodeUint, nil
	case reflect.Uint64:
		if rt == bytesizeType {
			return encodeByteSize, nil
		}
		return encodeUint, nil
	}
	return nil, fmt.Errorf("eon: could not create encoder for %s", rt)
}
