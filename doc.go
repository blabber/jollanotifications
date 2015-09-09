// "THE BEER-WARE LICENSE" (Revision 42):
// <tobias.rehbein@web.de> wrote this file. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you
// think this stuff is worth it, you can buy me a beer in return.
//                                                             Tobias Rehbein

// jollanotifications serves a Jolla phone's notifications via a web interface.
//
// It sniffs for notification events on the dbus and serves a web view
// displaying the last notifications. By default this web view is served via
// "/index.html" on all network interfaces on port 8080.
//
// A JSON encoded representation of the displayed notifications can be accessed
// via "/notifications".
//
// Flags:
//
//	-html string
//	      directory containing the web interface (default "./html")
//	-listen string
//	      network address to listen on (default ":8080")
//	-max int
//	      maximum number of notifications to serve (default 10)
//	-verbose
//	      verbose logging
//
package main
