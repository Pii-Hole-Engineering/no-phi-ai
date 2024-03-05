package rrr

import "github.com/pkg/errors"

var (
	ErrChunkFileToRequestsFailed    = errors.New("failed to chunk non-empty file into one or more requests")
	ErrChunkFileToRequestsInFileNil = errors.New("cannot chunk input file with nil pointer")
	ErrMaxChunkSizeInvalid          = errors.New("invalid max chunk size")
	ErrNewRequestEmptyCommitID      = errors.New("cannot create a new request with an empty commit ID")
	ErrNewRequestEmptyObjectID      = errors.New("cannot create a new request with an empty object ID")
	ErrNewRequestEmptyRepositoryID  = errors.New("cannot create a new request with an empty repository ID")
	ErrNewRequestEmptyText          = errors.New("cannot create a new request with an empty text")
)
