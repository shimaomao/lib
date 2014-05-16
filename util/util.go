package util

import (
	"encoding/json"
	"fmt"
	"github.com/speedland/wcg"
	v "github.com/speedland/wcg/validation"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var logger wcg.Logger

// Returns a singleton logger for non HTTP context.
// If you are in request context, you should use wcg.NewLogger(req)
func GetLogger() wcg.Logger {
	if logger == nil {
		logger = wcg.NewLogger(nil)
	}
	return logger
}

var DefaultWaitForTimeout = 30

func FormatJson(v interface{}) string {
	bytes, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func QueryString(params map[string]string) string {
	q := url.Values{}
	for k, v := range params {
		q.Add(k, v)
	}
	return q.Encode()
}

func WaitFor(f func() bool, seconds int) error {
	if seconds < 0 {
		seconds = DefaultWaitForTimeout
	}
	for c := 0; c < seconds; c++ {
		if res := f(); res {
			return nil
		}
		time.Sleep(time.Duration(1) * time.Second)
	}
	return fmt.Errorf("Timeout in %d seconds", seconds)
}

const ISO8601 = "2006-01-02T15:04:05Z"

// Use ISO8601 format: 2006-01-02T15:04:05Z
func FormatDateTime(t time.Time) string {
	return t.Format(ISO8601)
}

// Use ISO8601 format: 2006-01-02T15:04:05Z
func ParseDateTime(str string) (time.Time, error) {
	return time.Parse(ISO8601, str)
}

func NormalizeDateTime(t time.Time) time.Time {
	return t.Add(time.Duration(-t.Nanosecond()))
}

func NormalizeDate(t time.Time) time.Time {
	return t.Add(time.Duration(-(t.Nanosecond() + t.Second() + t.Minute() + t.Hour())))
}

func ValidateDateTimeString(val interface{}) *v.FieldValidationError {
	switch t := val.(type) {
	case string:
		if _, err := ParseDateTime(val.(string)); err != nil {
			return v.NewFieldValidationError("Datetime format must be ISO8601 in UTC timezone (ending with 'Z')", nil)
		}
	default:
		return v.NewFieldValidationError(
			"not string but {{.Type}}",
			map[string]interface{}{"Type": t},
		)
	}
	return nil
}

func WithTempDir(f func(path string)) {
	p, _ := ioutil.TempDir("", "speedland-apps-temp")
	defer os.RemoveAll(p)
	f(p)
}

func AbsPath(p string) string {
	if strings.HasPrefix(p, "/") {
		return p
	} else {
		current, _ := os.Getwd()
		return filepath.Join(current, p)
	}
}
