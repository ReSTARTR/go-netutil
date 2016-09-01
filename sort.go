package dialx

import (
	"math/rand"
	"net"
	"sort"
	"time"
)

// RandomSort sorts ips at random.
func RandomSort(ips []net.IP) []net.IP {
	sort.Sort(atRandom(ips))
	return ips
}

type atRandom []net.IP

func (b atRandom) Len() int      { return len(b) }
func (b atRandom) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b atRandom) Less(_, _ int) bool {
	rand.Seed(time.Now().UnixNano())
	l := b.Len()
	return rand.Intn(l) > rand.Intn(l)
}
