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
