package logger

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
	"sort"
	"strings"
	"time"
)

type Formatter struct {
	TimestampFormat string

	CallerFirst bool
	// CustomCallerFormatter - set custom formatter for caller info
	CustomCallerFormatter func(*runtime.Frame) string
}

// Format an log entry
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	levelColor := getColorByLevel(entry.Level)

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = time.RFC3339
	}

	// output buffer
	b := &bytes.Buffer{}

	fmt.Fprintf(b, "\x1b[%dm", levelColor)

	f.writeCaller(b, entry)

	// write level
	b.WriteString(" [")
	level := strings.ToUpper(entry.Level.String())
	b.WriteString(level[:4])
	b.WriteString("] ")
	// write time
	b.WriteString(entry.Time.Format(timestampFormat))

	b.WriteString(" ")

	// write fields
	f.writeFields(b, entry)

	b.WriteString(" ")

	// write message
	b.WriteString(strings.TrimSpace(entry.Message))

	fmt.Fprintf(b, " \x1b[0m ")

	b.WriteByte('\n')

	return b.Bytes(), nil
}

func (f *Formatter) writeCaller(b *bytes.Buffer, entry *logrus.Entry) {
	if entry.HasCaller() {
		if f.CustomCallerFormatter != nil {
			fmt.Fprintf(b, f.CustomCallerFormatter(entry.Caller))
		} else {
			fmt.Fprintf(
				b,
				" (%s:%d %s)",
				entry.Caller.File,
				entry.Caller.Line,
				entry.Caller.Function,
			)
		}
	}
}

func (f *Formatter) writeFields(b *bytes.Buffer, entry *logrus.Entry) {
	if len(entry.Data) != 0 {
		fields := make([]string, 0, len(entry.Data))
		for field := range entry.Data {
			fields = append(fields, field)
		}

		sort.Strings(fields)

		for _, field := range fields {
			f.writeField(b, entry, field)
		}
	}
}

func (f *Formatter) writeField(b *bytes.Buffer, entry *logrus.Entry, field string) {
	fmt.Fprintf(b, "%s=%v", field, entry.Data[field])

	b.WriteString(" ")
}

const (
	colorRed    = 31
	colorYellow = 33
	colorBlue   = 36
	colorGray   = 37
)

func getColorByLevel(level logrus.Level) int {
	switch level {
	case logrus.DebugLevel, logrus.TraceLevel:
		return 0
	case logrus.WarnLevel:
		return colorYellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return colorRed
	default:
		return 0
	}
}
