package format

import (
	"errors"
)

// Size of the channel buffer (in runes) used in parsing format strings.
const parseBufLen = 1 << 8

type FormatPart interface {
	Write(interface{}) (string, error)
	Read(string, interface{}) error
}

type rawPart struct {
	data string
}

func (rp rawPart) Write(data interface{}) (string, error) {
	return rp.data, nil
}

func (rp rawPart) Read(input string, data interface{}) error {
	// TODO
	return nil
}

type fieldPart struct {
	name string
}

func (fp fieldPart) Write(data interface{}) (string, error) {
	// TODO
	return "", nil
}

func (fp fieldPart) Read(input string, data interface{}) error {
	// TODO
	return nil
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
		fmt = append(fmt, rawPart{buf})
		buf = ""
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
	return nil, nil
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
