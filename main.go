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
//
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"sync"
	"time"

	"github.com/blabber/jollanotifications/internal/jn"
)

var (
	s                state
	maxNotifications = flag.Int("max", 10, "maximum number of notifications to serve")
	networkAddress   = flag.String("listen", ":8080", "network address to listen on")
	htmlDir          = flag.String("html", "./html", "directory containing the web interface")
)

// state represents the shared state used by the sniffer and served via web. An
// embedded sync.RWMutex is used for synchronization.
type state struct {
	sync.RWMutex

	// Notifications is a slice of the 10 last notifications, represented
	// by *Notification, or less if fewer notifications occured.
	Notifications []*Notification
}

// Notification represents a time stamped notification.
type Notification struct {
	*jn.Notification

	// Time is a string representation of the time when the notification
	// occured.
	Time string
}

func main() {
	flag.Parse()

	c := make(chan *Notification)
	go sniffDbus(dbusReader, c)

	go func() {
		for n := range c {
			s.Lock()
			// prepend the new *Notification to s.Notifications
			s.Notifications = append([]*Notification{n}, s.Notifications...)
			// trim s.Notifications to maximum size
			if len(s.Notifications) >= *maxNotifications {
				s.Notifications = s.Notifications[:*maxNotifications]
			}
			s.Unlock()
		}
	}()

	http.HandleFunc("/notifications", func(w http.ResponseWriter, r *http.Request) {
		s.RLock()
		j, err := json.Marshal(s)
		if err != nil {
			panic(err)
		}
		s.RUnlock()
		w.Write(j)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path.Join(*htmlDir, r.URL.Path))
	})

	panic(http.ListenAndServe(*networkAddress, nil))
}

// dbusReaderFunc is expected to return an io.ReadCloser providing the ouput of
// the monitor-dbus command or an error if something went wrong.
type dbusReaderFunc func() (io.ReadCloser, error)

// dbusReaderMock is a mock function for testing purposes. Opens the file
// "dbus.log" which should contain the captured output of the command line
//
//    dbus-monitor "interface='org.freedesktop.Notifications',member='Notify'"
func dbusReaderMock() (io.ReadCloser, error) {
	return os.Open("dbus.log")
}

// dbusReader starts dbus-monitor configured to sniff for notification events.
func dbusReader() (io.ReadCloser, error) {
	c := exec.Command("dbus-monitor", "interface='org.freedesktop.Notifications',member='Notify'")
	r, err := c.StdoutPipe()
	if err != nil {
		return nil, err
	}
	err = c.Start()
	if err != nil {
		return nil, err
	}

	return r, nil
}

// sniffDbus scans the io.ReadCloser returned by rf for records and returns
// *Notification via out.
func sniffDbus(rf dbusReaderFunc, out chan<- *Notification) {
	r, err := rf()
	if err != nil {
		panic(err)
	}
	defer r.Close()

	s := bufio.NewScanner(r)
	s.Split(jn.ScanNotifications)
	for s.Scan() {
		n, err := jn.NewNotificationFromMonitorString(s.Text())
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERR: %v\n", err)
			continue
		}

		if n.IsEmpty() {
			continue
		}

		out <- &Notification{
			n,
			time.Now().Format(time.RFC822),
		}
	}
	if err := s.Err(); err != nil {
		panic(err)
	}

	close(out)
}
