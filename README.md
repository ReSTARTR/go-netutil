go-netutil
====

`netutil.RRDialer` has the same interface of `net.Dial`.
With this package, You can request to a specific IP address with priority as you like.

In golang, `net.LookupIP` returns IP addresses after `sortBy6724`, that is defined in [src/net/addrselect.go](https://golang.org/src/net/addrselect.go).
Because of this, address for `net.Dial` might connect to a fixed ip address in some case.
For example, source IP and destination IP are in the same local area network like "10.0.0.0/16" or "192.168.0.0/24".
In my case, in a AWS VPC network, http.Client always connect to a specific IP address of ELB from an EC2 instance.
This will ruin the effects of DNS round robin.

You can change the priority of each address with this package.

`netutil.DefaultRRDialer` sort a list of address by random priority.

Installation
----

```bash
go get github.com/ReSTARTR/go-netutil
```

Usage
----

### netutil.RRDialer

`netutil.DefaultRRDialer` uses `netutil.DefaultSort` for sorting IP addresses.
`netutil.DefaultSort` sorts IPs randomly.

```go
import (
	"net/http"
	"github.com/ReSTARTR/go-netutil"
)

func main() {
	req, _ := http.NewRequest("GET",  "https://example.com/foo/bar", nil)
	client := http.Client{
		Transport: &http.Transport {
			Dial: netutil.DefaultRRDialer.Dial,
		},
	}
	res, _ := client.Do(req)
	_ = res
}
```

You can set a specific `netutil.RRDialer.Sort` function.

```go
dialer := netutil.Dialer{
	Sort: func(ips []net.IP) ([]net.IP) {
		sort.Sort(ByFoo(ips))
		return ips
	},
}
client := http.Client{
	Transport: &http.Transport {
		Dial: dialer.Dial,
	},
}
client.Do(req)
```


Contribution
----

- Fork (https://github.com/ReSTARTR/go-netutil/fork)
- Create a feature branch
- Commit your changes
- Rebase your local changes against the master branch
- Run test suite with the make test command and confirm that it passes
- Create a new Pull Request

Author
----

[ReSTARTR](https://github.com/ReSTARTR)
