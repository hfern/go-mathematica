package mathematica

import (
	"bufio"
	"bytes"
	"testing"
)

func expectValue(t *testing.T, v interface{}, expected string) {
	bbuf := new(bytes.Buffer)
	bio := bufio.NewWriter(bbuf)
	e := encodeValue(bio, v)
	if e != nil {
		t.Fatal(e)
		return
	}

	bio.Flush()
	encodedTo := bbuf.String()
	if encodedTo != expected {
		t.Fatalf("Expected `%s`, got `%s` instead", expected, encodedTo)
	}

}

func Test_encodeBool(t *testing.T) {
	expectValue(t, true, "True")
	expectValue(t, false, "False")
}

func Test_encodeFloat(t *testing.T) {
	expectValue(t, 1.2, "1.2")
	expectValue(t, 3.14159, "3.14159")
	expectValue(t, 1.0, "1")
	expectValue(t, -1.0, "-1")
	expectValue(t, 0.0, "0")
}

func Test_encodeInteger(t *testing.T) {
	expectValue(t, 0, "0")
	expectValue(t, 1, "1")
	expectValue(t, 2, "2")
	expectValue(t, -1, "-1")
}

func Test_encodeArrays(t *testing.T) {
	expectValue(t, []string{}, `{}`)
	expectValue(t, []int{}, `{}`)
	expectValue(t, []int{0}, `{0}`)
	expectValue(t, []int{0, -1}, `{0,-1}`)
	expectValue(t, []string{"hunter"}, `{"hunter"}`)
	expectValue(t, []string{"1\n2"}, `{"1\n2"}`)
	expectValue(t, []interface{}{0, []int{1, 2, 3}}, `{0,{1,2,3}}`)
}

func Test_encodeStructs(t *testing.T) {
	expectValue(t, struct{ X int }{X: 1}, `<|"X"->1|>`)
	expectValue(t, struct{ X, Y int }{X: 1}, `<|"X"->1,"Y"->0|>`)

	expectValue(t, struct {
		X int
		Y string
	}{Y: "Hunter"}, `<|"X"->0,"Y"->"Hunter"|>`)

	expectValue(t, struct {
		X int
		Y string
		Z []int
	}{Y: "k", Z: []int{1, 2, 3}}, `<|"X"->0,"Y"->"k","Z"->{1,2,3}|>`)

	expectValue(t, struct {
		X int `mathematica:"Y"`
	}{X: 1}, `<|"Y"->1|>`)
}
