package livereload

import (
	"context"
	"os"
	"sync"
	"time"
)

type FileInfo struct {
	Name         string
	LastModified time.Time
	Delay        time.Duration
}

type FileWatcher struct {
	mtx   *sync.RWMutex
	files []*FileInfo
}

func NewFileWatcher(files []*FileInfo) *FileWatcher {
	fw := &FileWatcher{
		mtx:   new(sync.RWMutex),
		files: files,
	}

	return fw
}

func (w *FileWatcher) AddFile(name string, delay time.Duration) {
	w.files = append(w.files, &FileInfo{Name: name, Delay: delay})
}
func (w *FileWatcher) isModified(file *FileInfo) bool {
	fi, err := os.Stat(file.Name)
	if err != nil {
		return false
	}

	if file.LastModified.IsZero() {
		file.LastModified = fi.ModTime()
		return false
	}

	diff := fi.ModTime().Sub(file.LastModified)
	file.LastModified = fi.ModTime()
	if diff > file.Delay {
		// log.Println("modified! file:", file.Name, "modTime:", fi.ModTime(), "dif:", diff)
		return true
	}

	return false
}

func (w *FileWatcher) filesHaveBeenModified() bool {

	for i := range w.files {
		if w.isModified(w.files[i]) {
			return true
		}
	}

	return false
}

// Run loops forever until it finds a modified file then it sends "reload" to the supplied channel
// run will exit if the supplied context is done
func (w *FileWatcher) Run(ctx context.Context, outbox chan string) {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// if any of the watched files have been modified, send reload
			if w.filesHaveBeenModified() {
				outbox <- "reload"
			}
		}
	}
}
