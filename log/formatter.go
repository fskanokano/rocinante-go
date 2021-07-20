package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
)

type Formatter struct {
	MessageFormat string
	TimeFormat    string
	ForceColor    bool
}

func DefaultFormatter() *Formatter {
	return &Formatter{
		MessageFormat: defaultMessageFormat,
		TimeFormat:    defaultTimeFormat,
		ForceColor:    true,
	}
}

func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	if f.MessageFormat == ""{
		f.MessageFormat = defaultMessageFormat
	}

	if f.TimeFormat == ""{
		f.TimeFormat = defaultTimeFormat
	}

	messageFormat := f.MessageFormat

	time := entry.Time.Format(f.TimeFormat)
	if f.ForceColor {
		levelColor := getColorByLevel(entry.Level)

		messageFormat = strings.Replace(messageFormat, "|", fmt.Sprintf("\u001B[%dm|\u001B[0m", purpleRed), -1)

		messageFormat = strings.Replace(messageFormat, "{time}", fmt.Sprintf("\u001B[%dm%s\u001B[0m", green, time), -1)

		messageFormat = strings.Replace(messageFormat, "{level}", fmt.Sprintf("\u001B[%dm%s\u001B[0m", levelColor, strings.ToUpper(entry.Level.String())), -1)

		messageFormat = strings.Replace(messageFormat, "{message}", fmt.Sprintf("\u001B[%dm%s\u001B[0m", levelColor, entry.Message), -1)
	} else {
		messageFormat = strings.Replace(messageFormat, "{time}", fmt.Sprintf("%s", time), -1)

		messageFormat = strings.Replace(messageFormat, "{level}", fmt.Sprintf("%s", strings.ToUpper(entry.Level.String())), -1)

		messageFormat = strings.Replace(messageFormat, "{message}", fmt.Sprintf("%s", entry.Message), -1)
	}

	messageFormat += "\n"

	return []byte(messageFormat), nil
}

func getColorByLevel(level logrus.Level) int {
	switch level {
	case logrus.DebugLevel, logrus.TraceLevel:
		return blue
	case logrus.WarnLevel:
		return yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return red
	default:
		return cyanBlue
	}
}

const (
	defaultMessageFormat = "{time} | {level} | {message}"
	defaultTimeFormat    = "2006-01-02 15:04:05"
)

const (
	red = iota + 31
	green
	yellow
	blue
	purpleRed
	cyanBlue
)
