// "THE BEER-WARE LICENSE" (Revision 42):
// <tobias.rehbein@web.de> wrote this file. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you
// think this stuff is worth it, you can buy me a beer in return.
//                                                             Tobias Rehbein

package jn

import (
	"bufio"
	"strings"
)

// Notification represents a notification.
type Notification struct {
	// Summary is the summary of the notification. This is misleading
	// though, as this normally describes the source of the notification.
	Summary string

	// Body is the body of the notification.
	Body string
}

// NewNotificationFromMonitorString returns the *Notification represented by
// the dbus-monitor output string ms.
func NewNotificationFromMonitorString(ms string) (*Notification, error) {
	body := false
	summary := false

	n := &Notification{}

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
			if len(n.Summary) == 0 {
				summary = true
			}
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

// IsEmpty returns true if *Notification n does not contain usable data. This
// is the case if n does not contain a non-empty Body.
func (n *Notification) IsEmpty() bool {
	if len(n.Body) > 0 {
		return false
	}
	return true
}

// extractString removes the substring prefixing the first and suffixing the
// last quotation mark including the quotation marks in s.
func extractString(s string) string {
	return s[strings.Index(s, `"`)+1 : strings.LastIndex(s, `"`)]
}
