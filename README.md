# Monolog-Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/MonoLog-Connector/monolog-go.svg)](https://pkg.go.dev/github.com/MonoLog-Connector/monolog-go/client)

Monolog-Go is a lightweight SDK for tracking logs in Go applications, especially for integrating with the `gin-gonic` framework. It provides a middleware that can be used with the Gin router to track and log requests in a specified log file.

## Installation

To use the SDK, you need to first install it using `go get`.

```bash
go get github.com/MonoLog-Connector/monolog-go/client
```

## Usage Example

Below is an example of how you can integrate the Monolog-Go SDK with a Gin router in your Go application.

### Step 1: Import the package

```
import (
    "github.com/MonoLog-Connector/monolog-go/client"
    "github.com/gin-gonic/gin"
)
```
### Step 2: Initialize the SDK

Initialize the SDK by passing the path to the log file where the request details will be stored.

```
func main() {
    router := gin.Default()

    // Initialize SDK with the path to the log file
    monologsdk := client.NewSDK("/path/to/your/log/file.log")
    // Use the middleware provided by the SDK
    router.Use(monologsdk.GinTrackerMiddleware())
    router.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })
    router.Run(":8080")
}
```

### Step 3: Now test it

After starting the server, you can test the setup by making a request to the endpoint:

```
curl http://localhost:8080/ping
```

### Step 4: View the logs

Once the request is made, the logs will be written to the specified log file. Below is an example of the generated logs:

```
{
    "CPU Delta": 74.19527896995683,
    "DateTime": "2024-10-06T01:38:19+05:30",
    "Latency": "529.292µs",
    "Memory Delta (MB)": 0.390625,
    "RequestMethod": "GET",
    "RequestURL": "/ping",
    "Status": 404,
    "level": "info",
    "msg": "Request details logged",
    "time": "2024-10-06T01:38:19+05:30"
}
{
    "level": "info",
    "msg": "Shutdown signal received, cleaning up...",
    "time": "2024-10-06T01:39:06+05:30"
}
{
    "CPU Delta": 1.1155807145198469,
    "DateTime": "2024-10-06T01:40:53+05:30",
    "Latency": "5.702657166s",
    "Memory Delta (MB)": 9.84375,
    "RequestMethod": "GET",
    "RequestURL": "/ping",
    "Status": 200,
    "level": "info",
    "msg": "Request details logged",
    "time": "2024-10-06T01:40:59+05:30"
}


```