# Contextual Logger

[![codecov](https://codecov.io/gh/godepo/logger/graph/badge.svg?token=psNjauhUqW)](https://codecov.io/gh/godepo/logger)
[![Go Report Card](https://goreportcard.com/badge/godepo/logger)](https://goreportcard.com/report/godepo/logger)
[![License](https://img.shields.io/badge/License-MIT%202.0-blue.svg)](https://github.com/godepo/elephant/blob/main/LICENSE)


This library solves the problem of contextual logging in business code built on components or layers. Most often, 
we want to get traceable log entries linked within each unique query. Achieving this by standard means complicates the 
code and increases possible errors.

The library allows you to use the standard slog library and place one of the variants of the logger object in a context 
from the go standard library. In addition, the library provides a type-safe interface to prevent untyped values for each 
structural tag from being passed to the log entry.
 
## Quick start

1. Configure your preferred logger in slog package with adapter.
2. Register your logger in slog.
3. Use this workaround library for logging your entries.

Create root logger:

```go
rootLogger := logger.From(context.Background())
```

Create logger for entire service:

```
serviceLogger := rootLogger.With(
    slog.String("service", "service_name), 
    slog.String("pod", "pod_name"), 
    slog.String("node", "node_address"),
)
```

Write middleware for http endpoint: 

```go
func LoggerMiddleware(next http.Handler, log logger.Logger, routeName string) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Add the logger to the context using our special Logger
		// and add to this logger contextual tags
		ctx := log.With(r.Context(), slog.String("route", routeName), slog.String("endpoint", "http"))
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

Use this middleware like this:

```go
func main() {
    log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	slog.SetDefault(log)
	
    serviceLogger := logger.From(context.Background()).With(
        slog.String("service", "service_name), 
        slog.String("pod", "pod_name"),
        slog.String("node", "node_address"),
    )
	
    loginHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Fetch the Logger from Context
        ctx := r.Context()
        // Grab the Logger from the context
        logger := ctx.Value(LoggerCtxKey).(*slog.Logger)
        logger.ErrorContext(ctx, "failed to login")

        w.WriteHeader(http.StatusOK)
    })

    // Wrap your handler with the LoggerMiddleware
    http.Handle("/login", LoggerMiddleware(loginHandler, serviceLogger, "login"))
}
```

## Enhance Developer Experience

The constant use of basic slog attributes clutters up your business code a bit. Therefore, it is recommended to create 
your own module with declarative tags to specify your structural tags.

For example:

```go
package logtag

import (
	"log/slog"

	"github.com/google/uuid"
)

func ServiceName(name string) slog.Attr {
	return slog.String("service", name)
}

func PodName(name string) slog.Attr {
	return slog.String("pod", name)
}

func UserID(userID uuid.UUID) slog.Attr {
	return slog.String("user_id", userID.String())
}

func Amount(amount uint64) slog.Attr {
	return slog.Uint64("amount", amount)
}

```

This is help write you log entries like this:

```go
logging.Warn(ctx, "possible fraud detected", logtag.Amount(req.Amount), logtag.UserID(req.UserID))
```

This style is more expressive and cleaner, allowing you to bring some of the complexity to the DSL syntax level.