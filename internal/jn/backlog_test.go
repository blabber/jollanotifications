// "THE BEER-WARE LICENSE" (Revision 42):
// <tobias.rehbein@web.de> wrote this file. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you
// think this stuff is worth it, you can buy me a beer in return.
//                                                             Tobias Rehbein

package jn

import (
	"testing"
)

var tests = []*Notification{
	{"Time_1", "Summary_1", "Body_1"},
	{"Time_2", "Summary_2", "Body_2"},
	{"Time_3", "Summary_3", "Body_3"},
	{"Time_4", "Summary_4", "Body_4"},
	{"Time_5", "Summary_5", "Body_5"},
}

func testBacklog(s int, t *testing.T) {
	b := NewBacklog(s)

	for _, n := range tests {
		b.Add(n)
	}

	bn := b.Notifications()
	if len(bn) > s {
		t.Errorf("backlog length %v != %v", len(bn), s)
	}

	for i := 0; i < s; i++ {
		actual := bn[i]
		expected := tests[len(tests)-1-i]

		if actual != expected {
			t.Errorf("%v != %v", actual, expected)
		}
	}
}

func TestBacklogSize2(t *testing.T) {
	testBacklog(2, t)
}

func TestBacklogSize3(t *testing.T) {
	testBacklog(3, t)
}

func TestBacklogSize4(t *testing.T) {
	testBacklog(4, t)
}
