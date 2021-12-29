# rolling
Rolling writer

## Example

```go
package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"

	"github/gota33/rolling"
)

const (
	content = `"License" shall mean the terms and conditions for use, reproduction, and distribution as defined by Sections 1 through 9 of this document.`
)

func main() {
	w := rolling.NewWriter(rolling.Config{
		Dir:        filepath.Join(os.TempDir(), "demo"),
		Name:       "demo-%03d.log",
		VolumeSize: 1024 * 1024, // 1MB
		Listener: rolling.ListenerFunc(func(status rolling.Status) (err error) {
			log.Printf("%02d %12d %s", status.TotalNum, status.TotalSize, status.Path)
			return
		}),
	})

	for i := 0; i < math.MaxInt16; i++ {
		_, err := fmt.Fprintf(w, "%s\n", content)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := w.Close(); err != nil {
		log.Fatal(err)
	}
}
```