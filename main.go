// jollanotifications sniffs for notification events on the dbus and serves a
// JSON encoded representation of the last 10 notifications via
// ":8080/notifications". ":8080" serves a web view displaying these
// notifications.
//
// This is used to access a Jolla phones notifications via a web interface.
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/blabber/jollanotifications/internal/jn"
)

var (
	s                state
	maxNotifications int = 10
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
	c := make(chan *Notification)

	go sniffDbus(dbusReader, c)

	go func() {
		for n := range c {
			s.Lock()
			// prepend the new *Notification to s.Notifications
			s.Notifications = append([]*Notification{n}, s.Notifications...)
			// trim s.Notifications to maximum size
			if len(s.Notifications) >= maxNotifications {
				s.Notifications = s.Notifications[:maxNotifications]
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
		http.ServeFile(w, r, strings.Join([]string{"html", r.URL.Path}, "/"))
	})

	panic(http.ListenAndServe(":8080", nil))
}

// dbusReaderFunc is expected to return an io.ReadCloser providing the ouput of
// the monitor-dbus command or an error if something went wrong.
type dbusReaderFunc func() (io.ReadCloser, error)

// dbusReaderMock is a mock function for testing purposes. Opens the file "dbus.log" which should contain
// the captured output of the command line
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
