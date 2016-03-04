package mathematica

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

var (
	CantEncodeValue = errors.New("Can't encode this value!")
)

type Marshaler interface {
	MarshalMathematica() ([]byte, error)
}

func Marshal(v interface{}) ([]byte, error) {
	bbuf := new(bytes.Buffer)
	bio := bufio.NewWriter(bbuf)
	e := encodeValue(bio, v)
	if e != nil {
		return []byte{}, e
	}

	bio.Flush()
	return e, bbuf.Bytes()
}

func encodeValue(buf *bufio.Writer, v interface{}) (e error) {
	if m, ok := v.(Marshaler); ok {
		if m != nil {
			var byt []byte
			byt, e = m.MarshalMathematica()
			buf.Write(byt)
			return
		}
	}

	// quick-cases to skip reflection
	switch t := v.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		_, e = fmt.Fprint(buf, v)
		return

	case float32, float64:
		_, e = fmt.Fprint(buf, v)
		return

	case bool:
		return encodeBool(buf, t)

	case string:
		return encodeString(buf, t)

	case fmt.Stringer:
		return encodeString(buf, t.String())
	}

	vt := reflect.ValueOf(v)
	return encodeValueReflect(buf, vt)
}

func encodeBool(buf *bufio.Writer, v bool) (e error) {
	switch v {
	case true:
		_, e = buf.WriteString("True")
	case false:
		_, e = buf.WriteString("False")
	}
	return
}

func encodeString(buf *bufio.Writer, v string) (e error) {
	_, e = buf.WriteString(strconv.Quote(v))
	return
}

func encodeValueReflect(buf *bufio.Writer, vt reflect.Value) (e error) {
	// deal with arrays and structs and other non-standard values
	switch vt.Kind() {
	case reflect.Array, reflect.Slice:
		return encodeArray(buf, vt)
	case reflect.Struct:
		return encodeStruct(buf, vt)
	default:
		return CantEncodeValue
	}
	return nil
}

func encodeArray(buf *bufio.Writer, vt reflect.Value) (e error) {
	l := vt.Len()
	buf.WriteByte('{')

	for i := 0; i < l; i++ {
		if i != 0 {
			buf.WriteByte(',')
		}
		e = encodeValue(buf, vt.Index(i).Interface())
		if e != nil {
			return
		}
	}

	buf.WriteByte('}')

	return
}

func encodeStruct(buf *bufio.Writer, vt reflect.Value) (e error) {

	tt := vt.Type()
	n := tt.NumField()

	first := true

	buf.WriteString("<|")

	for i := 0; i < n; i++ {
		f := tt.Field(i)

		if f.PkgPath != "" {
			// this field is unexported!!
			continue
		}

		fname := f.Tag.Get("mathematica")
		if fname == "" {
			fname = f.Name
		}

		if first {
			first = false
		} else {
			buf.WriteByte(',')
		}

		encodeString(buf, fname)

		buf.WriteString("->")

		e = encodeValue(buf, vt.Field(i).Interface())
		if e != nil {
			return
		}
	}

	buf.WriteString("|>")

	return
}
