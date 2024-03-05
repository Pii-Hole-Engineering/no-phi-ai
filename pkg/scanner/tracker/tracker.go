package tracker

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/Pii-Hole-Engineering/no-phi-ai/pkg/scanner/rrr"
)

// KeyCodeToState() function is used to convert the provided key code
// (int) to a string representation of the state.
func KeyCodeToState(code int) string {
	switch code {
	case KeyCodeComplete:
		return KeyStateComplete
	case KeyCodeIgnore:
		return KeyStateIgnore
	case KeyCodeInit:
		return KeyStateInit
	case KeyCodeError:
		return KeyStateError
	case KeyCodePending:
		return KeyStatePending
	default:
		return KeyStateInit
	}
}

// KeyCodeValidate() function is used to validate that the provided
// key code is within the expected range.
func KeyCodeValidate(code int) error {
	if (code < KeyCodeInit) || (code > KeyCodeComplete) {
		return ErrKeyCodeInvalid
	}

	return nil
}

// KeyData struct is used to track the essential data for a given key
// in a KeyTracker instance.
type KeyData struct {
	// Children is a map of child keys that are associated with the
	// current key. The map is used to track scan progress for children
	// as a way to determine when the current key can be marked as complete.
	Children map[string]bool `json:"children"`
	// Code is used to represent the state of the object as an integer,
	// where the value is expected to increase from KeyCodeInit to one
	// of KeyCodeError, KeyCodePending, or KeyCodeComplete.
	Code int `json:"code"`
	// Message is an optional string message that can be used to provide
	// additional context about the state of the object. If the Code is
	// KeyCodeError or KeyCodePending, then the Message should be set.
	Message string `json:"message"`
	// State is a string representation of the Code and is automatically
	// set based on the value of Code.
	State string `json:"state"`
	// TimestampFirst is the timestamp of the first time the associated
	// key was seen in the scan.
	TimestampFirst int64 `json:"timestamp_first"`
	// TimestampLatest is the timestamp of the most recent time the
	// associated key was seen in the scan.
	TimestampLatest int64 `json:"timestamp_latest"`
}

// NewKeyData() function initializes a new KeyData struct with the
// provided code and message. Returns an empty KeyData struct and a
// non-nil error if the code is invalid.
func NewKeyData(code int, message string, child_keys []string) (KeyData, error) {
	if err := KeyCodeValidate(code); err != nil {
		return KeyData{}, errors.Wrapf(err, "failed to create new KeyData with code %d", code)
	}

	now := rrr.TimestampNow()

	children := make(map[string]bool)
	child_state := false
	if code == KeyCodeComplete {
		child_state = true
	}
	for _, key := range child_keys {
		children[key] = child_state
	}

	return KeyData{
		Children:        children,
		Code:            code,
		Message:         message,
		State:           KeyCodeToState(code),
		TimestampFirst:  now,
		TimestampLatest: now,
	}, nil
}

// KeyDataCounts struct is used to count the number of objects in each
// state known to the KeyTracker.
type KeyDataCounts struct {
	Complete int `json:"complete"`
	Error    int `json:"error"`
	Ignore   int `json:"ignore"`
	Init     int `json:"init"`
	Pending  int `json:"pending"`
}

// NewKeyDataCounts() function initializes a new KeyDataCounts struct
// with all counts set to zero and returns the struct.
func NewKeyDataCounts() KeyDataCounts {
	return KeyDataCounts{
		Complete: 0,
		Error:    0,
		Ignore:   0,
		Init:     0,
		Pending:  0,
	}
}

// KeyTracker struct is used to track the state of objects as they are
// scanned in order to prevent duplicate work and to provide a mechanism
// for tracking the progress of the scan.
type KeyTracker struct {
	keys   map[string]KeyData
	kind   string
	logger *zerolog.Logger
	mu     *sync.RWMutex
}

// NewKeyTracker() function initializes a new KeyTracker struct and
// returns a pointer to the struct. The kind parameter is used to
// specify the type of object that the KeyTracker will be used to track.
func NewKeyTracker(kind string, logger *zerolog.Logger) (*KeyTracker, error) {
	switch kind {
	case ScanObjectTypeCommit:
		kind = ScanObjectTypeCommit
	case ScanObjectTypeFile:
		kind = ScanObjectTypeFile
	case ScanObjectTypeRequestResponse:
		kind = ScanObjectTypeRequestResponse
	default:
		return nil, ErrKeyTrackerInvalidKind
	}

	return &KeyTracker{
		keys:   make(map[string]KeyData, 0),
		kind:   kind,
		logger: logger,
		mu:     &sync.RWMutex{},
	}, nil
}

// CheckAllComplete() method checks if the state of all keys in the
// KeyTracker is KeyCodeComplete. Only returns true if all keys are
// complete, otherwise returns false.
func (kt *KeyTracker) CheckAllComplete() bool {
	kt.mu.RLock()
	defer kt.mu.RUnlock()

	for _, key_data := range kt.keys {
		if key_data.Code == KeyCodeInit || key_data.Code == KeyCodePending {
			return false
		}
	}

	return true
}

// Get() method gets the KeyData for the provided key, if it exists in the
// KeyTracker, and returns the KeyData and a boolean indicating whether
// the key exists in the tracker.
func (kt *KeyTracker) Get(key string) (key_data KeyData, exists bool) {
	if key == "" {
		return
	}

	// lock the tracker for reading
	kt.mu.RLock()
	// unlock the tracker after the function returns
	defer kt.mu.RUnlock()

	key_data, exists = kt.keys[key]
	return
}

func (kt *KeyTracker) GetCounts() KeyDataCounts {
	// lock the tracker for reading
	kt.mu.RLock()
	// unlock the tracker after the function returns
	defer kt.mu.RUnlock()

	counts := NewKeyDataCounts()
	for _, key_data := range kt.keys {
		switch key_data.Code {
		case KeyCodeComplete:
			counts.Complete++
		case KeyCodeError:
			counts.Error++
		case KeyCodeIgnore:
			counts.Ignore++
		case KeyCodeInit:
			counts.Init++
		case KeyCodePending:
			counts.Pending++
		}
	}

	return counts
}

// GetKeys() method gets the list of keys known to the KeyTracker, returned
// as a slice of strings.
func (st *KeyTracker) GetKeys() (keys []string) {
	// lock the tracker for reading
	st.mu.RLock()
	// unlock the tracker after the function returns
	defer st.mu.RUnlock()

	for key := range st.keys {
		keys = append(keys, key)
	}

	return
}

// GetKeysData() method gets the map of keys and their associated KeyData.
func (st *KeyTracker) GetKeysData() map[string]KeyData {
	// lock the tracker for reading
	st.mu.RLock()
	// unlock the tracker after the function returns
	defer st.mu.RUnlock()

	// make a copy of the map to return
	out := make(map[string]KeyData, 0)
	for key, key_data := range st.keys {
		out[key] = key_data
	}

	return out
}

// GetKeysDataForCode() method returns a filtered copy of the KeyTracker
// keys map, where the only keys in the returned map are those with a
// Code of KeyCodePending.
func (st *KeyTracker) GetKeysDataForCode(code int) (map[string]KeyData, error) {
	if err := KeyCodeValidate(code); err != nil {
		return nil, errors.Wrapf(err, "failed to get keys data for code %d", code)
	}

	out := make(map[string]KeyData, 0)
	for key, key_data := range st.GetKeysData() {
		if key_data.Code == code {
			out[key] = key_data
		}
	}

	return out, nil
}

func (st *KeyTracker) PrintCodes() []int {
	codes := make([]int, 0)
	for key, key_data := range st.GetKeysData() {
		codes = append(codes, key_data.Code)
		st.logger.Debug().Msgf(
			"PrintCodes :: KIND=%s : KEY=%s : CODE=%d : STATE=%s",
			st.kind,
			key,
			key_data.Code,
			key_data.State,
		)
	}
	return codes
}

func (kt *KeyTracker) PrintCounts() KeyDataCounts {
	counts := kt.GetCounts()
	kt.logger.Debug().Msgf(
		"PrintCounts :: KIND=%s : INIT=%d : ERROR=%d : IGNORE=%d : PENDING=%d : COMPLETE=%d",
		kt.kind,
		counts.Init,
		counts.Error,
		counts.Ignore,
		counts.Pending,
		counts.Complete,
	)
	return counts
}

// Update() method updates the KeyData for the given key with the provided
// code and message. If the key does not exist in the KeyTracker, then it
// will be added.
func (kt *KeyTracker) Update(key string, code_in int, message string, child_keys []string) (code_out int, e error) {
	if key == "" {
		e = ErrKeyUpdateKeyEmpty
		return
	}
	if e = KeyCodeValidate(code_in); e != nil {
		e = errors.Wrapf(e, "failed to update data for key=%s", key)
		return
	}
	// use a read-write lock to update the key data
	kt.mu.Lock()
	// release the lock after the function returns
	defer kt.mu.Unlock()
	key_data, exists := kt.keys[key]
	// check if the key already exists in the kt.keys map
	if !exists {
		// add the key if it does not exist
		k_data, k_err := NewKeyData(code_in, message, child_keys)
		if k_err != nil {
			e = errors.Wrapf(k_err, "failed to update data for key=%s", key)
			return
		}
		kt.keys[key] = k_data
		code_out = k_data.Code
		kt.logger.Trace().Msgf("KIND=%s : created new key=%s with code=%d", kt.kind, key, k_data.Code)
		return
	}
	// refuse to go back to a lower state
	if code_in < key_data.Code {
		code_out = key_data.Code
		return
	}
	// overwrite the message of the existing key data
	key_data.Message = message
	// update the latest timestamp for the existing key data
	key_data.TimestampLatest = rrr.TimestampNow()

	// process the child_keys for "pending" and "complete" states
	switch code_in {
	case KeyCodePending:
		// add/update the "pending" child keys
		for _, child_key := range child_keys {
			// set the value to false to indicate the child key is pending
			key_data.Children[child_key] = false
		}
	case KeyCodeComplete:
		// update the "complete" child keys
		for _, child_key := range child_keys {
			// set the value to true to indicate the child key is complete
			key_data.Children[child_key] = true
		}
		is_complete := true
		for _, is_child_complete := range key_data.Children {
			if !is_child_complete {
				is_complete = false
				break
			}
		}
		// check if all children are complete before marking (this) key_data
		// as complete
		if is_complete {
			key_data.Code = code_in
			key_data.State = KeyCodeToState(code_in)
		}
		// update the existing key data for the completed key state
		kt.keys[key] = key_data
		code_out = key_data.Code
		return
	}

	if code_in != key_data.Code {
		key_data.Code = code_in
		key_data.State = KeyCodeToState(code_in)
	}

	// update the existing key data
	kt.keys[key] = key_data
	code_out = key_data.Code

	return
}
