package format

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

type Format []FormatPart

func (fmt Format) Write(data interface{}) (string, error) {
	return "", nil
}

func (fmt Format) Read(input string, data interface{}) error {
	return nil
}

func ParseFormat(fmtStr string) (Format, error) {
	return nil, nil
}

func Write(fmtStr string, data interface{}) (string, error) {
	fmt, err := ParseFormat(fmtStr)
	if err != nil {
		return "", err
	}
	return fmt.Write(data)
}

func Read(fmtStr string, input string, data interface{}) error {
	fmt, err := ParseFormat(fmtStr)
	if err != nil {
		return err
	}
	return fmt.Read(input, data)
}

