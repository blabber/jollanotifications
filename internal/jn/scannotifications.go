package jn

import (
	"bytes"
)

// ScanNotifications is a split function for bufio.Scanner that returns each
// record of dbus-monitor output.
func ScanNotifications(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if len(data) == 0 && atEOF {
		return 0, nil, nil
	}

	target := "int32 -1\n"

	i := bytes.Index(data, []byte(target))
	if i == -1 && !atEOF {
		// request more data
		return 0, nil, nil
	}

	if atEOF {
		return len(data), data, nil
	}

	return i + len(target), data[0 : i+len(target)-1], nil
}
