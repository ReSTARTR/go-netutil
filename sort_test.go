package netutil

import (
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

func TestRandomSort(t *testing.T) {
	tests := []struct {
		ips    []net.IP
		repeat int
		want   float32
	}{
		{
			[]net.IP{
				net.IP{192, 168, 0, 3},
				net.IP{192, 168, 1, 31},
				net.IP{192, 168, 2, 193},
			},
			1000,
			0.3,
		},
		{
			[]net.IP{
				net.IP{10, 0, 1, 131},
				net.IP{10, 0, 2, 25},
				net.IP{10, 0, 3, 100},
				net.IP{10, 0, 4, 8},
			},
			1000,
			0.22,
		},
	}
	for _, test := range tests {
		results := map[string]int{}
		for i := 0; i < test.repeat; i += 1 {
			ipss := test.ips[:] // copy
			RandomSort(ipss)
			first := ipss[0].String()
			if _, ok := results[first]; !ok {
				results[first] = 0
			}
			results[first] += 1
		}
		for k, v := range results {
			p := float32(v) / float32(test.repeat) * 100.0
			if p < test.want {
				t.Errorf("%s is under 30%%: %f %", k, p)
			}
		}
	}
}
