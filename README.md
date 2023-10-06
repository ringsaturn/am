# `am`: Apple Maps Server API SDK for Go [![Go Reference](https://pkg.go.dev/badge/github.com/ringsaturn/am.svg)](https://pkg.go.dev/github.com/ringsaturn/am)

```bash
go install github.com/ringsaturn/am
```

## Usage

### Create a client

```go
package main

import (
	"fmt"

	"github.com/ringsaturn/am"
)

func main() {
	client := am.NewClient("your_auth_token")
	fmt.Println(client)
}
```
