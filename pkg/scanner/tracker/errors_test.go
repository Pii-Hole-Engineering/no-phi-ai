package tracker

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
			err:  ErrKeyAddKeyEmpty,
			name: "ErrKeyAddKeyEmpty",
		},
		{
			err:  ErrKeyAddKeyExists,
			name: "ErrKeyAddKeyExists",
		},
		{
			err:  ErrKeyCodeInvalid,
			name: "ErrKeyCodeInvalid",
		},
		{
			err:  ErrKeyTrackerInvalidKind,
			name: "ErrKeyTrackerInvalidKind",
		},
		{
			err:  ErrKeyUpdateKeyEmpty,
			name: "ErrKeyUpdateKeyEmpty",
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
