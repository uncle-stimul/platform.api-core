package logger

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

type CustomStringFormat struct{}

func (format *CustomStringFormat) Format(logger *logrus.Entry) ([]byte, error) {
	var msg string

	if logger.Level <= logrus.WarnLevel && logger.Data[logrus.ErrorKey] != nil {
		errRaw := fmt.Sprintf("%v", logger.Data[logrus.ErrorKey])
		errRawString := strings.Split(errRaw, "\n")
		errMsg := strings.Join(errRawString, "\n\t")
		msg = fmt.Sprintf(
			"%s %s: %s \n\t%v\n",
			logger.Time.Format("2006-01-02 15:04:05"),
			logger.Level.String(),
			logger.Message,
			errMsg,
		)
	} else {
		msg = fmt.Sprintf(
			"%s %s: %s \n",
			logger.Time.Format("2006-01-02 15:04:05"),
			logger.Level.String(),
			logger.Message,
		)
	}

	b := &bytes.Buffer{}
	b.WriteString(msg)
	return b.Bytes(), nil
}
