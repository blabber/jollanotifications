// "THE BEER-WARE LICENSE" (Revision 42):
// <tobias.rehbein@web.de> wrote this file. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you
// think this stuff is worth it, you can buy me a beer in return.
//                                                             Tobias Rehbein

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
