package jn

import (
	"bufio"
	"fmt"
	"strings"
	"testing"
)

var monitorStrings = []string{threemaMonitorString, clockMonitorString, commhistorydMonitorString}
var completeMonitorString = strings.Join(monitorStrings, "\n")

func TestScanNotifications(t *testing.T) {
	s := bufio.NewScanner(strings.NewReader(completeMonitorString))
	s.Split(ScanNotifications)
	i := 0
	for s.Scan() {
		m := monitorStrings[i]
		i++
		if s.Text() != m {
			t.Errorf("%#v != %#v", s.Text(), m)
		}
	}
	if err := s.Err(); err != nil {
		fmt.Printf("%v\n", err)
		t.Error(err)
	}
}
