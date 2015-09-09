// "THE BEER-WARE LICENSE" (Revision 42):
// <tobias.rehbein@web.de> wrote this file. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you
// think this stuff is worth it, you can buy me a beer in return.
//                                                             Tobias Rehbein

package jn

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"time"
)

// TimeFormatter is a function that returns a string representation of
// time.Time t.
type TimeFormatter func(time time.Time) string

// Notification represents a notification.
type Notification struct {
	// Time is a string representing the time when the notification
	// occured.
	Time string

	// Summary is the summary of the notification. This is misleading
	// though, as this normally describes the source of the notification.
	Summary string

	// Body is the body of the notification.
	Body string
}

// NewNotificationFromMonitorString returns the *Notification represented by
// the dbus-monitor output string ms. tf is called with the current time and
// the returned string is used as the value for the Time field of the newly
// created *Notification.
func NewNotificationFromMonitorString(ms string, tf TimeFormatter) (*Notification, error) {
	body := false
	summary := false

	n := &Notification{
		Time: tf(time.Now()),
	}

	s := bufio.NewScanner(strings.NewReader(ms))
	for s.Scan() {
		switch {
		case body:
			n.Body = extractString(s.Text())
			body = false
		case summary:
			n.Summary = extractString(s.Text())
			summary = false
		case strings.Contains(s.Text(), `string "x-nemo-preview-body"`):
			body = true
		case strings.Contains(s.Text(), `string "x-nemo-owner"`):
			if len(n.Summary) == 0 {
				summary = true
			}
		case strings.Contains(s.Text(), `string "x-nemo-preview-summary"`):
			summary = true
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	return n, nil
}

// IsEmpty returns true if *Notification n does not contain usable data. This
// is the case if n does not contain a non-empty Body.
func (n *Notification) IsEmpty() bool {
	if len(n.Body) > 0 {
		return false
	}
	return true
}

// String returns a string rerpresentation of *Notification n.
func (n *Notification) String() string {
	if n.IsEmpty() {
		return "Empty notification"
	}

	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("Notification time: \"%s\" ", n.Time))
	if len(n.Summary) > 0 {
		b.WriteString(fmt.Sprintf("summary: \"%s\" ", n.Summary))
	}
	b.WriteString(fmt.Sprintf("body: \"%s\"", n.Body))

	return b.String()
}

// extractString removes the substring prefixing the first and suffixing the
// last quotation mark including the quotation marks in s.
func extractString(s string) string {
	return s[strings.Index(s, `"`)+1 : strings.LastIndex(s, `"`)]
}
