package server

import (
	"encoding/json"
	"os"
	"sync"
)

// ResultWriter is a writer for writing json-line results to a file.
type ResultWriter struct {
	f   *os.File
	m   sync.Mutex
	enc *json.Encoder
}

// NewResultWriter returns a new ResultWriter.
func NewResultWriter(filename string) (*ResultWriter, error) {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	return &ResultWriter{f: f, enc: json.NewEncoder(f)}, err
}

// AppendResults appends PerformanceResults to a log file.
func (rw *ResultWriter) AppendResults(results *PerformanceResults) error {
	rw.m.Lock()
	defer rw.m.Unlock()
	defer rw.f.Sync()

	return rw.enc.Encode(results)
}

// Close closes the underlying file.
func (rw *ResultWriter) Close() error {
	return rw.f.Close()
}
