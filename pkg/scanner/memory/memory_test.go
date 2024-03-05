package memory

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/Pii-Hole-Engineering/no-phi-ai/pkg/scanner/rrr"
)

// TestNewMemoryResultRecordIO() unit test function tests the
// NewMemoryResultRecordIO() function.
func TestNewMemoryResultRecordIO(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	// call the NewMemoryResultRecordIO function
	resultIO := NewMemoryResultRecordIO(ctx)

	assert.Equal(t, zerolog.Ctx(ctx), resultIO.logger)

	assert.NotNil(t, resultIO.mutex)

	// assert that the result_records map is initialized
	assert.NotNil(t, resultIO.result_records)
	assert.Empty(t, resultIO.result_records)
}

// TestMemoryResultRecordIO_Delete() unit test function tests
// the Delete() method of MemoryResultRecordIO struct.
func TestMemoryResultRecordIO_Delete(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	// create a new MemoryResultRecordIO instance
	resultIO := NewMemoryResultRecordIO(ctx)

	expected_err := resultIO.Delete("")
	assert.ErrorIs(t, expected_err, ErrMemoryResultRecordIODeleteEmptyID)

	// add a result record to the map
	id := "123"
	resultIO.result_records[id] = rrr.ResultRecord{}

	// call the Delete method
	err := resultIO.Delete(id)

	// assert that the result record is deleted
	assert.NoError(t, err)
	assert.NotContains(t, resultIO.result_records, id)
}

// TestMemoryResultRecordIO_List unit test function tests the
// List() method of MemoryResultRecordIO struct.
func TestMemoryResultRecordIO_List(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	// create a new MemoryResultRecordIO instance
	resultIO := NewMemoryResultRecordIO(ctx)

	// add some result records to the map
	result1 := rrr.ResultRecord{
		MetadataRequestResponse: rrr.MetadataRequestResponse{
			ID: "1",
		},
		Hash: "hash-1",
	}
	resultIO.result_records["hash-1"] = result1
	result2 := rrr.ResultRecord{
		MetadataRequestResponse: rrr.MetadataRequestResponse{
			ID: "2",
		},
		Hash: "hash-2",
	}
	resultIO.result_records["hash-2"] = result2
	result3 := rrr.ResultRecord{
		MetadataRequestResponse: rrr.MetadataRequestResponse{
			ID: "3",
		},
		Hash: "hash-3",
	}
	resultIO.result_records["hash-3"] = result3

	// call the List method
	results, err := resultIO.List()

	// assert that the list of results is correct
	assert.NoError(t, err)
	assert.Len(t, results, 3)
	assert.Contains(t, results, result1)
	assert.Contains(t, results, result2)
	assert.Contains(t, results, result3)
}

// TestMemoryResultRecordIO_Read unit test function tests the
// Read() method of MemoryResultRecordIO struct.
func TestMemoryResultRecordIO_Read(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	// create a new MemoryResultRecordIO instance
	resultIO := NewMemoryResultRecordIO(ctx)

	// add a result record to the map
	id := "123"
	result := rrr.ResultRecord{
		MetadataRequestResponse: rrr.MetadataRequestResponse{
			ID: "123",
		},
		Hash: "hash-123",
	}
	resultIO.result_records[id] = result

	// call the Read method
	record, err := resultIO.Read(id)

	// assert that the result record is read correctly
	assert.NoError(t, err)
	assert.Equal(t, result, record)
}

// TestMemoryResultRecordIO_ReadEmptyID unit test function tests the
// Read() method of MemoryResultRecordIO struct when an empty ID is provided.
func TestMemoryResultRecordIO_ReadEmptyID(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	// create a new MemoryResultRecordIO instance
	resultIO := NewMemoryResultRecordIO(ctx)

	// call the Read method with an empty ID
	record, err := resultIO.Read("")

	// assert that an error is returned
	assert.Error(t, err)
	assert.Equal(t, ErrMemoryResultRecordIOReadEmptyID, err)
	assert.Equal(t, rrr.ResultRecord{}, record)
}

// TestMemoryResultRecordIO_ReadFailed unit test function tests the
// Read() method of MemoryResultRecordIO struct when the ID is not found.
func TestMemoryResultRecordIO_ReadFailed(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	// create a new MemoryResultRecordIO instance
	resultIO := NewMemoryResultRecordIO(ctx)

	// call the Read method with a non-existent ID
	record, err := resultIO.Read("non-existent-id")

	// assert that an error is returned
	assert.Error(t, err)
	assert.Equal(t, ErrMemoryResultRecordIOReadFailed, err)
	assert.Equal(t, rrr.ResultRecord{}, record)
}

// TestMemoryResultRecordIO_Write unit test function tests the
// Write() method of MemoryResultRecordIO struct.
func TestMemoryResultRecordIO_Write(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	// create a new MemoryResultRecordIO instance
	resultIO := NewMemoryResultRecordIO(ctx)

	// create some result records
	result1 := rrr.ResultRecord{
		MetadataRequestResponse: rrr.MetadataRequestResponse{
			ID: "1",
		},
		Hash: "hash-1",
	}
	result2 := rrr.ResultRecord{
		MetadataRequestResponse: rrr.MetadataRequestResponse{
			ID: "2",
		},
		Hash: "hash-2",
	}
	result3 := rrr.ResultRecord{
		MetadataRequestResponse: rrr.MetadataRequestResponse{
			ID: "3",
		},
		Hash: "hash-3",
	}
	result_records := []rrr.ResultRecord{result1, result2, result3}

	// call the Write method
	err := resultIO.Write(result_records)

	// assert that the result records are written correctly
	assert.NoError(t, err)
	assert.Len(t, resultIO.result_records, 3)
	assert.Equal(t, result1, resultIO.result_records["hash-1"])
	assert.Equal(t, result2, resultIO.result_records["hash-2"])
	assert.Equal(t, result3, resultIO.result_records["hash-3"])
}
