package jn

import (
	"bufio"
	"strings"
)

// Notification represents a notification.
type Notification struct {
	// Body is the body of the notification.
	Body string

	// Summary is the summary of the notification. This is misleading
	// though, as this normally describes the source of the notification.
	Summary string
}

// NewNotificationFromMonitorString returns the *Notification represented by
// the dbus-monitor output string ms.
func NewNotificationFromMonitorString(ms string) (*Notification, error) {
	body := false
	summary := false

	n := &Notification{
		Body:    "",
		Summary: "",
	}

	s := bufio.NewScanner(strings.NewReader(ms))
	for s.Scan() {
		if body {
			n.Body = extractString(s.Text())
			body = false
			continue
		}
		if summary {
			n.Summary = extractString(s.Text())
			summary = false
			continue
		}

		if strings.Contains(s.Text(), `string "x-nemo-preview-body"`) {
			body = true
			continue
		}
		if strings.Contains(s.Text(), `string "x-nemo-owner"`) {
			summary = true
			continue
		}
		if strings.Contains(s.Text(), `string "x-nemo-preview-summary"`) {
			summary = true
			continue
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	return n, nil
}

// extractString removes the substring prefixing the first and suffixing the
// last quotation mark including the quotation marks in s.
func extractString(s string) string {
	return s[strings.Index(s, `"`)+1 : strings.LastIndex(s, `"`)]
}
