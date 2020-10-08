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
    - `ALT4_MODE` This is either `release` or `debug` with the default set to `release`. Under `debug` mode, a log entry will be written to stderr/(specified emit file) in addition to being sent to alt4.
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
This client emulates golang's built in `log` package as much as possible.
The example below demonstrates this usage:
```go
package main

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
