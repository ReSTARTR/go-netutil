package main

import (
	"context"
	_ "crypto/tls"
	_ "crypto/x509"
	"flag"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/ReSTARTR/go-http-requestutil"
	"github.com/tcnksm/go-httptraceutils"
)

var u = flag.String("url", "https://www.google.com/", "URL string")

func main() {
	flag.Parse()

	ctx := httptraceutils.WithClientTrace(context.Background())

	// selector := requestutil.DefaultSelector
	url := *u
	req, _ := http.NewRequest("GET", url, nil)

	req = req.WithContext(ctx)

	c := requestutil.Client{}
	res, _ := c.Do(req)
	// res, _ := requestutil.RequestViaProxy(req)
	/*
		req, err := requestutil.NewRequestWithSelector("GET", url, nil, selector)
		if err != nil {
			log.Fatal(err)
		}
	*/

	b, err := httputil.DumpRequest(req, false)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Header:\n%s", b)

	/*
		var client *http.Client
		if req.URL.Scheme == "https" {
			client = &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{RootCAs: x509.NewCertPool()},
				},
			}
		} else {
			client = http.DefaultClient
		}

		res, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
	*/
	_ = res
}
