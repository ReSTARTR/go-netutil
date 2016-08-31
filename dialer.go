package dialx

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"sort"
	"time"
)

// Dialer contains options for connection to an address with a specific IP.
//
// TODO: connection pool
type Dialer struct {
	dialer *net.Dialer
	Sort   func([]net.IP) []net.IP
}

// DefaultDialer
var DefaultDialer = &Dialer{
	dialer: &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	},
}

// RandomSort sorts ips at random.
func RandomSort(ips []net.IP) []net.IP {
	sort.Sort(atRandom(ips))
	return ips
}

// Dial connects to the address with an specific IP
// on the named network.
func (d *Dialer) Dial(network, address string) (net.Conn, error) {
	return d.DialContext(context.Background(), network, address)
}

// DialContext connects to the address with an specific IP
// on the named network using the provided context.
func (d *Dialer) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}

	addrs, err := net.LookupIP(host)
	if err != nil {
		return nil, err
	}

	sort := d.Sort
	if sort == nil {
		sort = RandomSort
	}

	for _, addr := range sort(addrs) {
		ip4 := addr.To4()
		if ip4 == nil {
			continue
		}

		ipStr := ip4.String()

		name := net.JoinHostPort(ipStr, port)
		pc, err := d.dialer.Dial(network, name)
		if err == nil {
			return pc, nil
		}
	}
	return nil, fmt.Errorf("cannot get connection: %s://%s", network, addr)
}

type atRandom []net.IP

func (b atRandom) Len() int      { return len(b) }
func (b atRandom) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b atRandom) Less(_, _ int) bool {
	rand.Seed(time.Now().UnixNano())
	l := b.Len()
	return rand.Intn(l) > rand.Intn(l)
}
