// Package requestutil is helper package for HTTP request.
//
package requestutil

import (
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var pickupIPFunc func(host string, fn Selector) (ip string, err error) = nil

type Selector func(addrs []net.IP) net.IP

// DefaultSelector returns an addr randomly.
func DefaultSelector(addrs []net.IP) net.IP {
	as := []net.IP{}
	for _, addr := range addrs {
		ip4 := addr.To4()
		if ip4 != nil {
			log.Printf("%+v\n", addr.To4().String())
			as = append(as, ip4)
		}
	}
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(as))
	return as[i]
}

// NewRequestWithSelector returns a request object.
// This function switches request.Host into an IP address that
// is looked up from DNS service and is picked by a selector function.
//
// In go, list of addresses that is looked up from DNS serivce are
// sorted by sortByRFC6724 function.
// (https://golang.org/src/net/addrselect.go#L13)
//
// Because of this, if both of source and destination exist in the same local network,
// ip address for requesting might be fixed to an certain address
// no matter whether DNS uses round robin or not.
//
// So, this function resolve the sort problem by switcing hostname for http request to an
// address that is picked by an specific function of selection.
//
// NOTE: request of HTTPS might fail caused by its certification method.
//
// Call without selector function(use DefaultSelector)
//
//   req, err := NewRequestWithSelector("GET", "http://example.com/path/to/page", nil, nil)
//
// or set a selector function in last argument.
//
//   selector := func(addrs []net.IP) net.IP {
//   	// always select first addr
//   	return addrs[0]
//   }
//   req, err := NewRequestWithSelector("GET", "http://example.com/path/to/page", nil, selector)
//
// then, values are
//
//   req.Host     // hostname
//   req.URL.Host // ip address
func NewRequestWithSelector(method, urlStr string, body io.Reader, fn Selector) (*http.Request, error) {
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}

	var addr string
	if pickupIPFunc == nil {
		addr, err = defaultPickupIP(req.URL.Host, fn)
	} else {
		addr, err = pickupIPFunc(req.URL.Host, fn)
	}
	if err != nil {
		return nil, err
	}

	hp := strings.Split(":", req.URL.Host)
	host := addr
	if len(hp) == 2 {
		host += ":" + hp[1]
	}

	// set original Host header,
	// and switch host to ipaddr
	req.Host = req.URL.Host
	req.URL.Host = host
	return req, nil
}

func defaultPickupIP(host string, fn Selector) (ip string, err error) {
	addrs, err := net.LookupIP(host)
	if err != nil {
		return "", err
	}
	if fn == nil {
		fn = DefaultSelector
	}
	return fn(addrs).To4().String(), nil
}

type Client struct {
	Selector Selector
	PickupIP func(host string, fn Selector) (ip string, err error)
}

var DefaultClient = Client{}

func transport(proxyURL *url.URL) *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
}

func Get(urlStr string) (*http.Response, error) {
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	return DefaultClient.Do(req)
}

func Head(url string) (*http.Response, error) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return nil, err
	}
	return DefaultClient.Do(req)
}

func Post(url string, bodyType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", bodyType)
	return DefaultClient.Do(req)
}

func PostForm(url string, data url.Values) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return DefaultClient.Do(req)
}

func (c *Client) Do(r *http.Request) (*http.Response, error) {
	selector := c.Selector
	if selector == nil {
		selector = DefaultSelector
	}

	pickup := c.PickupIP
	if pickup == nil {
		pickup = defaultPickupIP
	}

	addr, err := pickup(r.URL.Host, selector)
	if err != nil {
		return nil, err
	}

	host := addr
	hp := strings.Split(r.URL.Host, ":")
	if len(hp) == 2 {
		host += hp[1]
	}

	pu, err := url.Parse(r.URL.Scheme + "://" + host)
	client := &http.Client{
		Transport: transport(pu),
	}
	return client.Do(r)
}
