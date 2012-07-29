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

type test_Write struct {
	formatString string
}

var tests_Write []test_Write = []test_Write{}

func TestWrite(t *testing.T) {
	for _, test := range tests_Write {
		println(test.formatString)
	}
}

func TestRead(t *testing.T) {

}
