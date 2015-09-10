// "THE BEER-WARE LICENSE" (Revision 42):
// <tobias.rehbein@web.de> wrote this file. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you
// think this stuff is worth it, you can buy me a beer in return.
//                                                             Tobias Rehbein

package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/blabber/jollanotifications/internal/jn"
)

const version = "v0.1.0+"

var (
	s state

	verbose          = flag.Bool("verbose", false, "verbose logging")
	maxNotifications = flag.Int("max", 10, "maximum number of notifications to serve")
	networkAddress   = flag.String("listen", ":8080", "network address to listen on")
	htmlDir          = flag.String("html", "./html", "directory containing the web interface")
)

// state represents the shared state used by the sniffer and served via web. An
// embedded sync.RWMutex is used for synchronization.
type state struct {
	backlog    *jn.Backlog
	websockets *jn.WebsocketManager
}

func main() {
	log.Printf("jollanotifications (%v)", version)

	flag.Parse()

	s.backlog = jn.NewBacklog(*maxNotifications)
	s.websockets = jn.NewWebsocketManager()

	go func() {
		c := sniffDbus(dbusReader)

		for n := range c {
			s.backlog.Add(n)
			s.websockets.Send(n)
		}
	}()

	http.Handle("/", rootHandler())
	http.Handle("/notifications", backlogHandler())
	http.Handle("/websocket", websocketHandler())

	log.Printf("Listening on %v", *networkAddress)
	log.Panic(http.ListenAndServe(*networkAddress, nil))
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

// timeFormatter is passed to jn.NewNotificationFromMonitorString to format the
// time when the notification was received.
func timeFormatter(t time.Time) string {
	return time.Now().Format(time.RFC822)
}

// sniffDbus scans the io.ReadCloser returned by rf for records and returns
// *Notification via the returned channel.
func sniffDbus(rf dbusReaderFunc) <-chan *jn.Notification {
	out := make(chan *jn.Notification)

	go func() {
		r, err := rf()
		if err != nil {
			log.Panic(err)
		}
		defer r.Close()

		s := bufio.NewScanner(r)
		s.Split(jn.ScanNotifications)
		for s.Scan() {
			if *verbose {
				log.Printf("D-Bus record: %v", s.Text())
			}

			n, err := jn.NewNotificationFromMonitorString(s.Text(), timeFormatter)
			if err != nil {
				log.Printf("Error: NewNotificationFromMonitorString: %v", err)
				continue
			}

			if *verbose {
				log.Printf("New Notification: %v", n)
			}

			if n.IsEmpty() {
				continue
			}

			out <- n
		}
		if err := s.Err(); err != nil {
			log.Panic(err)
		}

		close(out)
	}()

	return out
}
