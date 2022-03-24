<a href="https://alt4.dev"><img src="https://alt4.dev/banner.svg" alt="" height="120"></a>


## Golang Logging Client, [Docs](https://pkg.go.dev/mod/github.com/alt4dev/gcloud)

### Install
```shell script
go get github.com/alt4dev/gcloud
```

#### Authentication and Config
To authenticate with gcloud see: [https://cloud.google.com/logging/docs/reference/libraries#setting_up_authentication](https://cloud.google.com/logging/docs/reference/libraries#setting_up_authentication).
An environment variable `PROJECT_ID` is required to set up the logging client.

Different modes are applicable:
- `release`: This is the default mode and, it sends logs to Google cloud logging without emitting them to `stderr`
- `debug`: This mode sends logs to Google cloud logging and emits them to `stderr`
- `testing`: This mode only emits logs to `stderr` without attempting to send them to Google cloud logging
- `silent`: This mode will silently ignore all the logs, they won't be sent to Google cloud logging or emitted to `stderr`

You can set a mode as shown below or by setting the environment variable `ALT4_MODE`
```go
package main
import (
    alt4Service "github.com/alt4dev/gcloud/service"
)

alt4Service.SetMode("release")
```


### Usage
This client emulates golang's built in `log` package as much as possible. Logs will be written asynchronously to Google cloud logging.
If you're running on a system that doesn't allow background processes(goroutines) e.g. google cloud run,
we recommend using [log grouping](#grouping) and making sure you defer close group. This will wait for all writes to complete.
```go
package main
import (
    "github.com/alt4dev/gcloud/log"
    "time"
)

func main() {
    log.Println("Normal logging as you're used to")
    log.Debugf("A formatted log entry, current time %s", time.Now())
    log.Warning("Create a log with a Warning severity level")
    log.Error("Create a log with an error severity level. This won't exit after.")
    log.Fatal("Logs with a critical severity level then exits with status 1.")
    log.Panic("Logs with a critical severity level then panics.")
}
```

#### Labels
Labels are extra data that you want to relate to a log entry but aren't part of the log message.
This data can be used while filtering for logs from alt4. Claims implement all methods implemented by the `log` package.
```go
package main
import "github.com/alt4dev/gcloud/log"

func main() {
    log.Labels{
        "user_id": "user triggering this entry",
    }.Println("A normal log message")

    log.Labels{
        "id": "some_id",
        "name": "Some name here",
    }.Warning("Just a warning")
    
    // You can reuse labels for multiple logs
    labels := log.Labels{
        "org": "Test Org",
    }
    
    labels.Debug("Operation one")
    labels.Error("The Op failed")
}
```

#### Grouping
Grouping can help you resolve issues faster by grouping related logs together.
Google cloud logging groups logs based on a http request, this library uses goroutines to determine logs in the same request.

This example demonstrates the setup of a logging middle ware.
```go
package main

import (
    "github.com/alt4dev/gcloud/log"
    "net/http"
)

func loggingMiddleWare(next http.Handler) http.Handler {
    return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
        // Start a logging group for this request, Calling `SetRequest` is not necessary
        defer log.Group(request.URL.String(), request.Method).SetRequest(request).Close()

        // All logs after this line will be grouped

        next.ServeHTTP(writer, request)
    })
}

```

#### Set Default Logger to Write to Google cloud logging
This is the quickest way to get started with Google cloud logging without importing the library in every file that you do log from.
This is the recommended path for a pre-existing code base without the intention to use claims in logs.

All logs will have the default level.

This achieved by providing the following writer interfaces:
- **`github.com/alt4dev/gcloud/service.Writer`** This writer receives a log message and writes it Google cloud logging.

The example below demonstrates this:
```go
package main
import (
    alt4Service "github.com/alt4dev/gcloud/service"
    "log"
)

func main() {
    // Set default logger to use Google cloud logging
    log.SetOutput(alt4Service.Writer)
    log.Println("This writes to Google cloud logging")
    
    // Set a custom logger to use Google cloud logging
    logger := log.New(alt4Service.Writer, "[my custom logger]", 0)
    logger.Println("This writes Google cloud logging")
}
```

