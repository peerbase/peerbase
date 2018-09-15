package eon

import (
	"math"
	"strings"
	"testing"
	"time"

	"peerbase.net/go/bytesize"
)

type dummyMarshaler struct {
	val string
}

func (d dummyMarshaler) MarshalEON(scratch []byte, opts EncodeOpts) ([]byte, error) {
	return []byte(strings.ToUpper(d.val)), nil
}

func TestEncodeBool(t *testing.T) {
	type elem struct {
		v      bool
		expect string
	}
	for _, elem := range []elem{
		{true, "true"},
		{false, "false"},
	} {
		out, err := Marshal(elem.v)
		if err != nil {
			t.Errorf("unexpected error when encoding bool %v: %s", elem.v, err)
			continue
		}
		if elem.expect != string(out) {
			t.Errorf("mismatching encoded value for bool: expected %q, got %q", elem.expect, out)
		}
	}
}

func TestEncodeByteSize(t *testing.T) {
	type elem struct {
		v      bytesize.Value
		expect string
	}
	for _, elem := range []elem{
		{bytesize.KB, "1KB"},
		{bytesize.MB, "1MB"},
		{327029 * bytesize.Byte, "327029B"},
		{327 * bytesize.GB, "327GB"},
		{1024 * bytesize.MB, "1GB"},
	} {
		out, err := Marshal(elem.v)
		if err != nil {
			t.Errorf("unexpected error when encoding bytesize %v: %s", elem.v, err)
			continue
		}
		if elem.expect != string(out) {
			t.Errorf("mismatching encoded value for bytesize: expected %q, got %q", elem.expect, out)
		}
	}
}

func TestEncodeComplex(t *testing.T) {
	_, err := Marshal(5i)
	if err == nil {
		t.Fatalf("failed to receive expected error when marshalling complex number")
	}
}

func TestEncodeDuration(t *testing.T) {
	type elem struct {
		v      time.Duration
		expect string
	}
	for _, elem := range []elem{
		{time.Second, "1s"},
		{24 * time.Second, "24s"},
		{100 * time.Nanosecond, "100ns"},
		{24 * time.Hour, "24h0m0s"},
		{327 * time.Minute, "5h27m0s"},
		{1024 * time.Millisecond, "1.024s"},
	} {
		out, err := Marshal(elem.v)
		if err != nil {
			t.Errorf("unexpected error when encoding duration %v: %s", elem.v, err)
			continue
		}
		if elem.expect != string(out) {
			t.Errorf("mismatching encoded value for duration: expected %q, got %q", elem.expect, out)
		}
	}
}

func TestEncodeFloat32(t *testing.T) {
	type elem struct {
		v      float32
		expect string
	}
	for _, elem := range []elem{
		{1e20, "100000000000000000000"},
		{1e-6, "0.000001"},
		{1e-7, "0.0000001"},
		{1.234567e+06, "1234567"},
		{-0.0, "0"},
		{1.538237820e+22, "15382378000000000000000"},
		{999999999999999868928, "1000000000000000000000"},
	} {
		out, err := Marshal(elem.v)
		if err != nil {
			t.Errorf("unexpected error when encoding float32 %v: %s", elem.v, err)
			continue
		}
		if elem.expect != string(out) {
			t.Errorf("mismatching encoded value for float32: expected %q, got %q", elem.expect, out)
		}
	}
}

func TestEncodeFloat64(t *testing.T) {
	type elem struct {
		v      float64
		expect string
		err    bool
	}
	for _, elem := range []elem{
		{1e20, "100000000000000000000", false},
		{1e-6, "0.000001", false},
		{1e-7, "0.0000001", false},
		{1.234567e+06, "1234567", false},
		{-0.0, "0", false},
		{1.538237820e+22, "15382378200000000000000", false},
		{999999999999999868928, "999999999999999900000", false},
		{math.NaN(), "", true},
		{math.Inf(1), "", true},
		{math.Inf(-1), "", true},
	} {
		out, err := Marshal(elem.v)
		if err != nil {
			if elem.err {
				continue
			}
			t.Errorf("unexpected error when encoding float64 %v: %s", elem.v, err)
			continue
		}
		if elem.expect != string(out) {
			t.Errorf("mismatching encoded value for float64: expected %q, got %q", elem.expect, out)
		}
	}
}

func TestEncodeInt(t *testing.T) {
	type elem struct {
		v      int
		expect string
	}
	for _, elem := range []elem{
		{0, "0"},
		{1, "1"},
		{5, "5"},
		{10, "10"},
		{75, "75"},
		{127, "127"},
		{-128, "-128"},
		{-0, "0"},
		{32767, "32767"},
		{-32768, "-32768"},
		{2147483647, "2147483647"},
		{-2147483648, "-2147483648"},
		{3137093109, "3137093109"},
		{75979721093313246, "75979721093313246"},
		{9223372036854775807, "9223372036854775807"},
		{-0, "0"},
		{-75927941794, "-75927941794"},
		{-9223372036854775808, "-9223372036854775808"},
	} {
		out, err := Marshal(elem.v)
		if err != nil {
			t.Errorf("unexpected error when encoding int %v: %s", elem.v, err)
			continue
		}
		if elem.expect != string(out) {
			t.Errorf("mismatching encoded value for int: expected %q, got %q", elem.expect, out)
		}
	}
}

func TestEncodeInt8(t *testing.T) {
	type elem struct {
		v      int8
		expect string
	}
	for _, elem := range []elem{
		{0, "0"},
		{1, "1"},
		{5, "5"},
		{10, "10"},
		{75, "75"},
		{127, "127"},
		{-128, "-128"},
		{-0, "0"},
	} {
		out, err := Marshal(elem.v)
		if err != nil {
			t.Errorf("unexpected error when encoding int8 %v: %s", elem.v, err)
			continue
		}
		if elem.expect != string(out) {
			t.Errorf("mismatching encoded value for int8: expected %q, got %q", elem.expect, out)
		}
	}
}

func TestEncodeInt16(t *testing.T) {
	type elem struct {
		v      int16
		expect string
	}
	for _, elem := range []elem{
		{0, "0"},
		{1, "1"},
		{5, "5"},
		{10, "10"},
		{75, "75"},
		{127, "127"},
		{-128, "-128"},
		{-0, "0"},
		{32767, "32767"},
		{-32768, "-32768"},
	} {
		out, err := Marshal(elem.v)
		if err != nil {
			t.Errorf("unexpected error when encoding int16 %v: %s", elem.v, err)
			continue
		}
		if elem.expect != string(out) {
			t.Errorf("mismatching encoded value for int16: expected %q, got %q", elem.expect, out)
		}
	}
}

func TestEncodeInt32(t *testing.T) {
	type elem struct {
		v      int32
		expect string
	}
	for _, elem := range []elem{
		{0, "0"},
		{1, "1"},
		{5, "5"},
		{10, "10"},
		{75, "75"},
		{127, "127"},
		{-128, "-128"},
		{-0, "0"},
		{32767, "32767"},
		{-32768, "-32768"},
		{2147483647, "2147483647"},
		{-2147483648, "-2147483648"},
	} {
		out, err := Marshal(elem.v)
		if err != nil {
			t.Errorf("unexpected error when encoding int32 %v: %s", elem.v, err)
			continue
		}
		if elem.expect != string(out) {
			t.Errorf("mismatching encoded value for int32: expected %q, got %q", elem.expect, out)
		}
	}
}

func TestEncodeInt64(t *testing.T) {
	type elem struct {
		v      int64
		expect string
	}
	for _, elem := range []elem{
		{0, "0"},
		{1, "1"},
		{5, "5"},
		{10, "10"},
		{75, "75"},
		{127, "127"},
		{-128, "-128"},
		{-0, "0"},
		{32767, "32767"},
		{-32768, "-32768"},
		{2147483647, "2147483647"},
		{-2147483648, "-2147483648"},
		{9223372036854775807, "9223372036854775807"},
		{-9223372036854775808, "-9223372036854775808"},
	} {
		out, err := Marshal(elem.v)
		if err != nil {
			t.Errorf("unexpected error when encoding int64 %v: %s", elem.v, err)
			continue
		}
		if elem.expect != string(out) {
			t.Errorf("mismatching encoded value for int64: expected %q, got %q", elem.expect, out)
		}
	}
}

func TestEncodeMarshaler(t *testing.T) {
	type elem struct {
		v      string
		expect string
	}
	for _, elem := range []elem{
		{"hello", "HELLO"},
		{"hello-world", "HELLO-WORLD"},
	} {
		v := &dummyMarshaler{elem.v}
		out, err := Marshal(v)
		if err != nil {
			t.Errorf("unexpected error when encoding marshaler with value %q: %s", elem.v, err)
			continue
		}
		if elem.expect != string(out) {
			t.Errorf("mismatching encoded value for marshaler: expected %s, got %s", elem.expect, out)
		}
	}
}

func TestEncodeSlice(t *testing.T) {
	type elem struct {
		v      interface{}
		expect string
	}
	for _, elem := range []elem{
		{[]int{1, 2, 3}, `[1 2 3]`},
	} {
		out, err := Marshal(elem.v)
		if err != nil {
			t.Errorf("unexpected error when encoding slice %v: %s", elem.v, err)
			continue
		}
		if elem.expect != string(out) {
			t.Errorf("mismatching encoded value for slice: expected %s, got %s", elem.expect, out)
		}
	}
}
func TestEncodeString(t *testing.T) {
	type elem struct {
		v      string
		expect string
	}
	for _, elem := range []elem{
		{"hello", `"hello"`},
		{"hello world", `"hello world"`},
		{"\x00\t\r\"", `"\x00\t\r\""`},
		{"héllo \xff world", `"héllo \ufffd world"`},
	} {
		out, err := Marshal(elem.v)
		if err != nil {
			t.Errorf("unexpected error when encoding string %q: %s", elem.v, err)
			continue
		}
		if elem.expect != string(out) {
			t.Errorf("mismatching encoded value for string: expected %s, got %s", elem.expect, out)
		}
	}
}
func TestEncodeUint32(t *testing.T) {
	type elem struct {
		v      uint32
		expect string
	}
	for _, elem := range []elem{
		{0, "0"},
		{1, "1"},
		{5, "5"},
		{10, "10"},
		{75, "75"},
		{127, "127"},
		{-0, "0"},
		{32767, "32767"},
		{2147483647, "2147483647"},
	} {
		out, err := Marshal(elem.v)
		if err != nil {
			t.Errorf("unexpected error when encoding uint32 %v: %s", elem.v, err)
			continue
		}
		if elem.expect != string(out) {
			t.Errorf("mismatching encoded value for uint32: expected %q, got %q", elem.expect, out)
		}
	}
}

func TestEncodeUint64(t *testing.T) {
	type elem struct {
		v      uint64
		expect string
	}
	for _, elem := range []elem{
		{0, "0"},
		{1, "1"},
		{5, "5"},
		{10, "10"},
		{75, "75"},
		{127, "127"},
		{-0, "0"},
		{32767, "32767"},
		{2147483647, "2147483647"},
		{9223372036854775807, "9223372036854775807"},
	} {
		out, err := Marshal(elem.v)
		if err != nil {
			t.Errorf("unexpected error when encoding uint64 %v: %s", elem.v, err)
			continue
		}
		if elem.expect != string(out) {
			t.Errorf("mismatching encoded value for uint64: expected %q, got %q", elem.expect, out)
		}
	}
}
