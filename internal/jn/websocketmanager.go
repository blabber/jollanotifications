// "THE BEER-WARE LICENSE" (Revision 42):
// <tobias.rehbein@web.de> wrote this file. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you
// think this stuff is worth it, you can buy me a beer in return.
//                                                             Tobias Rehbein

package jn

import (
	"sync"
)

// WebsocketManager manages a collection of channels, one for each websocket
// connection. The Send method can be used to broadcast a notification to all
// websockets.
//
// WebsocketManager methods are synchronized and can be called from concurrent
// goroutines without additional synchronization.
type WebsocketManager struct {
	sync.RWMutex
	maxId    uint
	channels map[uint]chan<- *Notification
}

// NewWebsocketManager returns a new *WebsocketManager instance.
func NewWebsocketManager() *WebsocketManager {
	return &WebsocketManager{channels: make(map[uint]chan<- *Notification)}
}

// Add adds channel c to the collection of websocket channels. An ID is
// returned that can be used to remove the channel from the collection by
// calling Remove.
func (m *WebsocketManager) Add(c chan<- *Notification) uint {
	defer func() {
		m.maxId++
	}()

	m.Lock()
	defer m.Unlock()

	m.channels[m.maxId] = c

	return m.maxId
}

// Remove removes the channel identified by id.
func (m *WebsocketManager) Remove(id uint) {
	m.Lock()

	delete(m.channels, id)

	m.Unlock()
}

// Send broadcasts the notification n on all channels managed by
// WebsocketManager.
func (m *WebsocketManager) Send(n *Notification) {
	m.RLock()

	for _, c := range m.channels {
		c <- n
	}

	m.RUnlock()
}
