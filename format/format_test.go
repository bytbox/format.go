package format

import (
	"reflect"
	"testing"
)

type test_ParseFormat struct {
	formatString string
	format       Format
}

var tests_ParseFormat []test_ParseFormat = []test_ParseFormat{
	{`a`, Format{rawPart{`a`}}},
	{`abc`, Format{rawPart{`abc`}}},
	{`a$$b`, Format{rawPart{`a$b`}}},
	{`${a}`, Format{fieldPart{`a`}}},
	{`b${a}c`, Format{rawPart{`b`}, fieldPart{`a`}, rawPart{`c`}}},
}

func TestParseFormat(t *testing.T) {
	for _, test := range tests_ParseFormat {
		fmt, err := ParseFormat(test.formatString)
		if err != nil {
			t.Errorf("Unexpected error: %s", err.Error())
		}
		if !reflect.DeepEqual(fmt, test.format) {
			t.Errorf("`%s`: expected %v - got %v",
				test.formatString, test.format, fmt)
		}
	}
}

func BenchParseFormat(b *testing.B) {

}

type testRecord struct {
	A string
	B string
}

type test_Write struct {
	formatString string
	data         interface{}
	result       string
}

var tests_Write []test_Write = []test_Write{
	{``, nil, ``},
	{`ab`, nil, `ab`},
	{`a${B}b${A}`, testRecord{"x", "yz"}, `ayzbx`},
	{`a${B}b${A}`, map[string]string{"A": "x", "B": "yz"}, `ayzbx`},
}

func TestWrite(t *testing.T) {
	for _, test := range tests_Write {
		r, err := Write(test.formatString, test.data)
		if err != nil {
			t.Errorf("Unexpected error: %s", err.Error())
		}
		if r != test.result {
			t.Errorf("`%s`: expected `%s` - got `%s`",
				test.formatString, test.result, r)
		}
	}
}

type test_Read struct {
	formatString string
	input        string
	ptr          interface{}
	data         interface{}
}

var tests_Read []test_Read = []test_Read{
	{`ab`, `ab`, new(testRecord), &testRecord{}},
	{`ab${A}`, `abc`, new(testRecord), &testRecord{`c`, ``}},
	{`ab${A}c`, `abc`, new(testRecord), &testRecord{``, ``}},
	{`a${A}c`, `abc`, new(testRecord), &testRecord{`b`, ``}},
	{`ab${A}c`, `abdc`, &map[string]string{}, &map[string]string{"A": "d"}},
}

func TestRead(t *testing.T) {
	for _, test := range tests_Read {
		err := Read(test.formatString, test.input, test.ptr)
		if err != nil {
			t.Errorf("Unexpected error: %s", err.Error())
		}
		if !reflect.DeepEqual(test.data, test.ptr) {
			t.Errorf("`%s`, `%s`: expected %v - got %v",
				test.formatString, test.input, test.data, test.ptr)
		}
	}
}
