<a href="https://alt4.dev"><img src="https://alt4.dev/banner.svg" alt="" height="120"></a>

# Advanced Logging and Tracking 4 Developers
We stand by the principle of making developers' work easier, being able to debug your application easily is one of those requirements.

If you write logs through Google Cloud AppEngine python then you get your log requests automatically grouped for you per request by Google.
This functionality is not available on modern container based servers even the ones running on GCP. If you need the functionality you have to implement it yourself.

This library bridges that gap and provides a familiar API allowing you to easily group your logs. If you need assistance integrating this library reach out to: [me@billcountry.tech](mailto:me@billcountry.tech)

## Google Cloud Logging Client, [Docs](https://pkg.go.dev/mod/github.com/alt4dev/gcloud)

### Install
```shell script
go get github.com/alt4dev/gcloud
```

### Configuration
#### Authentication
To authenticate with gcloud see: [https://cloud.google.com/logging/docs/reference/libraries#setting_up_authentication](https://cloud.google.com/logging/docs/reference/libraries#setting_up_authentication).
An environment variable `PROJECT_ID` is required to set up the logging client.

#### Mode
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

#### Monitored Resource
Set the monitored resource that the logs will be produced from. Below is an example of a monitored resource for a Google cloud run service.
```json
{
    "labels": {
        "configuration_name": "test",
        "location": "us-central1",
        "project_id": "alt4dev",
        "revision_name": "test-april-fools",
        "service_name": "test-service"
    },
    "type": "cloud_run_revision"
}
```

This library provides a function `service.SetMonitoredResource` to set a monitored resource for your service.
The easiest way to get these values is to check a previous log entry printed by the same resource to CGP logging.

Below is how you'd set the monitored resource for the above example:
```go
package main

import (
	"os"
	alt4Service "github.com/alt4dev/gcloud/service"
)

func init() {
	// The variables K_SERVICE, K_REVISION, K_CONFIGURATION are automatically added to a running Cloud run container
	alt4Service.SetMonitoredResource("cloud_run_revision", map[string]string{
		"configuration_name": os.Getenv("K_CONFIGURATION"),
		"location": "us-central1",
		"project_id": os.Getenv("PROJECT_ID"),
		"revision_name": os.Getenv("K_REVISION"),
		"service_name": os.Getenv("K_SERVICE"),
    })   
}

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
		labels := make(map[string]interface{})
        // Start a logging group for this request
        group := log.Group(request, labels)
		defer group.Close()
		
		// You can set additional labels to the group or even set the status of the request
		// These steps are not necessary
		group.SetLabel("key", "value")
		group.SetStatus(200)

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

