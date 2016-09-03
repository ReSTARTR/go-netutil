package netutil

import (
	"fmt"
	"log"
	"net"
	"time"

	"golang.org/x/net/context"
)

// Dialer contains options for connection to an address with a specific IP.
//
// TODO: connection pool
type RRDialer struct {
	dialer *net.Dialer
	Sort   func([]net.IP)
}

// DefaultDialer has a dialer that includes the same as
// DefaultTransport use
var DefaultRRDialer = &RRDialer{
	dialer: &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	},
}

// Dial connects to the address with an specific IP
// on the named network.
func (d *RRDialer) Dial(network, address string) (net.Conn, error) {
	return d.DialContext(context.Background(), network, address)
}

// DialContext connects to the address with an specific IP
// on the named network using the provided context.
func (d *RRDialer) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}

	addrs, err := net.LookupIP(host)
	if err != nil {
		return nil, err
	}

	if d.Sort == nil {
		RandomSort(addrs)
	} else {
		d.Sort(addrs)
	}

	for _, addr := range addrs {
		ip4 := addr.To4()
		if ip4 == nil {
			continue
		}

		ipStr := ip4.String()

		name := net.JoinHostPort(ipStr, port)
		pc, err := d.dialer.Dial(network, name)
		if err != nil {
			// retry connect with a next address
			log.Printf("failed to net.Dial(%s, %s): %s", network, name, err)
			continue
		}

		return pc, nil
	}
	return nil, fmt.Errorf("cannot get connection: %s://%s", network, addr)
}
