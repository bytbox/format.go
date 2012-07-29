package format

import (
	"testing"
)

type test_ParseFormat struct {
	formatString string
	format       Format
}

var tests_ParseFormat []test_ParseFormat = []test_ParseFormat{}

func TestParseFormat(t *testing.T) {
	for _, test := range tests_ParseFormat {
		println(test.formatString)
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
