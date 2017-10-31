package tools

import (
	"os"
	"sync"
	"time"
)

type RotateWriter struct {
	lock      sync.Mutex
	filename  string // should be set to the actual filename
	fp        *os.File
	createdAt time.Time
	//Now()
}

// Make a new RotateWriter. Return nil if error occurs during setup.
func NewRotateWriter(filename string) *RotateWriter {
	w := &RotateWriter{filename: filename}
	err := w.rotate()
	w.createdAt = time.Now()
	if err != nil {
		return nil
	}
	return w
}

// Write satisfies the io.Writer interface.
func (w *RotateWriter) Write(output []byte) (int, error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	if w.ShouldRotate() {
		w.rotate()
		w.createdAt = time.Now()
	}
	return w.fp.Write(output)
}

// lets expose this so we can defer close on termination
// to avoid leaking file descriptors
func (w *RotateWriter) Close() {
	w.fp.Close()
}

// helper method to check if its time to rotate some logs
func (w *RotateWriter) ShouldRotate() bool {
	if time.Now().Sub(w.createdAt).Hours() >= 24.0 {
		return true
	}
	return false
}

// Perform the actual act of rotating and reopening file.
func (w *RotateWriter) rotate() (err error) {
	//	w.lock.Lock()
	//	defer w.lock.Unlock()

	// Close existing file if open
	if w.fp != nil {
		err = w.fp.Close()
		w.fp = nil
		if err != nil {
			return
		}
	}
	// Rename dest file if it already exists
	_, err = os.Stat(w.filename)
	if err == nil {
		err = os.Rename(w.filename, w.filename+"."+time.Now().Format(time.RFC3339))
		if err != nil {
			return
		}
	}

	// Create a file.
	w.fp, err = os.Create(w.filename)
	return
}
