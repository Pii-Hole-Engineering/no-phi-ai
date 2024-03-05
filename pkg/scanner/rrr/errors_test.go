package rrr

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestErrors(t *testing.T) {
	tests := []struct {
		err  error
		name string
	}{
		{
			err:  ErrChunkFileToRequestsFailed,
			name: "ErrChunkFileToRequestsFailed",
		},
		{
			err:  ErrNewRequestEmptyCommitID,
			name: "ErrNewRequestEmptyCommitID",
		},
		{
			err:  ErrNewRequestEmptyObjectID,
			name: "ErrNewRequestEmptyObjectID",
		},
		{
			err:  ErrNewRequestEmptyRepositoryID,
			name: "ErrNewRequestEmptyRepositoryID",
		},
		{
			err:  ErrNewRequestEmptyText,
			name: "ErrNewRequestEmptyText",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			new_err := errors.New(test.err.Error())
			assert.Error(t, test.err)
			assert.Equal(t, test.err.Error(), new_err.Error())
		})
	}
}
