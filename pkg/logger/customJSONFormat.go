package logger

import (
	"bytes"

	"github.com/sirupsen/logrus"
)

type CustomJSONFormat struct{}

func (format *CustomJSONFormat) Format(logger *logrus.Entry) ([]byte, error) {
	var msg string
	b := &bytes.Buffer{}
	b.WriteString(msg)
	return b.Bytes(), nil
}
