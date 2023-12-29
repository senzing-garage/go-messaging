# go-messaging

If you are beginning your journey with
[Senzing](https://senzing.com/),
please start with
[Senzing Quick Start guides](https://docs.senzing.com/quickstart/).

You are in the
[Senzing Garage](https://github.com/senzing-garage)
where projects are "tinkered" on.
Although this GitHub repository may help you understand an approach to using Senzing,
it's not considered to be "production ready" and is not considered to be part of the Senzing product.
Heck, it may not even be appropriate for your application of Senzing!

## :warning: WARNING: go-messaging is still in development :warning: _

At the moment, this is "work-in-progress" with Semantic Versions of `0.n.x`.
Although it can be reviewed and commented on,
the recommendation is not to use it yet.

## Synopsis

The Senzing `go-messaging` packages are used to create structured messages.

[![Go Reference](https://pkg.go.dev/badge/github.com/senzing-garage/go-messaging.svg)](https://pkg.go.dev/github.com/senzing-garage/go-messaging)
[![Go Report Card](https://goreportcard.com/badge/github.com/senzing-garage/go-messaging)](https://goreportcard.com/report/github.com/senzing-garage/go-messaging)
[![License](https://img.shields.io/badge/License-Apache2-brightgreen.svg)](https://github.com/senzing-garage/go-messaging/blob/main/LICENSE)

[![gosec.yaml](https://github.com/senzing-garage/go-messaging/actions/workflows/gosec.yaml/badge.svg)](https://github.com/senzing-garage/go-messaging/actions/workflows/gosec.yaml)
[![go-test-linux.yaml](https://github.com/senzing-garage/go-messaging/actions/workflows/go-test-linux.yaml/badge.svg)](https://github.com/senzing-garage/go-messaging/actions/workflows/go-test-linux.yaml)
[![go-test-darwin.yaml](https://github.com/senzing-garage/go-messaging/actions/workflows/go-test-darwin.yaml/badge.svg)](https://github.com/senzing-garage/go-messaging/actions/workflows/go-test-darwin.yaml)
[![go-test-windows.yaml](https://github.com/senzing-garage/go-messaging/actions/workflows/go-test-windows.yaml/badge.svg)](https://github.com/senzing-garage/go-messaging/actions/workflows/go-test-windows.yaml)

## Overview

`go-messaging` generates structured messages in multiple formats.
Currently, the JSON format and an
[slog](https://pkg.go.dev/golang.org/x/exp/slog)-friendly format are supported.

## Use

```go
import "github.com/senzing-garage/go-messaging/messenger"

aMessenger, _ := messenger.New()
fmt.Println(aMessenger.NewJson(0001, "Bob", "Mary"))
fmt.Println(aMessenger.NewSlog(0001, "Bob", "Mary"))
```

Output:

```console
{"time":"YYYY-MM-DDThh:mm:ss.nnnnnnnnn-00:00","level":"TRACE","id":"senzing-99990001","details":{"1":"Bob","2":"Mary"}}
[id senzing-99990001 details map[1:Bob 2:Mary]]
```

For more examples, see
[main.go](main.go)

## References

- [API documentation](https://pkg.go.dev/github.com/senzing-garage/go-messaging)
- [Development](docs/development.md)
- [Errors](docs/errors.md)
- [Examples](docs/examples.md)
