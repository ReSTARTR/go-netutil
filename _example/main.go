package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"

	"github.com/ReSTARTR/go-netutil"
	"github.com/tcnksm/go-httptraceutils"
)

var url = flag.String("url", "https://www.google.com/", "")

func main() {
	flag.Parse()

	req, _ := http.NewRequest("GET", *url, nil)
	ctx := httptraceutils.WithClientTrace(context.Background())
	req = req.WithContext(ctx)

	c := http.Client{
		Transport: &http.Transport{
			Proxy: nil,
			Dial:  netutil.DefaultRRDialer.Dial,
		},
	}

	for i := 0; i < 10; i++ {
		res, err := c.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		_ = res
	}
}
