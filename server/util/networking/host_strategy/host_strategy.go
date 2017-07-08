// Package hoststrategy Allows for getting the next host out of a pool based on the current strategy
package hoststrategy

import (
	"net"
	"net/url"
	"time"
)

var (
	//Dialer is a default dialer with a much shorter timeout so that we
	//can detect very slow hosts
	//Hosts that are no longer slow will connect the next round
	dialer = net.Dialer{
		Timeout: time.Second * 10,
	}
)

//HostStrategy supports a load balancing strategy
type HostStrategy interface {
	GetNextHost() *url.URL
}

//hostTCPDialer returns true if a connection was established
func hostTCPDialer(h string) bool {
	conn, err := dialer.Dial("tcp", h)
	if err != nil {
		return false
	}
	conn.Close()

	return true
}
