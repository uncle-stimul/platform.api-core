package logger

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

type customStringFormat struct{}

func (format *customStringFormat) Format(logger *logrus.Entry) ([]byte, error) {
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

type customJSONFormat struct{}

func (format *customJSONFormat) Format(logger *logrus.Entry) ([]byte, error) {
	var b bytes.Buffer

	var errString string
	if logger.Level == logrus.ErrorLevel && logger.Data[logrus.ErrorKey] != nil {
		errRaw := fmt.Sprintf("%v", logger.Data[logrus.ErrorKey])
		errLines := strings.Split(errRaw, "\n")
		errString = strings.Join(errLines, "\n\t")
	}

	b.WriteString("{")
	b.WriteString(fmt.Sprintf(`"time": "%s", `, logger.Time.Format("2006-01-02 15:04:05")))
	b.WriteString(fmt.Sprintf(`"level": "%s", `, logger.Level.String()))
	b.WriteString(fmt.Sprintf(`"message": "%s"`, logger.Message))

	if errString != "" {
		b.WriteString(fmt.Sprintf(`, "error_message": "%s"`, errString))
	}

	b.WriteString("}\n")
	return b.Bytes(), nil
}
