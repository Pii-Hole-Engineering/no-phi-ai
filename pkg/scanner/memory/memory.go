package memory

import (
	"context"
	"sync"

	"github.com/rs/zerolog"

	"github.com/Pii-Hole-Engineering/no-phi-ai/pkg/scanner/rrr"
)

// MemoryResultRecordIO struct provides an in-memory implementation
// of the ResultRecordIO interface. Useful for development, testing,
// and scans of small-scale repositories.
type MemoryResultRecordIO struct {
	rrr.ResultRecordIO

	logger         *zerolog.Logger
	mutex          *sync.RWMutex
	result_records map[string]rrr.ResultRecord
}

// NewMemoryResultRecordIO() function initializes a new MemoryResultRecordIO object.
func NewMemoryResultRecordIO(ctx context.Context) MemoryResultRecordIO {
	return MemoryResultRecordIO{
		logger:         zerolog.Ctx(ctx),
		mutex:          &sync.RWMutex{},
		result_records: make(map[string]rrr.ResultRecord),
	}
}

// Delete() method deletes the result with matching id from the memory store.
func (io MemoryResultRecordIO) Delete(id string) error {
	if id == "" {
		return ErrMemoryResultRecordIODeleteEmptyID
	}
	io.logger.Debug().Msgf("deleting result id=%s from memory store", id)
	io.mutex.Lock()
	defer io.mutex.Unlock()

	delete(io.result_records, id)
	return nil
}

// List() method returns a list of all results in the memory store.
func (io MemoryResultRecordIO) List() ([]rrr.ResultRecord, error) {
	io.logger.Debug().Msg("listing results from memory store")
	io.mutex.RLock()
	current_results := io.result_records
	io.mutex.RUnlock()

	var out []rrr.ResultRecord
	for _, result := range current_results {
		out = append(out, result)
	}
	return out, nil
}

// Read() method returns the result with matching id from the memory store.
// Returns a non-nil error if unable to find a result with matching id.
func (io MemoryResultRecordIO) Read(id string) (rrr.ResultRecord, error) {
	if id == "" {
		return rrr.ResultRecord{}, ErrMemoryResultRecordIOReadEmptyID
	}
	io.logger.Debug().Msgf("reading result id=%s from memory store", id)
	io.mutex.RLock()
	defer io.mutex.RUnlock()

	r, ok := io.result_records[id]
	if !ok {
		return rrr.ResultRecord{}, ErrMemoryResultRecordIOReadFailed
	}
	return r, nil
}

// Write() method writes the slice of results to the memory store. Returns a
// non-nil error if unable to write any result to the store.
func (io MemoryResultRecordIO) Write(result_records []rrr.ResultRecord) error {
	io.logger.Debug().Msgf("writing %d result(s) to memory store", len(result_records))
	io.mutex.Lock()
	defer io.mutex.Unlock()

	for _, r := range result_records {
		io.result_records[r.Hash] = r
	}
	return nil
}
