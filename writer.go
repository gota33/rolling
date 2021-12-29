package rolling

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Listener interface {
	OnRoll(status Status) (err error)
}

type ListenerFunc func(status Status) (err error)

func (l ListenerFunc) OnRoll(status Status) error { return l(status) }

type Status struct {
	TotalSize int
	TotalNum  int
	Path      string
}

type Config struct {
	Dir        string
	Name       string
	VolumeSize int
	Listener   Listener
}

type Writer struct {
	status Status
	config Config

	fd *os.File
}

func NewWriter(c Config) *Writer {
	return &Writer{config: c}
}

func (w *Writer) Write(p []byte) (n int, err error) {
	if w.reachThreshold(len(p)) {
		if err = w.roll(); err != nil {
			return
		}
		w.status.TotalNum++
	}

	if n, err = w.fd.Write(p); err != nil {
		return
	}
	w.status.TotalSize += n
	return
}

func (w *Writer) Close() (err error) {
	return w.close()
}

func (w *Writer) reachThreshold(size int) bool {
	return w.status.TotalSize+size > w.status.TotalNum*w.config.VolumeSize
}

func (w *Writer) close() (err error) {
	if w.fd == nil {
		return
	}
	if err = w.fd.Close(); err != nil {
		return
	}
	if w.config.Listener == nil {
		return
	}
	return w.config.Listener.OnRoll(w.status)
}

func (w *Writer) roll() (err error) {
	if err = w.close(); err != nil {
		return
	}

	name := fmt.Sprintf(w.config.Name, w.status.TotalNum)
	path := filepath.Join(w.config.Dir, name)
	dir := filepath.Dir(path)

	if _, statErr := os.Stat(dir); errors.Is(statErr, os.ErrNotExist) {
		if err = os.MkdirAll(dir, 0777); err != nil {
			return
		}
	}

	if w.fd, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0666); err != nil {
		return
	}

	w.status.Path = path
	return
}
