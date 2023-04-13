# go-messaging

## :warning: WARNING: go-messaging is still in development :warning: _

At the moment, this is "work-in-progress" with Semantic Versions of `0.n.x`.
Although it can be reviewed and commented on,
the recommendation is not to use it yet.

## Synopsis

The Senzing `go-messaging` packages are used to create structured messages.

[![Go Reference](https://pkg.go.dev/badge/github.com/senzing/go-messaging.svg)](https://pkg.go.dev/github.com/senzing/go-messaging)
[![Go Report Card](https://goreportcard.com/badge/github.com/senzing/go-messaging)](https://goreportcard.com/report/github.com/senzing/go-messaging)
[![go-test.yaml](https://github.com/Senzing/go-messaging/actions/workflows/go-test.yaml/badge.svg)](https://github.com/Senzing/go-messaging/actions/workflows/go-test.yaml)
[![License](https://img.shields.io/badge/License-Apache2-brightgreen.svg)](https://github.com/Senzing/go-messaging/blob/main/LICENSE)

## Overview

`go-messaging` generates structured messages in multiple formats.
Currently, the JSON format and an
[slog](https://pkg.go.dev/golang.org/x/exp/slog)-friendly format are supported.

## Use

```go
import "github.com/senzing/go-messaging/messenger"

aMessenger, _ := messenger.New()
fmt.Println(aMessenger.NewJson(0001, "Bob", "Mary"))
fmt.Println(aMessenger.NewSlog(0001, "Bob", "Mary"))
```

Output

```console
{"time":"YYYY-MM-DDThh:mm:ss.nnnnnnnnn-00:00","level":"TRACE","id":"senzing-99990001","details":{"1":"Bob","2":"Mary"}}
[id senzing-99990001 details map[1:Bob 2:Mary]]
```

## References

- [API documentation](https://pkg.go.dev/github.com/senzing/go-messaging)
- [Development](docs/development.md)
- [Errors](docs/errors.md)
- [Examples](docs/examples.md)
