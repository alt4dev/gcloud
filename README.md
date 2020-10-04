<a href="https://alt4.dev"><img src="https://alt4.dev/banner.svg" alt="" height="120"></a>

## Golang Logging Client

**Go lang logging client for <a href="https://alt4.dev">Alt4.dev</a>**

### Install
```shell script
go get github.com/alt4dev/go
```

### Usage
#### Authentication and Config
Alt4 automatically reads config information from the environment. You can configure alt4 in 3 ways:
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
 "os"
)

alt4Service.SetAuthToken("YOUR TOKEN HERE")
alt4Service.SetMode("release")
alt4Service.SetDebugOutput(os.Stderr)
```