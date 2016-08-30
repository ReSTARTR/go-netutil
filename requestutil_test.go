package requestutil

import (
	"log"
	"net"
	"testing"
)

var (
	addrs []net.IP
)

type testcase struct {
	name  string
	addrs []net.IP
}

func TestNewRequestWithSelector(t *testing.T) {
	// setUp: avoid to call DNS service
	pickupIPFunc = func(host string, _ Selector) (ip string, err error) {
		return DefaultSelector(addrs).String(), nil
	}

	tests := []testcase{
		testcase{
			"case 01",
			[]net.IP{
				net.IP{192, 168, 0, 3},
			},
		},
		testcase{
			"case 02",
			[]net.IP{
				net.IP{192, 168, 0, 3},
				net.IP{192, 168, 1, 31},
				net.IP{192, 168, 2, 193},
			},
		},
	}
	for _, test := range tests {
		addrs = test.addrs
		i := 0
		picks := map[string]int{}
		for i < 100 {
			i += 1
			req, err := NewRequestWithSelector("GET", "http://example.com/ping", nil, nil)
			if err != nil {
				t.Errorf("%s: %s", err)
				continue
			}
			if req.Host != "example.com" {
				t.Errorf("%s: req.Host is not hostname")
				continue
			}
			if req.URL.Host == "example.com" {
				t.Errorf("%s: req.URL.Host is not ip address")
				continue
			}
			if _, ok := picks[req.URL.Host]; !ok {
				picks[req.URL.Host] = 0
			}
			picks[req.URL.Host] += 1
		}

		if len(picks) != len(test.addrs) {
			t.Errorf("%s: There are some of addresses that is not picked up")
			continue
		}

		log.Println("Test:", test.name)
		for k, v := range picks {
			log.Println("", k, v)
		}
	}
	// tearDown
	pickupIPFunc = nil
}
