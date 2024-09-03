# go-messaging

If you are beginning your journey with [Senzing],
please start with [Senzing Quick Start guides].

You are in the [Senzing Garage] where projects are "tinkered" on.
Although this GitHub repository may help you understand an approach to using Senzing,
it's not considered to be "production ready" and is not considered to be part of the Senzing product.
Heck, it may not even be appropriate for your application of Senzing!

## :warning: WARNING: go-messaging is still in development :warning: _

At the moment, this is "work-in-progress" with Semantic Versions of `0.n.x`.
Although it can be reviewed and commented on,
the recommendation is not to use it yet.

## Synopsis

The Senzing `go-messaging` packages are used to create structured messages.

[![Go Reference Badge]][Package reference]
[![Go Report Card Badge]][Go Report Card]
[![License Badge]][License]
[![go-test-linux.yaml Badge]][go-test-linux.yaml]
[![go-test-darwin.yaml Badge]][go-test-darwin.yaml]
[![go-test-windows.yaml Badge]][go-test-windows.yaml]

[![golangci-lint.yaml Badge]][golangci-lint.yaml]

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

1. [API documentation]
1. [Development]
1. [Errors]
1. [Examples]
1. [Package reference]
1. Related artifacts:
    1. [DockerHub]
1. [JSON TypeDef](https://jsontypedef.com/)

[API documentation]: https://pkg.go.dev/github.com/senzing-garage/go-messaging
[Development]: docs/development.md
[DockerHub]: https://hub.docker.com/r/senzing/go-messaging
[Errors]: docs/errors.md
[Examples]: docs/examples.md
[Go Reference Badge]: https://pkg.go.dev/badge/github.com/senzing-garage/go-messaging.svg
[Go Report Card Badge]: https://goreportcard.com/badge/github.com/senzing-garage/go-messaging
[Go Report Card]: https://goreportcard.com/report/github.com/senzing-garage/go-messaging
[go-test-darwin.yaml Badge]: https://github.com/senzing-garage/go-messaging/actions/workflows/go-test-darwin.yaml/badge.svg
[go-test-darwin.yaml]: https://github.com/senzing-garage/go-messaging/actions/workflows/go-test-darwin.yaml
[go-test-linux.yaml Badge]: https://github.com/senzing-garage/go-messaging/actions/workflows/go-test-linux.yaml/badge.svg
[go-test-linux.yaml]: https://github.com/senzing-garage/go-messaging/actions/workflows/go-test-linux.yaml
[go-test-windows.yaml Badge]: https://github.com/senzing-garage/go-messaging/actions/workflows/go-test-windows.yaml/badge.svg
[go-test-windows.yaml]: https://github.com/senzing-garage/go-messaging/actions/workflows/go-test-windows.yaml
[golangci-lint.yaml Badge]: https://github.com/senzing-garage/go-messaging/actions/workflows/golangci-lint.yaml/badge.svg
[golangci-lint.yaml]: https://github.com/senzing-garage/go-messaging/actions/workflows/golangci-lint.yaml
[License Badge]: https://img.shields.io/badge/License-Apache2-brightgreen.svg
[License]: https://github.com/senzing-garage/go-messaging/blob/main/LICENSE
[Package reference]: https://pkg.go.dev/github.com/senzing-garage/go-messaging
[Senzing Garage]: https://github.com/senzing-garage
[Senzing Quick Start guides]: https://docs.senzing.com/quickstart/
[Senzing]: <https://senzing.com/-> [JSON TypeDef on Github](https://github.com/jsontypedef)
