package tracker

import "github.com/pkg/errors"

var (
	ErrKeyAddKeyEmpty        = errors.New("cannot add key : key is empty")
	ErrKeyAddKeyExists       = errors.New("cannot add key : key already exists")
	ErrKeyCodeInvalid        = errors.New("invalid key code")
	ErrKeyTrackerInvalidKind = errors.New("invalid kind for KeyTracker")
	ErrKeyUpdateKeyEmpty     = errors.New("cannot update key : key is empty")
)
