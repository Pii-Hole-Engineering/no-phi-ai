package memory

import "github.com/pkg/errors"

var (
	ErrMemoryResultRecordIODeleteEmptyID = errors.New("memory store failed to delete result record : empty ID")
	ErrMemoryResultRecordIOReadEmptyID   = errors.New("memory store failed to read result record : empty ID")
	ErrMemoryResultRecordIOReadFailed    = errors.New("memory store failed to read result record")
)
