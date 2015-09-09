// "THE BEER-WARE LICENSE" (Revision 42):
// <tobias.rehbein@web.de> wrote this file. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you
// think this stuff is worth it, you can buy me a beer in return.
//                                                             Tobias Rehbein

package jn

import (
	"sync"
)

// Backlog contains the current backlog of notifications. It holds a maximum
// number of notifications, as defined by the size parameter to NewBacklog.
//
// When a new item is added to the backlog and the number of items exceeds the
// maximum number, the oldest item in the backlog is deleted.
//
// To items in the backlog can be retrieved by calling Notifications.
//
// Add and Notifications are synchronized and can be called from concurrent go
// routines without additional synchronization.
type Backlog struct {
	mtx           sync.RWMutex
	size          int
	notifications []*Notification
}

// NewBacklog creates a new *Backlog, configured to hold a maximum number of
// examinations equal to size.
func NewBacklog(size int) *Backlog {
	return &Backlog{size: size}
}

// Add adds a new *Notification to the backlog. If the new number of items
// exceeds the maximum number of items, the oldest item is removed.
func (b *Backlog) Add(n *Notification) {
	b.mtx.Lock()
	// prepend the new *Notification to s.Notifications
	b.notifications = append([]*Notification{n}, b.notifications...)
	// trim s.Notifications to maximum size
	if len(b.notifications) >= b.size {
		b.notifications = b.notifications[:b.size]
	}
	b.mtx.Unlock()
}

// Notifications returns the items contained in the backlog as a slice of
// *Notification. The items are sorted. The newest item has index 0, the oldest
// index size-1.
func (b *Backlog) Notifications() []*Notification {
	b.mtx.RLock()
	defer b.mtx.RUnlock()

	return b.notifications
}
