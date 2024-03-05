package scanner

import "github.com/pkg/errors"

const (
	ErrMsgAddScanRepository     = "failed to add ScanRepository"
	ErrMsgCloneRepository       = "failed to clone repository"
	ErrMsgErrorChannelNil       = "received nil error channel as input"
	ErrMsgResultWriteFailed     = "failed to write result"
	ErrMsgScanRepositoryCreate  = "failed to create new ScanRepository object"
	ErrMsgScanRepositoryScan    = "failed to scan repository"
	ErrMsgScanTrackerUpdateFile = "failed to update tracker for file %s"
	ErrMsgScannerCreate         = "failed to create new Scanner"
	ErrMsgTrackerUpdateCommit   = "failed to update tracker for commit %s"
)

var (
	ErrProcessRequestNoID               = errors.New("cannot process a request without a valid ID")
	ErrProcessResponseNoID              = errors.New("cannot process a response without a valid ID")
	ErrScannerAddScanRepositoryEmptyID  = errors.New("cannot add a ScanRepository with an empty ID")
	ErrScannerAddScanRepositoryNil      = errors.New("cannot add a nil ScanRepository to scanner")
	ErrScannerGetScanRepositoryNotFound = errors.New("ScanRepository not found")
	ErrScannerRepositoryNil             = errors.New("Scanner cannot scan repository with nil pointer")
	ErrScanRepositoryChannelErrorsNil   = errors.New("ScanRepository errors channel is nil")
	ErrScanRepositoryChannelRequestsNil = errors.New("ScanRepository requests channel is nil")
	ErrScanRepositoryConfigNil          = errors.New("ScanRepository config is nil")
	ErrScanRepositoryContextNil         = errors.New("ScanRepository requires a non-nil context")
	ErrScanRepositoryCloneGitManagerNil = errors.New("ScanRepository git manager is nil")
	ErrScanRepositoryRepositoryNil      = errors.New("ScanRepository repository pointer is nil")
)
