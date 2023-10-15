package utils

import (
	"time"
)

func ISODate(t time.Time) string {
	return t.Format(time.RFC3339)
}

func ISODateNow() string {
	return ISODate(time.Now().UTC())
}

func ToLocalDate(date string) string {
	if date == "" {
		return ""
	}
	t, _ := time.Parse(time.RFC3339, date)
	return ISODate(t.Local())
}
