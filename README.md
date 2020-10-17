<a href="https://alt4.dev"><img src="https://alt4.dev/banner.svg" alt="" height="120"></a>

## Golang Logging Client

**Go lang logging client for <a href="https://alt4.dev">Alt4.dev</a>**

### Install
```shell script
go get github.com/alt4dev/go
```

#### Authentication and Config
Alt4 automatically reads config information from the environment. Auth tokens can be generated from `Settings -> Manage Access -> Application Keys` on <a target="_blank" href="https://app.alt4.dev">app.alt4.dev</a>
 You can configure alt4 in 3 ways:
1. **Using a config file:** Pass in the path to the config file via env variable `ALT4_CONFIG`. A config file can be downloaded during token generation.

2. **Using Environment Variables:** The following options can be configured. If set they'll override options set using the config file.
    - `ALT4_AUTH_TOKEN` A token string used to authorize your request to alt4
    - `ALT4_MODE` This controls the amount of ouptut alt4 emits to `stderr`. The following modes can be used:
        - `release`(default) - Under this mode all logs will be sent to alt4 without logging the logs to `stderr`
        unless there's a connection/authentication error.
        - `debug` Under this mode logs will be sent to alt4 and written to `stderr`.
        - `testing` Under this mode logs will only be emitted to `stderr` and not sent to alt4. We use this mode to develop our products locally.
        - `json`(coming soon) - Under this mode, logs will be written to a JSON file which can later be uploaded to alt4
    - `ALT4_SINK` A string specifying a sink to log the logs under. By default, logs entries will be logged to the sink `default`.
3. **Set options from the code**
```go
package main
import (
    alt4Service "github.com/alt4dev/go/service"
)

alt4Service.SetAuthToken("YOUR TOKEN HERE")
alt4Service.SetMode("release")
alt4Service.SetSink("default")
```

### Usage
This client emulates golang's built in `log` package as much as possible. Logs will be written asynchronously to alt4.
If you're running on a system that doesn't allow background processes(goroutines) e.g. google cloud run,
we recommend using [log grouping](#grouping) and making sure you defer close group. This will wait for all writes to complete.
You can call the function `Result` after a log to wait for the writing to finish.
```go
package main
import (
    "github.com/alt4dev/go/log"
    "time"
)

func main() {
    log.Println("Normal logging as you're used to")
    log.Debugf("A formatted log entry, current time %s", time.Now())
    log.Warning("Create a log with a Warning severity level")
    log.Error("Create a log with an error severity level. This won't exit after.")
    log.Fatal("Logs with a critical severity level then exits with status 1.")
    log.Panicln("Logs with a critical severity level then panics.")
}
```

#### Claims
Claims are extra data that you want to relate to a log entry but aren't part of the log message.
This data can be used while filtering for logs from alt4. Claims implement all methods implemented by the `log` package.
```go
package main
import "github.com/alt4dev/go/log"

func main() {
    log.Claims{
        "user_id": "user triggering this entry",
    }.Println("A normal log message")

    log.Claims{
        "id": "some_id",
        "name": "Some name here",
    }.Warning("Just a warning")
}
```

#### Grouping
Grouping can help you resolve issues faster by grouping related logs together.
Alt4 groups logs based on if they're running from the same goroutine.
```go
package main
import (
    "fmt"
    "github.com/alt4dev/go/log"
)

func processUrl(url string) {
    // We want logs from each process grouped together even if they run in parallel
    defer log.Group(fmt.Sprintf("Processing %s", url)).Close()
    log.Println("This log and those after will be grouped under the current routine")

    // You can also create a group with claims
    defer log.Claims{"url": url}.Group("Processing ", url).Close()
}

func main() {
    urls := []string{
        "http://github.com",
        "http://google.com",
        "http://stackoverflow.com",
        "http://medium.com",
        "http://gitlab.com",
    }
    // Process these URLs in different processes.
    for _, url := range urls {
        go processUrl(url)
    }
}
```

#### Set Default Logger to Write to Alt4
This is the quickest way to get started with alt4 without importing the library in every file that you do log from.
This is the recommended path for a pre-existing code base without the intention to use claims in logs.

This achieved by providing the following writer interfaces:
- **`github.com/alt4dev/go/service.Writer`** This writer receives a log message and writes it asynchronously to alt4.
- **`github.com/alt4dev/go/service.SyncWriter`** This writer writes synchronously, i.e. blocks as the log is written.
This is only recommended in environments where background processes aren't allowed or are discouraged, e.g. google cloud run.

The example below demonstrates the usage of both writers:
```go
package main
import (
    alt4Service "github.com/alt4dev/go/service"
    "log"
)

func main() {
    // Set default logger to use alt4
    log.SetOutput(alt4Service.Writer)
    log.Println("This writes to alt4")
    
    // Set a custom logger to use alt4
    logger := log.New(alt4Service.SyncWriter, "[my custom logger]", 0)
    logger.Println("This writes synchronously to alt4")
}
```
