package server

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// ResultWriter is a writer for writing json-line results to a file.
type ResultWriter struct {
	f   *os.File
	m   sync.Mutex
	enc *json.Encoder
}

// Metadata is additional data to be written to log.
type Metadata struct {
	DateTime time.Time `json:"datetime"`
	IP       string    `json:"ip"`
}

// NewResultWriter returns a new ResultWriter.
func NewResultWriter(filename string) (*ResultWriter, error) {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	return &ResultWriter{f: f, enc: json.NewEncoder(f)}, err
}

// AppendResults appends PerformanceResults to a log file.
func (rw *ResultWriter) AppendResults(metadata *Metadata, results *PerformanceResults) error {
	rw.m.Lock()
	defer rw.m.Unlock()

	type log struct {
		*Metadata
		*PerformanceResults
	}

	return rw.enc.Encode(&log{metadata, results})
}

// Close closes the underlying file.
func (rw *ResultWriter) Close() error {
	return rw.f.Close()
}
