# cacheuserlibrary

`cacheuserlibrary` is a Go package that provides a user caching mechanism. It includes functionality to cache user data, manage unique user IDs, and save the cache to a file.

## Usage

### Installation

Install the package using the following:

```bash
go get github.com/11marek/cacheuserlibrary
```

### Example


```go
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/your-username/cacheuserlibrary"
)

func main() {
	// Initialize a new caching handler
	cacheHandler := cacheuserlibrary.NewUserCacheHandler(100, &cacheuserlibrary.DatabaseConfig{})

	// Example usage of caching
	userID := "exampleUserID"
	if !cacheHandler.IsUserCached(userID) {
		if cacheHandler.HandleCaching(userID) {
			fmt.Println("User cached successfully.")
		}
	}
}


	
