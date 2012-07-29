package format

type FormatPart interface {
	Write(interface{}) string
	// TODO read
}

type rawPart struct {
	data string
}

type fieldPart struct {
	name string
}

type Format []FormatPart

func ParseFormat(fmtStr string) (Format, error) {
	return nil, nil
}

func Write(fmtStr string, data interface{}) (string, error) {
	return "", nil
}

func Read(fmtStr string, data interface{}) error {
	return nil
}

