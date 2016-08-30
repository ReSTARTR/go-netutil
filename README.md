go-http-requestutil
====

`requestutil.Client` has the same interface of `http.Client`.
With this package, You can request to a specific IP address with priority as you like.

`net.LookupIP` returns IP addresses after `sortBy6724`.
Because of this, address for `net.Dial` might be fixed always in some case.
(For example, source IP and destination IP are in the same local area network.)
You can change the priority of address selection by this package.

`requestutil.DefaultClient` select an address by random.

Feature
----

- TBD

Installation
----

```bash
go get github.com/ReSTARTR/go-http-requestutil
```

Usage
----

```go
import (
	"net/http"
	"github.com/ReSTARTR/go-http-requestutil
)

func main() {
	req, _ := http.NewRequest("GET",  "https://example.com/foo/bar", nil)
	client := requestutil.DefaultClient
	res, _ := client.Do(req)
	_ = res
}
```

You can set a specific `request.Selector` function.

```go
client := requestutil.Client {
	Selector: func(ips []net.IP) net.IP {
		return ips[0] // return the first address always.
	}
}
client.Do(req)
```


Contribution
----

- Fork (https://github.com/ReSTARTR/go-http-requestutil/fork)
- Create a feature branch
- Commit your changes
- Rebase your local changes against the master branch
- Run test suite with the make test command and confirm that it passes
- Create a new Pull Request

Author
----

[ReSTARTR](https://github.com/ReSTARTR)
