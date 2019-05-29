# rlog [![GoDoc][doc-img]][doc] [![Build Status][build-img]][build]

rlog logs incomming HTTP requests, and provides access to a logger from HTTP request's `context.Context`.

rlog uses [github.com/uber-go/zap](https://github.com/uber-go/zap) as a logger implementation.

Each log record has the key `requestID`.

## Installation

```bash
go get github.com/i-core/rlog
```

## Usage

```go
package main

import (
    "fmt"
    "net/http"
    "os"

    "go.uber.org/zap"
    "github.com/i-core/rlog"
)

func main() {
    // Create a logger.
    log, err := zap.NewProduction()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to create logger: %s", err)
        os.Exit(1)
    }

    // Create a logger middleware.
    logmw := rlog.NewMiddleware(log)

    // Wrap HTTP handler with the logger middleware.
    handler := logmw(newHandler())

    // Start a web server.
    log.Info("Server started")
    log.Fatal("Server finished", zap.Error(http.ListenAndServe(":8080", handler)))
}

func newHandler() http.Handler {
    return http.HandlerFunc(w http.ResponseWriter, r *http.Request) {
        // Get a request's logger.
        log := rlog.FromContext(r.Context())
        log.Info("Handle request")
    }
}
```

## Contributing

Thanks for your interest in contributing to this project.
Get started with our [Contributing Guide][contrib].

## License

The code in this project is licensed under [MIT license][license].

[build-img]: https://travis-ci.com/i-core/rlog.svg?branch=master
[build]: https://travis-ci.com/i-core/rlog

[doc-img]: https://godoc.org/github.com/i-core/rlog?status.svg
[doc]: https://godoc.org/github.com/i-core/rlog

[contrib]: https://github.com/i-core/.github/blob/master/CONTRIBUTING.md
[license]: LICENSE