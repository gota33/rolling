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
	content = `Each major Go release is supported until there are two newer major releases. For example, Go 1.5 was supported until the Go 1.7 release, and Go 1.6 was supported until the Go 1.8 release. We fix critical problems, including critical security problems, in supported releases as needed by issuing minor revisions (for example, Go 1.6.1, Go 1.6.2, and so on).`
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
