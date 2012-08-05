package format

import (
	"errors"
	"reflect"
)

// Size of the channel buffer (in runes) used in parsing format strings.
const parseBufLen = 1 << 8

type FormatPart interface {
	Write(interface{}) (string, error)

	// Match returns true if this FormatPart can be applied to the
	// beginning of the given string without error.
	Match(string) bool
	Read(string, interface{}, FormatPart) (string, error)
}

type rawPart struct {
	data string
}

func (rp rawPart) Write(data interface{}) (string, error) {
	return rp.data, nil
}

func (rp rawPart) Match(input string) bool {
	return input[:len(rp.data)] == rp.data
}

func (rp rawPart) Read(input string, data interface{}, next FormatPart) (string, error) {
	if len(input) < len(rp.data) || input[:len(rp.data)] != rp.data {
		return "", errors.New("Read failed (expected raw data not found)")
	}
	return input[len(rp.data):], nil
}

type fieldPart struct {
	name string
}

func (fp fieldPart) Write(data interface{}) (string, error) {
	dv := reflect.ValueOf(data)
	switch dv.Type().Kind() {
	case reflect.Struct:
		return dv.FieldByName(fp.name).String(), nil
	case reflect.Map:
		return dv.MapIndex(reflect.ValueOf(fp.name)).String(), nil
	}
	panic("unknown type passed to Write()")
}

func (fp fieldPart) Match(input string) bool {
	return true
}

func (fp fieldPart) Read(input string, data interface{}, next FormatPart) (string, error) {
	var i int
	for i = 0; i < len(input); i++ {
		if next.Match(input[i:]) {
			break
		}
	}
	dv := reflect.Indirect(reflect.ValueOf(data))
	switch dv.Type().Kind() {
	case reflect.Struct:
		dv.FieldByName(fp.name).SetString(input[:i])
		return input[i:], nil
	case reflect.Map:
		dv.SetMapIndex(reflect.ValueOf(fp.name), reflect.ValueOf(input[:i]))
		return input[i:], nil
	}
	panic("unknown type passed to Write()")
}

type terminalPart struct{}

func (tp terminalPart) Write(data interface{}) (string, error) {
	return "", nil
}

func (tp terminalPart) Match(input string) bool {
	return len(input) == 0
}

func (tp terminalPart) Read(input string, data interface{}, next FormatPart) (string, error) {
	if len(input) > 0 {
		return "", errors.New("unexpected data at end of string")
	}
	return "", nil
}

// A Format describes a particular way of representing a record as a string.
type Format []FormatPart

func (fmt Format) Write(data interface{}) (string, error) {
	s := ""
	for _, f := range fmt {
		ts, err := f.Write(data)
		s += ts
		if err != nil {
			return s, err
		}
	}
	return s, nil
}

func (fmt Format) Read(input string, data interface{}) error {
	for i := 0; i < len(fmt); i++ {
		var np FormatPart
		if i+1 == len(fmt) {
			np = terminalPart{}
		} else {
			np = fmt[i+1]
		}
		s, err := fmt[i].Read(input, data, np)
		if err != nil {
			return err
		}
		input = s
	}
	return nil
}

// Creates a Format object that can be used to efficiently read and write
// records from the given format string.
func ParseFormat(fmtStr string) (Format, error) {
	var fmt Format
	var buf string
	var k bool
	rchan := make(chan rune, parseBufLen)

	go func() {
		for _, r := range fmtStr {
			rchan <- r
		}
		close(rchan)
	}()

	var r rune
	buf = ""

raw:
	r, k = <-rchan
	if !k {
		if len(buf) > 0 {
			fmt = append(fmt, rawPart{buf})
		}
		goto done
	}
	if r == '$' {
		r, k = <-rchan
		if !k {
			return nil, errors.New("Format string ended unexpectedly")
		}
		if r == '$' {
			buf += string('$')
			goto raw
		}
		if r != '{' {
			return nil, errors.New("Invalid character sequence: $" + string(r))
		}
		if len(buf) > 0 {
			fmt = append(fmt, rawPart{buf})
			buf = ""
		}
		goto field
	}
	buf += string(r)
	goto raw

field:
	r, k = <-rchan
	if !k {
		return nil, errors.New("Format string ended unexpectedly")
	}
	if r == '}' {
		fmt = append(fmt, fieldPart{buf})
		buf = ""
		goto raw
	}
	buf += string(r)
	goto field

done:
	return fmt, nil
}

// See ParseFormat for a description of the format string.
func Write(fmtStr string, data interface{}) (string, error) {
	fmt, err := ParseFormat(fmtStr)
	if err != nil {
		return "", err
	}
	return fmt.Write(data)
}

// See ParseFormat for a description of the format string.
func Read(fmtStr string, input string, data interface{}) error {
	fmt, err := ParseFormat(fmtStr)
	if err != nil {
		return err
	}
	return fmt.Read(input, data)
}
