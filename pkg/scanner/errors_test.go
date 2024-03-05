package scanner

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
			err:  ErrProcessRequestNoID,
			name: "ErrProcessRequestNoID",
		},
		{
			err:  ErrProcessResponseNoID,
			name: "ErrProcessResponseNoID",
		},
		{
			err:  ErrScannerAddScanRepositoryEmptyID,
			name: "ErrScannerAddScanRepositoryEmptyID",
		},
		{
			err:  ErrScannerAddScanRepositoryNil,
			name: "ErrScannerAddScanRepositoryNil",
		},
		{
			err:  ErrScannerGetScanRepositoryNotFound,
			name: "ErrScannerGetScanRepositoryNotFound",
		},
		{
			err:  ErrScannerRepositoryNil,
			name: "ErrScannerRepositoryNil",
		},
		{
			err:  ErrScanRepositoryChannelRequestsNil,
			name: "ErrScanRepositoryChannelRequestsNil",
		},
		{
			err:  ErrScanRepositoryConfigNil,
			name: "ErrScanRepositoryConfigNil",
		},
		{
			err:  ErrScanRepositoryContextNil,
			name: "ErrScanRepositoryContextNil",
		},
		{
			err:  ErrScanRepositoryCloneGitManagerNil,
			name: "ErrScanRepositoryCloneGitManagerNil",
		},
		{
			err:  ErrScanRepositoryRepositoryNil,
			name: "ErrScanRepositoryRepositoryNil",
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
