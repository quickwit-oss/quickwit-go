## Quickwit Go Client

WARNING: This is a work in progress and can be used only for testing purposes.

## Installation

```bash
go get github.com/quickwit-oss/quickwit-go
```

## Testing the client

### Start a Quickwit instance

```bash
docker run -it --rm -p 7280:7280 quickwit/quickwit
```

### Execute a search query

```go
package main

import (
    "fmt"
    "github.com/quickwit-oss/quickwit-go"
)

func main() {
    qclient := quickwit.NewQuickwitClient("http://localhost:7280")
    searchRequest := quickwit.SearchRequest{Query: "severity_text:error"}
    // otel-logs-v0 is created when quickwit starts.
    searchResponse, err := qclient.Search("otel-logs-v0", searchRequest)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println("--------------------")
    fmt.Println("response", searchResponse)
}
```
