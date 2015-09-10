// "THE BEER-WARE LICENSE" (Revision 42):
// <tobias.rehbein@web.de> wrote this file. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you
// think this stuff is worth it, you can buy me a beer in return.
//                                                             Tobias Rehbein

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"path"

	"golang.org/x/net/websocket"

	"github.com/blabber/jollanotifications/internal/jn"
)

// backlogHandler returns a http.Handler serving a JSON representation of the
// current notification backlog.
func backlogHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logHTTPRequest(r)

		j, err := json.Marshal(s.backlog.Notifications())
		if err != nil {
			log.Panic(err)
		}

		w.Write(j)
	})
}

// rootHandler returns a http.Handler serving files from *htmlDir.
func rootHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logHTTPRequest(r)
		http.ServeFile(w, r, path.Join(*htmlDir, r.URL.Path))
	})
}

// logHTTPRequests logs *http.Request r.
func logHTTPRequest(r *http.Request) {
	if *verbose {
		log.Printf("Request from %v: %v %v", r.RemoteAddr, r.Method, r.URL.Path)
	}
}

// websocketHandler returns a websocket.Handler serving new notifications.
func websocketHandler() websocket.Handler {
	return websocket.Handler(func(ws *websocket.Conn) {
		if *verbose {
			log.Printf("Websocket connected")
		}

		c := make(chan *jn.Notification)
		id := s.websockets.Add(c)

		for n := range c {
			// This would be a good place to check if ws is still
			// connected. The question is how to check for a
			// connection. Until this is answered, rely on the
			// error returned when a message is send over a
			// unconnected websocket.
			err := websocket.JSON.Send(ws, n)
			if err != nil {
				log.Printf("Send: %v", err)
				break
			}
		}

		s.websockets.Remove(id)
		close(c)
		err := ws.Close()
		if err != nil {
			log.Printf("Close: %v", err)
		}
	})
}
