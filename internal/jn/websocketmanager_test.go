// "THE BEER-WARE LICENSE" (Revision 42):
// <tobias.rehbein@web.de> wrote this file. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you
// think this stuff is worth it, you can buy me a beer in return.
//                                                             Tobias Rehbein

package jn

import (
	"sync"
	"testing"
)

func TestAdd(t *testing.T) {
	times := 10

	m := NewWebsocketManager()

	for i := 1; i <= times; i++ {
		c := make(chan *Notification)
		id := m.Add(c)

		if len(m.channels) != i {
			t.Errorf("len(m.channels) => %v, want %v", len(m.channels), i)
		}

		if mc, ok := m.channels[id]; !ok {
			t.Error("m.channels[id] not existing")
		} else {
			if mc != c {
				t.Errorf("m.channels[id] => %v, want %v", mc, c)
			}
		}
	}
}

func TestRemove(t *testing.T) {
	times := 10

	m := NewWebsocketManager()

	cm := make(map[uint]chan *Notification)
	for i := 0; i < times; i++ {
		c := make(chan *Notification)
		id := m.Add(c)
		cm[id] = c
	}

	wl := times
	for id, c := range cm {
		if mc, ok := m.channels[id]; !ok {
			t.Error("m.channels[id] not existing")
		} else {
			if mc != c {
				t.Errorf("m.channels[id] => %v, want %v", mc, c)
			}
		}

		m.Remove(id)

		wl--
		if len(m.channels) != wl {
			t.Errorf("len(m.channels) => %v, want %v", len(m.channels), wl)
		}

		if _, ok := m.channels[id]; ok {
			t.Error("m.channels[id] still existing")
		}
	}
}

var testNotification = &Notification{
	"Time",
	"Summary",
	"Body",
}

func TestSend(t *testing.T) {
	times := 10

	var wg sync.WaitGroup
	m := NewWebsocketManager()
	success := make([]bool, times)

	for i := 0; i < times; i++ {
		c := make(chan *Notification)
		m.Add(c)

		go func(ii int, cc chan *Notification) {
			wg.Add(1)
			n := <-cc
			if n == testNotification {
				success[ii] = true
			} else {
				t.Errorf("<-cc => %v, want %v", n, testNotification)
			}
			wg.Done()
		}(i, c)
	}

	m.Send(testNotification)
	wg.Wait()

	for i := 0; i < times; i++ {
		if !success[i] {
			t.Error("success[i] false, want true")
		}
	}
}
