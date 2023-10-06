# `am`: Apple Maps Server API SDK for Go

```bash
go install github.com/ringsaturn/am
```

## Usage

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
