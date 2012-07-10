package thrift

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

type TestStruct2 struct {
	Str string `thrift:"1"`
}

func (t *TestStruct2) String() string {
	return fmt.Sprintf("{Str:%s}", t.Str)
}

type TestStruct struct {
	String  string            `thrift:"1"`
	Int     *int              `thrift:"2"`
	List    []string          `thrift:"3"`
	Map     map[string]string `thrift:"4"`
	Struct  *TestStruct2      `thrift:"5"`
	List2   []*string         `thrift:"6"`
	Struct2 TestStruct2       `thrift:"7"`
}

func TestKeepEmpty(t *testing.T) {
	s := struct {
		Str1 string `thrift:"1"`
	}{}
	buf := &bytes.Buffer{}
	p := &BinaryProtocol{Writer: buf}
	enc := &Encoder{Protocol: p}
	err := enc.WriteStruct(s)
	if err != nil {
		t.Fatal(err)
	}
	if buf.Len() != 1 || buf.Bytes()[0] != 0 {
		t.Fatal("keepempty didn't work as expected")
	}
}

func TestEncodeRequired(t *testing.T) {
	buf := &bytes.Buffer{}
	p := &BinaryProtocol{Writer: buf}
	enc := &Encoder{Protocol: p}

	s := struct {
		Str1 string `thrift:"1,required"`
	}{}
	err := enc.WriteStruct(s)
	if err != nil {
		t.Fatal(err)
	}
	if buf.Len() != 8 {
		t.Fatal("Non-pointer required fields that aren't 'keepempty' should be serialized empty")
	}

	buf.Reset()
	s2 := struct {
		Str1 *string `thrift:"1,required"`
	}{}
	err = enc.WriteStruct(s2)
	_, ok := err.(*MissingRequiredField)
	if !ok {
		t.Fatalf("Missing required field should throw MissingRequiredField instead of %+v", err)
	}
}

func TestBasics(t *testing.T) {
	i := 123
	str := "bar"
	ts2 := TestStruct2{"qwerty"}
	s := &TestStruct{
		"test",
		&i,
		[]string{"a", "b"},
		map[string]string{"a": "b", "1": "2"},
		&ts2,
		[]*string{&str},
		ts2,
		// nil,
	}
	buf := &bytes.Buffer{}
	p := &BinaryProtocol{Writer: buf, Reader: buf, StrictWrite: true, StrictRead: false}

	enc := &Encoder{Protocol: p}
	err := enc.WriteStruct(s)
	if err != nil {
		t.Fatal(err)
	}

	s2 := &TestStruct{}
	dec := &Decoder{Protocol: p}
	err = dec.ReadStruct(s2)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(s, s2) {
		t.Fatalf("encdec doesn't match: %+v != %+v", s, s2)
	}
}