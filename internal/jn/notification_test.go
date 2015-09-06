// "THE BEER-WARE LICENSE" (Revision 42):
// <tobias.rehbein@web.de> wrote this file. As long as you retain this notice
// you can do whatever you want with this stuff. If we meet some day, and you
// think this stuff is worth it, you can buy me a beer in return.
//                                                             Tobias Rehbein

package jn

import (
	"testing"
)

const threemaMonitorString = `method call sender=:1.22 -> dest=org.freedesktop.Notifications serial=136 path=/org/freedesktop/Notifications; interface=org.freedesktop.Notifications; member=Notify
   string "AndroidNotification"
   uint32 1664
   string "/data/notificationIcon/ch.threema.app2130838684.png"
   string "Herp Derp"
   string "3 neue Nachrichten"
   array [
      string "default"
      string ""
   ]
   array [
      dict entry(
         string "x-nemo-preview-icon"
         variant             string "/data/notificationIcon/ch.threema.app2130838684.png"
      )
      dict entry(
         string "x-nemo-preview-body"
         variant             string "3 neue Nachrichten"
      )
      dict entry(
         string "x-nemo-preview-summary"
         variant             string "Herp Derp"
      )
   ]
   int32 -1`

var threemaNotification = &Notification{
	Body:    "3 neue Nachrichten",
	Summary: "Herp Derp",
}

const clockMonitorString = `method call sender=:1.117 -> dest=org.freedesktop.Notifications serial=43 path=/org/freedesktop/Notifications; interface=org.freedesktop.Notifications; member=Notify
   string ""
   uint32 0
   string ""
   string ""
   string ""
   array [
   ]   
   array [
      dict entry(
         string "x-nemo-owner"
         variant             string "Uhr"
      )   
      dict entry(
         string "x-nemo-preview-body"
         variant             string "Verbleibende Zeit: 17 Stunden und 18 Minuten"
      )   
      dict entry(
         string "category"
         variant             string "x-jolla.settings.clock"
      )   
   ]   
   int32 -1`

var clockNotification = &Notification{
	Body:    "Verbleibende Zeit: 17 Stunden und 18 Minuten",
	Summary: "Uhr",
}

const commhistorydMonitorString = `method call sender=:1.36 -> dest=org.freedesktop.Notifications serial=117 path=/org/freedesktop/Notifications; interface=org.freedesktop.Notifications; member=Notify
   string "Nachrichten"
   uint32 0
   string ""
   string ""
   string ""
   array [
      string "default"
      string ""
      string "app"
      string ""
   ]
   array [
      dict entry(
         string "x-nemo-preview-summary"
         variant             string "Herp Derp"
      )
      dict entry(
         string "x-nemo-preview-body"
         variant             string "Test"
      )
      dict entry(
         string "x-nemo-owner"
         variant             string "commhistoryd"
      )
      dict entry(
         string "category"
         variant             string "x-nemo.messaging.sms"
      )
      dict entry(
         string "x-nemo-remote-action-app"
         variant             string "org.nemomobile.qmlmessages / org.nemomobile.qmlmessages showGroupsWindow"
      )
      dict entry(
         string "x-nemo-timestamp"
         variant             string "2015-09-03T10:50:33Z"
      )
      dict entry(
         string "x-commhistoryd-data"
         variant             array of bytes "DEADBEEF"
      )
      dict entry(
         string "x-nemo-remote-action-default"
         variant             string "org.nemomobile.qmlmessages / org.nemomobile.qmlmessages startConversation Some Stuff"
      )
   ]
   int32 -1`

var commhistorydNotification = &Notification{
	Body:    "Test",
	Summary: "Herp Derp",
}

func (na *Notification) Equals(nb *Notification) bool {
	if na.Body != nb.Body {
		return false
	}
	if na.Summary != nb.Summary {
		return false
	}
	return true
}

func testNewNotificationFromMonitorString(t *testing.T, ms string, n *Notification) {
	nn, err := NewNotificationFromMonitorString(ms)
	if err != nil {
		t.Error(err)
	}
	if !nn.Equals(n) {
		t.Errorf("%#v != %#v", nn, n)
	}
}

func TestNotificationThreema(t *testing.T) {
	testNewNotificationFromMonitorString(t, threemaMonitorString, threemaNotification)
}

func TestNotificationClock(t *testing.T) {
	testNewNotificationFromMonitorString(t, clockMonitorString, clockNotification)
}

func TestCommhistorydClock(t *testing.T) {
	testNewNotificationFromMonitorString(t, commhistorydMonitorString, commhistorydNotification)
}

var testEmptyTable = []struct {
	n        *Notification
	expected bool
}{
	{&Notification{"NonEmptySummary", "NonEmptyBody"}, false},
	{&Notification{"", "NonEmptyBody"}, false},
	{&Notification{"NonEmptySummary", ""}, true},
	{&Notification{"", ""}, true},
}

func TestEmpty(t *testing.T) {
	for _, tt := range testEmptyTable {
		if tt.n.IsEmpty() != tt.expected {
			t.Errorf("%#v.IsEmpty() != %v", tt.n, tt.expected)
		}
	}
}
