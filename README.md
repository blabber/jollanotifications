jollanotifications
==================
[![Build Status](https://travis-ci.org/blabber/jollanotifications.svg?branch=master)](https://travis-ci.org/blabber/mbox)
[![GoDoc](https://godoc.org/github.com/blabber/mbox?status.svg)](https://godoc.org/github.com/blabber/mbox)

A hackish solution to access my Jolla phone's notification via a web interface.

TODO
----

 * Replace the dbus-monitor command with a go package providing access to D-Bus
 * Make the whole thing configurable (DO NOT BIND TO ALL INTERFACES)
 * Some kind of logging
 * Do not fail silently in web interface if a refresh of the view model fails
