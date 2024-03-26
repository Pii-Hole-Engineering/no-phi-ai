package tracker

import (
	"os"
	"sync"
	"testing"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/Pii-Hole-Engineering/no-phi-ai/pkg/scanner/rrr"
)

var (
	test_message_complete = "test message complete"
	test_message_ignore   = "test message ignore"
	test_message_init     = "test message init"
	test_message_error    = "test message error"
	test_message_pending  = "test message pending"
)

// TestKeyCodeToState() unit test function is used to test the KeyCodeToState()
// function.
func TestKeyCodeToState(t *testing.T) {
	t.Parallel()
	tests := []struct {
		code     int
		expected string
		name     string
	}{
		{
			code:     KeyCodeComplete,
			expected: KeyStateComplete,
			name:     "KeyCodeComplete",
		},
		{
			code:     KeyCodeIgnore,
			expected: KeyStateIgnore,
			name:     "KeyCodeIgnore",
		},
		{
			code:     KeyCodeInit,
			expected: KeyStateInit,
			name:     "KeyCodeInit",
		},
		{
			code:     KeyCodeError,
			expected: KeyStateError,
			name:     "KeyCodeError",
		},
		{
			code:     KeyCodePending,
			expected: KeyStatePending,
			name:     "KeyCodePending",
		},
		{
			code:     123, // Replace with your custom code value
			expected: KeyStateInit,
			name:     "CustomCode",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := KeyCodeToState(test.code)
			assert.Equal(t, test.expected, result)
		})
	}
}

// TestKeyCodeValidate() unit test function tests the KeyCodeValidate() function.
func TestKeyCodeValidate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		code     int
		expected error
		name     string
	}{
		{
			code:     KeyCodeInit,
			expected: nil,
			name:     "KeyCodeInit",
		},
		{
			code:     KeyCodeError,
			expected: nil,
			name:     "KeyCodeError",
		},
		{
			code:     KeyCodeIgnore,
			expected: nil,
			name:     "KeyCodeIgnore",
		},
		{
			code:     KeyCodePending,
			expected: nil,
			name:     "KeyCodePending",
		},
		{
			code:     KeyCodeComplete,
			expected: nil,
			name:     "KeyCodeComplete",
		},
		{
			code:     KeyCodeInit - 1,
			expected: ErrKeyCodeInvalid,
			name:     "Code_Invalid_Low",
		},
		{
			code:     KeyCodeComplete + 1,
			expected: ErrKeyCodeInvalid,
			name:     "Code_Invalid_High",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := KeyCodeValidate(test.code)
			assert.Equal(t, test.expected, err)
		})
	}
}

// TestNewKeyData() unit test function tests the NewKeyData() function.
func TestNewKeyData(t *testing.T) {
	t.Parallel()
	tests := []struct {
		code             int
		expected_data    KeyData
		expected_err     error
		expected_message string
		message          string
		name             string
	}{
		{
			code: KeyCodeInit,
			expected_data: KeyData{
				Code:  KeyCodeInit,
				State: KeyStateInit,
			},
			expected_err:     nil,
			expected_message: test_message_init,
			message:          test_message_init,
			name:             "ValidCodeInit",
		},
		{
			code: KeyCodeInit,
			expected_data: KeyData{
				Code:  KeyCodeInit,
				State: KeyStateInit,
			},
			expected_err:     nil,
			expected_message: "",
			message:          "",
			name:             "ValidCodeInitMessageEmpty",
		},
		{
			code: KeyCodeError,
			expected_data: KeyData{
				Code:  KeyCodeError,
				State: KeyStateError,
			},
			expected_err:     nil,
			expected_message: test_message_error,
			message:          test_message_error,
			name:             "ValidCodeError",
		},
		{
			code: KeyCodeIgnore,
			expected_data: KeyData{
				Code:  KeyCodeIgnore,
				State: KeyStateIgnore,
			},
			expected_err:     nil,
			expected_message: test_message_ignore,
			message:          test_message_ignore,
			name:             "ValidCodeIgnore",
		},
		{
			code: KeyCodePending,
			expected_data: KeyData{
				Code:  KeyCodePending,
				State: KeyStatePending,
			},
			expected_err:     nil,
			expected_message: test_message_pending,
			message:          test_message_pending,
			name:             "ValidCodePending",
		},
		{
			code: KeyCodeComplete,
			expected_data: KeyData{
				Code:  KeyCodeComplete,
				State: KeyStateComplete,
			},
			expected_err:     nil,
			expected_message: test_message_complete,
			message:          test_message_complete,
			name:             "ValidCodeComplete",
		},
		{
			code:             KeyCodeInit - 1,
			expected_data:    KeyData{},
			expected_err:     ErrKeyCodeInvalid,
			expected_message: "",
			message:          "",
			name:             "InvalidCodeLow",
		},
		{
			code:             KeyCodeComplete + 1,
			expected_data:    KeyData{},
			expected_err:     ErrKeyCodeInvalid,
			expected_message: "",
			message:          "",
			name:             "InvalidCodeHigh",
		},
	}

	timestamp_test_min := rrr.TimestampNow()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data, err := NewKeyData(test.code, test.message, []string{})
			if test.expected_err != nil {
				assert.ErrorContains(t, err, test.expected_err.Error())
				return
			}
			timestamp_test_max := rrr.TimestampNow()
			assert.NoError(t, err)
			assert.Equal(t, test.expected_data.Code, data.Code)
			assert.Equal(t, test.expected_data.State, data.State)
			assert.Equal(t, test.expected_message, data.Message)
			assert.GreaterOrEqual(t, data.TimestampFirst, timestamp_test_min)
			assert.LessOrEqual(t, data.TimestampFirst, timestamp_test_max)
			assert.GreaterOrEqual(t, data.TimestampLatest, timestamp_test_min)
			assert.LessOrEqual(t, data.TimestampLatest, timestamp_test_max)
			assert.Exactly(t, data.TimestampFirst, data.TimestampLatest, "TimestampFirst and TimestampLatest should be equal for a new KeyData instance.")
		})
	}
}

// TestNewKeyDataCounts() unit test function tests the NewKeyDataCounts function.
func TestNewKeyDataCounts(t *testing.T) {
	t.Parallel()

	expectedCounts := KeyDataCounts{
		Complete: 0,
		Error:    0,
		Ignore:   0,
		Init:     0,
		Pending:  0,
	}

	counts := NewKeyDataCounts()

	assert.Equal(t, expectedCounts, counts)
}

// TestNewKeyTracker() unit test function tests the NewKeyTracker function.
func TestNewKeyTracker(t *testing.T) {
	t.Parallel()

	logger := zerolog.New(os.Stdout)
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	tests := []struct {
		kind         string
		expected     *KeyTracker
		expected_err error
		name         string
	}{
		{
			kind: ScanObjectTypeCommit,
			expected: &KeyTracker{
				Keys:   make(KeyDataMap),
				Kind:   ScanObjectTypeCommit,
				logger: &logger,
				mu:     &sync.RWMutex{},
			},
			expected_err: nil,
			name:         "ValidKindCommit",
		},
		{
			kind: ScanObjectTypeFile,
			expected: &KeyTracker{
				Keys:   make(KeyDataMap),
				Kind:   ScanObjectTypeFile,
				logger: &logger,
				mu:     &sync.RWMutex{},
			},
			expected_err: nil,
			name:         "ValidKindFile",
		},
		{
			kind: ScanObjectTypeRequestResponse,
			expected: &KeyTracker{
				Keys:   make(KeyDataMap),
				Kind:   ScanObjectTypeRequestResponse,
				logger: &logger,
				mu:     &sync.RWMutex{},
			},
			expected_err: nil,
			name:         "ValidKindRequestResponse",
		},
		{
			kind:         "InvalidKind",
			expected:     nil,
			expected_err: ErrKeyTrackerInvalidKind,
			name:         "InvalidKind",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tracker, err := NewKeyTracker(test.kind, &logger)
			if test.expected_err != nil {
				assert.ErrorIs(t, err, test.expected_err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, test.expected.Keys, tracker.Keys)
			assert.Equal(t, test.expected.Kind, tracker.Kind)
			assert.Equal(t, test.expected.logger, tracker.logger)
			assert.Equal(t, test.expected.mu, tracker.mu)
		})
	}
}

// TestKeyTracker_CheckAllComplete() unit test function tests the
// CheckAllComplete() method of the KeyTracker type.
func TestKeyTracker_CheckAllComplete(t *testing.T) {
	t.Parallel()

	logger := zerolog.New(os.Stdout)
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	type testData struct {
		key   string
		codes []int
	}

	tests := []struct {
		name            string
		data            []testData
		expect_complete bool
		kind            string
	}{
		{
			name: "Init_1",
			data: []testData{
				{
					codes: []int{KeyCodeInit},
					key:   "A",
				},
			},
			expect_complete: false,
			kind:            ScanObjectTypeCommit,
		},
		{
			name: "Init_2",
			data: []testData{
				{
					codes: []int{KeyCodeInit},
					key:   "A",
				},
			},
			expect_complete: false,
			kind:            ScanObjectTypeCommit,
		},
		{
			name: "Complete_1",
			data: []testData{
				{
					codes: []int{KeyCodeComplete},
					key:   "A",
				},
			},
			expect_complete: true,
			kind:            ScanObjectTypeCommit,
		},
		{
			name: "Complete_2",
			data: []testData{
				{
					codes: []int{KeyCodeComplete, 2},
					key:   "A",
				},
			},
			expect_complete: true,
			kind:            ScanObjectTypeCommit,
		},
		{
			name: "Error_1",
			data: []testData{
				{
					codes: []int{KeyCodeError},
					key:   "A",
				},
			},
			expect_complete: true,
			kind:            ScanObjectTypeCommit,
		},
		{
			name: "Error_2",
			data: []testData{
				{
					codes: []int{KeyCodeError, -1},
					key:   "A",
				},
			},
			expect_complete: true,
			kind:            ScanObjectTypeFile,
		},
		{
			name: "Mixed_1",
			data: []testData{
				{
					codes: []int{KeyCodeInit},
					key:   "A",
				},
				{
					codes: []int{KeyCodeInit, KeyCodeError, KeyCodeInit},
					key:   "B",
				},
				{
					codes: []int{KeyCodeIgnore},
					key:   "C",
				},
				{
					codes: []int{KeyCodeInit, KeyCodePending, KeyCodeError},
					key:   "D",
				},
				{
					codes: []int{KeyCodeInit, KeyCodeError, KeyCodeComplete},
					key:   "E",
				},
			},
			expect_complete: false,
			kind:            ScanObjectTypeFile,
		},
		{
			name: "Progression",
			data: []testData{
				{
					codes: []int{-2, -1, 0, 1, 2},
					key:   "A",
				},
				{
					codes: []int{
						KeyCodeInit,
						KeyCodeError,
						KeyCodeIgnore,
						KeyCodePending,
						KeyCodeComplete,
					},
					key: "B",
				},
			},
			expect_complete: true,
			kind:            ScanObjectTypeFile,
		},
		{
			name: "Regression",
			data: []testData{
				{
					codes: []int{2, 1, 0, -1, -2},
					key:   "A",
				},
				{
					codes: []int{
						KeyCodeComplete,
						KeyCodePending,
						KeyCodeIgnore,
						KeyCodeError,
						KeyCodeInit,
					},
					key: "B",
				},
			},
			expect_complete: true,
			kind:            ScanObjectTypeCommit,
		},
		{
			name: "ReInit",
			data: []testData{
				{
					codes: []int{
						-2,
						-1,
						0,
						1,
						2,
						-2,
						-2,
						-2,
					},
					key: "A",
				},
				{
					codes: []int{
						KeyCodeInit,
						KeyCodeError,
						KeyCodeIgnore,
						KeyCodePending,
						KeyCodeComplete,
						KeyCodeInit,
						KeyCodeInit,
						KeyCodeInit,
					},
					key: "B",
				},
			},
			expect_complete: true,
			kind:            ScanObjectTypeCommit,
		},
	}

	for _, test_i := range tests {
		t.Run(test_i.name, func(t *testing.T) {
			tracker, err := NewKeyTracker(test_i.kind, &logger)
			if !assert.NoError(t, err) {
				t.FailNow()
			}

			// for each key in the test data, apply the series of update codes
			for _, d := range test_i.data {
				for _, code := range d.codes {
					_, update_err := tracker.Update(d.key, code, "", []string{})
					assert.NoError(t, update_err)
				}
			}

			// get the data for the key after applying all updates
			is_complete := tracker.CheckAllComplete()
			assert.Equal(
				t,
				test_i.expect_complete,
				is_complete,
				"CheckAllComplete() did not return the expected value",
			)
		})
	}
}

// TestKeyTracker_Get() unit test function tests the Get() method of the
// KeyTracker type.
func TestKeyTracker_Get(t *testing.T) {
	t.Parallel()

	// Create a new KeyTracker instance
	tracker := &KeyTracker{
		Keys:   KeyDataMap{},
		Kind:   ScanObjectTypeFile,
		logger: nil,
		mu:     &sync.RWMutex{},
	}

	// add some test data to the tracker
	testKey := "testKey"
	testData := KeyData{
		Code:  KeyCodeInit,
		State: KeyStateInit,
	}
	tracker.Keys[testKey] = testData

	// Test case: Existing key
	t.Run("ExistingKey", func(t *testing.T) {
		keyData, exists := tracker.Get(testKey)
		assert.True(t, exists)
		assert.Equal(t, testData, keyData)
	})

	// Test case: Non-existing key
	t.Run("NonExistingKey", func(t *testing.T) {
		keyData, exists := tracker.Get("nonExistingKey")
		assert.False(t, exists)
		assert.Equal(t, KeyData{}, keyData)
	})

	// Test case: Empty key
	t.Run("EmptyKey", func(t *testing.T) {
		keyData, exists := tracker.Get("")
		assert.False(t, exists)
		assert.Equal(t, KeyData{}, keyData)
	})
}

// TestKeyTracker_GetCounts() unit test function tests the GetCounts()
// method of the KeyTracker type.
func TestKeyTracker_GetCounts(t *testing.T) {
	t.Parallel()

	logger := zerolog.New(os.Stdout)
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	type testData struct {
		key   string
		codes []int
	}

	tests := []struct {
		name      string
		data      []testData
		final     KeyDataCounts
		final_err error
		kind      string
	}{
		{
			name:      "InvalidKind",
			data:      []testData{},
			final:     KeyDataCounts{},
			final_err: ErrKeyTrackerInvalidKind,
			kind:      "InvalidKind",
		},
		{
			name: "Init_1",
			data: []testData{
				{
					codes: []int{KeyCodeInit},
					key:   "A",
				},
			},
			final: KeyDataCounts{
				Complete: 0,
				Error:    0,
				Ignore:   0,
				Init:     1,
				Pending:  0,
			},
			final_err: nil,
			kind:      ScanObjectTypeCommit,
		},
		{
			name: "Init_2",
			data: []testData{
				{
					codes: []int{KeyCodeInit},
					key:   "A",
				},
			},
			final: KeyDataCounts{
				Complete: 0,
				Error:    0,
				Ignore:   0,
				Init:     1,
				Pending:  0,
			},
			final_err: nil,
			kind:      ScanObjectTypeCommit,
		},
		{
			name: "Complete_1",
			data: []testData{
				{
					codes: []int{KeyCodeComplete},
					key:   "A",
				},
			},
			final: KeyDataCounts{
				Complete: 1,
				Error:    0,
				Ignore:   0,
				Init:     0,
				Pending:  0,
			},
			final_err: nil,
			kind:      ScanObjectTypeCommit,
		},
		{
			name: "Complete_2",
			data: []testData{
				{
					codes: []int{KeyCodeComplete, 2},
					key:   "A",
				},
			},
			final: KeyDataCounts{
				Complete: 1,
				Error:    0,
				Ignore:   0,
				Init:     0,
				Pending:  0,
			},
			final_err: nil,
			kind:      ScanObjectTypeCommit,
		},
		{
			name: "Error_1",
			data: []testData{
				{
					codes: []int{KeyCodeError},
					key:   "A",
				},
			},
			final: KeyDataCounts{
				Complete: 0,
				Error:    1,
				Ignore:   0,
				Init:     0,
				Pending:  0,
			},
			final_err: nil,
			kind:      ScanObjectTypeCommit,
		},
		{
			name: "Error_2",
			data: []testData{
				{
					codes: []int{KeyCodeError, -1},
					key:   "A",
				},
			},
			final: KeyDataCounts{
				Complete: 0,
				Error:    1,
				Ignore:   0,
				Init:     0,
				Pending:  0,
			},
			final_err: nil,
			kind:      ScanObjectTypeFile,
		},
		{
			name: "Mixed_1",
			data: []testData{
				{
					codes: []int{KeyCodeInit},
					key:   "A",
				},
				{
					codes: []int{KeyCodeInit, KeyCodeError, KeyCodeInit},
					key:   "B",
				},
				{
					codes: []int{KeyCodeIgnore},
					key:   "C",
				},
				{
					codes: []int{KeyCodeInit, KeyCodePending, KeyCodeError},
					key:   "D",
				},
				{
					codes: []int{KeyCodeInit, KeyCodeError, KeyCodeComplete},
					key:   "E",
				},
			},
			final: KeyDataCounts{
				Complete: 1,
				Error:    1,
				Ignore:   1,
				Init:     1,
				Pending:  1,
			},
			final_err: nil,
			kind:      ScanObjectTypeFile,
		},
		{
			name: "Progression",
			data: []testData{
				{
					codes: []int{-2, -1, 0, 1, 2},
					key:   "A",
				},
				{
					codes: []int{
						KeyCodeInit,
						KeyCodeError,
						KeyCodeIgnore,
						KeyCodePending,
						KeyCodeComplete,
					},
					key: "B",
				},
			},
			final: KeyDataCounts{
				Complete: 2,
				Error:    0,
				Ignore:   0,
				Init:     0,
				Pending:  0,
			},
			final_err: nil,
			kind:      ScanObjectTypeFile,
		},
		{
			name: "Regression",
			data: []testData{
				{
					codes: []int{2, 1, 0, -1, -2},
					key:   "A",
				},
				{
					codes: []int{
						KeyCodeComplete,
						KeyCodePending,
						KeyCodeIgnore,
						KeyCodeError,
						KeyCodeInit,
					},
					key: "B",
				},
			},
			final: KeyDataCounts{
				Complete: 2,
				Error:    0,
				Ignore:   0,
				Init:     0,
				Pending:  0,
			},
			final_err: nil,
			kind:      ScanObjectTypeCommit,
		},
		{
			name: "ReInit",
			data: []testData{
				{
					codes: []int{
						-2,
						-1,
						0,
						1,
						2,
						-2,
						-2,
						-2,
					},
					key: "A",
				},
				{
					codes: []int{
						KeyCodeInit,
						KeyCodeError,
						KeyCodeIgnore,
						KeyCodePending,
						KeyCodeComplete,
						KeyCodeInit,
						KeyCodeInit,
						KeyCodeInit,
					},
					key: "B",
				},
			},
			final: KeyDataCounts{
				Complete: 2,
				Error:    0,
				Ignore:   0,
				Init:     0,
				Pending:  0,
			},
			final_err: nil,
			kind:      ScanObjectTypeCommit,
		},
	}

	for _, test_i := range tests {
		t.Run(test_i.name, func(t *testing.T) {
			tracker, err := NewKeyTracker(test_i.kind, &logger)
			if test_i.final_err == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, test_i.final_err)
				return
			}

			// for each key in the test data, apply the series of update codes
			for _, d := range test_i.data {
				for _, code := range d.codes {
					_, update_err := tracker.Update(d.key, code, "", []string{})
					assert.NoError(t, update_err)
				}
			}

			// get the data for the key after applying all updates
			key_data_counts := tracker.GetCounts()
			assert.Equal(t, test_i.final, key_data_counts)
		})
	}
}

// TestKeyTracker_GetKeys unit test function tests the GetKeys() method of the KeyTracker type.
func TestKeyTracker_GetKeys(t *testing.T) {
	t.Parallel()

	// create a new KeyTracker instance
	tracker := &KeyTracker{
		Keys:   KeyDataMap{},
		Kind:   ScanObjectTypeFile,
		logger: nil,
		mu:     &sync.RWMutex{},
	}

	// add some test data to the tracker
	tracker.Keys["key1"] = KeyData{
		Code:  KeyCodeInit,
		State: KeyStateInit,
	}
	tracker.Keys["key2"] = KeyData{
		Code:  KeyCodeComplete,
		State: KeyStateComplete,
	}
	tracker.Keys["key3"] = KeyData{
		Code:  KeyCodeError,
		State: KeyStateError,
	}

	// call the GetKeys() method
	keys := tracker.GetKeys()

	// assert the expected keys
	expectedKeys := []string{"key1", "key2", "key3"}
	assert.ElementsMatch(t, expectedKeys, keys)
}

// TestKeyTracker_GetKeysData unit test function tests the GetKeysData() method of the KeyTracker type.
func TestKeyTracker_GetKeysData(t *testing.T) {
	t.Parallel()

	// create a new KeyTracker instance
	tracker := &KeyTracker{
		Keys:   KeyDataMap{},
		Kind:   ScanObjectTypeFile,
		logger: nil,
		mu:     &sync.RWMutex{},
	}

	// add some test data to the tracker
	test_key_init := "test_key_init"
	test_data_init := KeyData{
		Code:  KeyCodeInit,
		State: KeyStateInit,
	}
	tracker.Keys[test_key_init] = test_data_init
	test_key_error := "test_key_error"
	test_data_error := KeyData{
		Code:    KeyCodeError,
		Message: test_message_error,
		State:   KeyStateError,
	}
	tracker.Keys[test_key_error] = test_data_error
	test_key_ignore := "test_key_ignore"
	test_data_ignore := KeyData{
		Code:    KeyCodeIgnore,
		Message: test_message_ignore,
		State:   KeyStateIgnore,
	}
	tracker.Keys[test_key_ignore] = test_data_ignore
	test_key_pending := "test_key_pending"
	test_data_pending := KeyData{
		Code:    KeyCodePending,
		Message: test_message_pending,
		State:   KeyStatePending,
	}
	tracker.Keys[test_key_pending] = test_data_pending
	test_key_complete_1 := "test_key_complete_1"
	test_data_complete_1 := KeyData{
		Code:    KeyCodeComplete,
		Message: test_message_complete,
		State:   KeyStateComplete,
	}
	tracker.Keys[test_key_complete_1] = test_data_complete_1
	test_key_complete_2 := "test_key_complete_2"
	test_data_complete_2 := KeyData{
		Code:    KeyCodeComplete,
		Message: test_message_complete,
		State:   KeyStateComplete,
	}
	tracker.Keys[test_key_complete_2] = test_data_complete_2

	// call the GetKeysData method
	keysData := tracker.GetKeysData()
	// check if the returned map is equal to the original keys map
	assert.Equal(t, tracker.Keys, keysData)

	// get KeyData for each element where the Code is KeyCodeInit
	keys_data_init, init_err := tracker.GetKeysDataForCode(KeyCodeInit)
	assert.NoError(t, init_err)
	assert.Equal(t, 1, len(keys_data_init))
	assert.Equal(t, KeyDataMap{test_key_init: test_data_init}, keys_data_init)

	// get KeyData for each element where the Code is KeyCodeError
	keys_data_error, error_err := tracker.GetKeysDataForCode(KeyCodeError)
	assert.NoError(t, error_err)
	assert.Equal(t, 1, len(keys_data_error))
	assert.Equal(t, KeyDataMap{test_key_error: test_data_error}, keys_data_error)

	// get KeyData for each element where the Code is KeyCodeIgnore
	keys_data_ignore, ignore_err := tracker.GetKeysDataForCode(KeyCodeIgnore)
	assert.NoError(t, ignore_err)
	assert.Equal(t, 1, len(keys_data_ignore))
	assert.Equal(t, KeyDataMap{test_key_ignore: test_data_ignore}, keys_data_ignore)

	// get KeyData for each element where the Code is KeyCodePending
	keys_data_pending, pending_err := tracker.GetKeysDataForCode(KeyCodePending)
	assert.NoError(t, pending_err)
	assert.Equal(t, 1, len(keys_data_pending))
	assert.Equal(t, KeyDataMap{test_key_pending: test_data_pending}, keys_data_pending)

	// get KeyData for each element where the Code is KeyCodeComplete
	keys_data_complete, complete_err := tracker.GetKeysDataForCode(KeyCodeComplete)
	assert.NoError(t, complete_err)
	assert.Equal(t, 2, len(keys_data_complete))
	assert.Equal(
		t,
		KeyDataMap{
			test_key_complete_1: test_data_complete_1,
			test_key_complete_2: test_data_complete_2,
		},
		keys_data_complete,
	)

	_, err_low := tracker.GetKeysDataForCode(KeyCodeInit - 1)
	assert.ErrorContains(t, err_low, ErrKeyCodeInvalid.Error())

	_, err_high := tracker.GetKeysDataForCode(KeyCodeComplete + 1)
	assert.ErrorContains(t, err_high, ErrKeyCodeInvalid.Error())
}

// TestKeyTracker_PrintCodes unit test function tests the PrintCodes() method
// of the KeyTracker type.
func TestKeyTracker_PrintCodes(t *testing.T) {
	t.Parallel()

	logger := zerolog.New(os.Stdout)
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	// create a new KeyTracker instance
	tracker := &KeyTracker{
		Keys:   KeyDataMap{},
		Kind:   ScanObjectTypeFile,
		logger: &logger,
		mu:     &sync.RWMutex{},
	}

	// add some test data to the tracker
	tracker.Keys["key1"] = KeyData{
		Code:  KeyCodeInit,
		State: KeyStateInit,
	}
	tracker.Keys["key2"] = KeyData{
		Code:  KeyCodeComplete,
		State: KeyStateComplete,
	}
	tracker.Keys["key3"] = KeyData{
		Code:  KeyCodeError,
		State: KeyStateError,
	}

	expected_codes := []int{KeyCodeInit, KeyCodeComplete, KeyCodeError}

	// call the PrintCodes() methodV
	actual_codes := tracker.PrintCodes()
	for _, code := range actual_codes {
		//assert.Equal(t, expected_codes, actual_codes)
		assert.Contains(t, expected_codes, code)
	}
}

// TestKeyTracker_PrintCounts unit test function tests the PrintCounts() method
// of the KeyTracker type.
func TestKeyTracker_PrintCounts(t *testing.T) {
	t.Parallel()

	logger := zerolog.New(os.Stdout)
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	// create a new KeyTracker instance
	tracker := &KeyTracker{
		Keys:   KeyDataMap{},
		Kind:   ScanObjectTypeFile,
		logger: &logger,
		mu:     &sync.RWMutex{},
	}

	// add some test data to the tracker
	tracker.Keys["key1"] = KeyData{
		Code:  KeyCodeInit,
		State: KeyStateInit,
	}
	tracker.Keys["key2"] = KeyData{
		Code:  KeyCodeComplete,
		State: KeyStateComplete,
	}
	tracker.Keys["key3"] = KeyData{
		Code:  KeyCodeError,
		State: KeyStateError,
	}

	expected_key_data_counts := KeyDataCounts{
		Complete: 1,
		Error:    1,
		Ignore:   0,
		Init:     1,
		Pending:  0,
	}

	// call the PrintCounts() methodV
	actual_key_data_counts := tracker.PrintCounts()
	assert.Equal(t, expected_key_data_counts, actual_key_data_counts)
}

// TestKeyTracker_Restore unit test function tests the Restore() method
// of the KeyTracker type.
func TestKeyTracker_Restore(t *testing.T) {
	t.Parallel()

	// create a new KeyTracker instance
	trckr := &KeyTracker{
		Keys:   make(KeyDataMap),
		Kind:   ScanObjectTypeFile,
		logger: nil,
		mu:     &sync.RWMutex{},
	}

	// add some test data to the tracker
	test_key_init := "test_key_init"
	test_data_init := KeyData{
		Code:  KeyCodeInit,
		State: KeyStateInit,
	}
	trckr.Keys[test_key_init] = test_data_init
	test_key_error := "test_key_error"
	test_data_error := KeyData{
		Code:    KeyCodeError,
		Message: test_message_error,
		State:   KeyStateError,
	}
	trckr.Keys[test_key_error] = test_data_error
	test_key_ignore := "test_key_ignore"
	test_data_ignore := KeyData{
		Code:    KeyCodeIgnore,
		Message: test_message_ignore,
		State:   KeyStateIgnore,
	}
	trckr.Keys[test_key_ignore] = test_data_ignore
	test_key_pending := "test_key_pending"
	test_data_pending := KeyData{
		Code:    KeyCodePending,
		Message: test_message_pending,
		State:   KeyStatePending,
	}
	trckr.Keys[test_key_pending] = test_data_pending
	test_key_complete_1 := "test_key_complete_1"
	test_data_complete_1 := KeyData{
		Code:    KeyCodeComplete,
		Message: test_message_complete,
		State:   KeyStateComplete,
	}
	trckr.Keys[test_key_complete_1] = test_data_complete_1
	test_key_complete_2 := "test_key_complete_2"
	test_data_complete_2 := KeyData{
		Code:    KeyCodeComplete,
		Message: test_message_complete,
		State:   KeyStateComplete,
	}
	trckr.Keys[test_key_complete_2] = test_data_complete_2

	// call the GetKeysData method
	tracker_keys_data := trckr.GetKeysData()
	assert.Equal(t, len(trckr.Keys), len(tracker_keys_data))
	// check if the returned map is equal to the original keys map
	assert.Equal(t, trckr.Keys, tracker_keys_data)

	// create a new KeyTracker instance
	trckr_new := &KeyTracker{
		Keys:   make(KeyDataMap),
		Kind:   ScanObjectTypeFile,
		logger: nil,
		mu:     &sync.RWMutex{},
	}

	test_key_complete_3 := "test_key_complete_3"
	test_data_complete_3 := KeyData{
		Code:    KeyCodeComplete,
		Message: test_message_complete,
		State:   KeyStateComplete,
	}
	// write some new data to the new tracker before restore so that we can
	// check if it is correctly overwritten
	trckr_new.Keys[test_key_complete_3] = test_data_complete_3

	trckr_new.Restore(tracker_keys_data)
	assert.Equal(t, trckr.Kind, trckr_new.Kind)
	assert.Equal(t, len(trckr.Keys), len(trckr_new.Keys))
	assert.Equal(t, trckr.Keys, trckr_new.Keys)
}

// TestKeyTracker_Update() unit test function tests the Update() method
// of the KeyTracker type.
func TestKeyTracker_Update(t *testing.T) {
	t.Parallel()

	logger := zerolog.New(os.Stdout)
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	type testData struct {
		children        []string
		code            int
		expect_code     int
		expect_code_err error
		message         string
	}

	tests := []struct {
		name      string
		data      []testData
		final     KeyData
		init_err  error
		key_empty bool
		kind      string
	}{
		{
			name: "InvalidCodeLow",
			data: []testData{
				{
					children:        []string{},
					code:            KeyCodeInit - 1,
					expect_code:     0,
					expect_code_err: ErrKeyCodeInvalid,
					message:         "",
				},
			},
			final:     KeyData{},
			init_err:  nil,
			key_empty: false,
			kind:      ScanObjectTypeFile,
		},
		{
			name: "InvalidCodeHigh",
			data: []testData{
				{
					children:        []string{},
					code:            KeyCodeComplete + 1,
					expect_code:     0,
					expect_code_err: ErrKeyCodeInvalid,
					message:         "",
				},
			},
			final:     KeyData{},
			init_err:  nil,
			key_empty: false,
			kind:      ScanObjectTypeCommit,
		},
		{
			name: "ErrorOnEmptyKey",
			data: []testData{
				{
					children:        []string{},
					code:            KeyCodeComplete,
					expect_code:     0,
					expect_code_err: ErrKeyUpdateKeyEmpty,
					message:         "",
				},
			},
			final:     KeyData{},
			init_err:  nil,
			key_empty: true,
			kind:      ScanObjectTypeRequestResponse,
		},
		{
			name: "ErrorOnInvalidUpdateCodeLow",
			data: []testData{
				{
					children:        []string{},
					code:            KeyCodeInit - 1,
					expect_code:     0,
					expect_code_err: errors.New("failed to update data for key"),
					message:         "",
				},
			},
			final:     KeyData{},
			init_err:  nil,
			key_empty: false,
			kind:      ScanObjectTypeRequestResponse,
		},
		{
			name: "ErrorOnInvalidUpdateCodeHigh",
			data: []testData{
				{
					children:        []string{},
					code:            KeyCodeComplete + 1,
					expect_code:     0,
					expect_code_err: errors.New("failed to update data for key"),
					message:         "",
				},
			},
			final:     KeyData{},
			init_err:  nil,
			key_empty: false,
			kind:      ScanObjectTypeRequestResponse,
		},
		{
			name:      "InvalidKind",
			data:      []testData{},
			final:     KeyData{},
			init_err:  ErrKeyTrackerInvalidKind,
			key_empty: false,
			kind:      "InvalidKind",
		},
		{
			name: "CodeInit",
			data: []testData{
				{
					children:        []string{},
					code:            KeyCodeInit,
					expect_code:     KeyCodeInit,
					expect_code_err: nil,
					message:         "",
				},
			},
			final: KeyData{
				Children: map[string]bool{},
				Code:     KeyCodeInit,
				Message:  "",
				State:    KeyStateInit,
			},
			init_err:  nil,
			key_empty: false,
			kind:      ScanObjectTypeCommit,
		},
		{
			name: "CodeComplete",
			data: []testData{
				{
					children:        []string{},
					code:            KeyCodeInit,
					expect_code:     KeyCodeInit,
					expect_code_err: nil,
					message:         "",
				},
				{
					children:        []string{},
					code:            KeyCodeComplete,
					expect_code:     KeyCodeComplete,
					expect_code_err: nil,
					message:         "",
				},
			},
			final: KeyData{
				Children: map[string]bool{},
				Code:     KeyCodeComplete,
				Message:  "",
				State:    KeyStateComplete,
			},
			init_err:  nil,
			key_empty: false,
			kind:      ScanObjectTypeFile,
		},
		{
			name: "CodeError",
			data: []testData{
				{
					children:        []string{},
					code:            KeyCodeInit,
					expect_code:     KeyCodeInit,
					expect_code_err: nil,
					message:         "",
				},
				{
					children:        []string{},
					code:            KeyCodeError,
					expect_code:     KeyCodeError,
					expect_code_err: nil,
					message:         test_message_error,
				},
			},
			final: KeyData{
				Children: map[string]bool{},
				Code:     KeyCodeError,
				Message:  test_message_error,
				State:    KeyStateError,
			},
			init_err:  nil,
			key_empty: false,
			kind:      ScanObjectTypeFile,
		},
		{
			name: "Progression",
			data: []testData{
				{
					children:        []string{},
					code:            KeyCodeInit,
					expect_code:     KeyCodeInit,
					expect_code_err: nil,
					message:         "",
				},
				{
					children:        []string{},
					code:            KeyCodeError,
					expect_code:     KeyCodeError,
					expect_code_err: nil,
					message:         test_message_error,
				},
				{
					children:        []string{},
					code:            KeyCodeIgnore,
					expect_code:     KeyCodeIgnore,
					expect_code_err: nil,
					message:         test_message_ignore,
				},
				{
					children:        []string{"child1", "child2"},
					code:            KeyCodePending,
					expect_code:     KeyCodePending,
					expect_code_err: nil,
					message:         test_message_pending,
				},
				{
					children:        []string{"child1", "child2"},
					code:            KeyCodeComplete,
					expect_code:     KeyCodeComplete,
					expect_code_err: nil,
					message:         test_message_complete,
				},
			},
			final: KeyData{
				Children: map[string]bool{
					"child1": true,
					"child2": true,
				},
				Code:    KeyCodeComplete,
				Message: test_message_complete,
				State:   KeyStateComplete,
			},
			init_err:  nil,
			key_empty: false,
			kind:      ScanObjectTypeFile,
		},
		{
			name: "Regression",
			data: []testData{
				{
					children:        []string{},
					code:            KeyCodeComplete,
					expect_code:     KeyCodeComplete,
					expect_code_err: nil,
					message:         test_message_complete,
				},
				{
					children:        []string{},
					code:            KeyCodePending,
					expect_code:     KeyCodeComplete,
					expect_code_err: nil,
					message:         test_message_pending,
				},
				{
					children:        []string{},
					code:            KeyCodeIgnore,
					expect_code:     KeyCodeComplete,
					expect_code_err: nil,
					message:         test_message_ignore,
				},
				{
					children:        []string{},
					code:            KeyCodeError,
					expect_code:     KeyCodeComplete,
					expect_code_err: nil,
					message:         test_message_error,
				},
				{
					children:        []string{},
					code:            KeyCodeInit,
					expect_code:     KeyCodeComplete,
					expect_code_err: nil,
					message:         test_message_init,
				},
			},
			final: KeyData{
				Children: map[string]bool{},
				Code:     KeyCodeComplete,
				Message:  test_message_complete,
				State:    KeyStateComplete,
			},
			init_err:  nil,
			key_empty: false,
			kind:      ScanObjectTypeFile,
		},
		{
			name: "ReInit",
			data: []testData{
				{
					children:        []string{},
					code:            KeyCodeInit,
					expect_code:     KeyCodeInit,
					expect_code_err: nil,
					message:         "",
				},
				{
					children:        []string{},
					code:            KeyCodeComplete,
					expect_code:     KeyCodeComplete,
					expect_code_err: nil,
					message:         test_message_complete,
				},
				{
					children:        []string{},
					code:            KeyCodeInit,
					expect_code:     KeyCodeComplete,
					expect_code_err: nil,
					message:         test_message_init,
				},
			},
			final: KeyData{
				Children: map[string]bool{},
				Code:     KeyCodeComplete,
				Message:  test_message_complete,
				State:    KeyStateComplete,
			},
			init_err:  nil,
			key_empty: false,
			kind:      ScanObjectTypeFile,
		},
		{
			name: "RepeatUpdatePending",
			data: []testData{
				{
					children:        []string{"child1"},
					code:            KeyCodePending,
					expect_code:     KeyCodePending,
					expect_code_err: nil,
					message:         test_message_pending,
				},
				{
					children:        []string{"child2"},
					code:            KeyCodePending,
					expect_code:     KeyCodePending,
					expect_code_err: nil,
					message:         test_message_pending,
				},
				{
					children:        []string{"child3"},
					code:            KeyCodePending,
					expect_code:     KeyCodePending,
					expect_code_err: nil,
					message:         test_message_pending,
				},
				{
					children:        []string{"child4"},
					code:            KeyCodePending,
					expect_code:     KeyCodePending,
					expect_code_err: nil,
					message:         test_message_pending,
				},
				{
					children:        []string{"child5"},
					code:            KeyCodePending,
					expect_code:     KeyCodePending,
					expect_code_err: nil,
					message:         test_message_pending,
				},
				{
					children:        []string{"child6"},
					code:            KeyCodePending,
					expect_code:     KeyCodePending,
					expect_code_err: nil,
					message:         test_message_pending,
				},
				{
					children:        []string{"child7"},
					code:            KeyCodePending,
					expect_code:     KeyCodePending,
					expect_code_err: nil,
					message:         test_message_pending,
				},
				{
					children:        []string{"child8"},
					code:            KeyCodePending,
					expect_code:     KeyCodePending,
					expect_code_err: nil,
					message:         test_message_pending,
				},
				{
					children:        []string{"child9"},
					code:            KeyCodePending,
					expect_code:     KeyCodePending,
					expect_code_err: nil,
					message:         test_message_pending,
				},
				{
					children:        []string{"child10"},
					code:            KeyCodePending,
					expect_code:     KeyCodePending,
					expect_code_err: nil,
					message:         test_message_pending,
				},
			},
			final: KeyData{
				Children: map[string]bool{
					"child1":  false,
					"child2":  false,
					"child3":  false,
					"child4":  false,
					"child5":  false,
					"child6":  false,
					"child7":  false,
					"child8":  false,
					"child9":  false,
					"child10": false,
				},
				Code:    KeyCodePending,
				Message: test_message_pending,
				State:   KeyStatePending,
			},
			init_err:  nil,
			key_empty: false,
			kind:      ScanObjectTypeFile,
		},
		{
			name: "RepeatUpdateComplete",
			data: []testData{
				{
					children:        []string{"child1"},
					code:            KeyCodeComplete,
					expect_code:     KeyCodeComplete,
					expect_code_err: nil,
					message:         test_message_complete,
				},
				{
					children:        []string{"child2"},
					code:            KeyCodeComplete,
					expect_code:     KeyCodeComplete,
					expect_code_err: nil,
					message:         test_message_complete,
				},
				{
					children:        []string{"child3"},
					code:            KeyCodeComplete,
					expect_code:     KeyCodeComplete,
					expect_code_err: nil,
					message:         test_message_complete,
				},
				{
					children:        []string{"child4"},
					code:            KeyCodeComplete,
					expect_code:     KeyCodeComplete,
					expect_code_err: nil,
					message:         test_message_complete,
				},
				{
					children:        []string{"child5"},
					code:            KeyCodeComplete,
					expect_code:     KeyCodeComplete,
					expect_code_err: nil,
					message:         test_message_complete,
				},
				{
					children:        []string{"child6"},
					code:            KeyCodeComplete,
					expect_code:     KeyCodeComplete,
					expect_code_err: nil,
					message:         test_message_complete,
				},
				{
					children:        []string{"child7"},
					code:            KeyCodeComplete,
					expect_code:     KeyCodeComplete,
					expect_code_err: nil,
					message:         test_message_complete,
				},
				{
					children:        []string{"child8"},
					code:            KeyCodeComplete,
					expect_code:     KeyCodeComplete,
					expect_code_err: nil,
					message:         test_message_complete,
				},
				{
					children:        []string{"child9"},
					code:            KeyCodeComplete,
					expect_code:     KeyCodeComplete,
					expect_code_err: nil,
					message:         test_message_complete,
				},
				{
					children:        []string{"child10"},
					code:            KeyCodeComplete,
					expect_code:     KeyCodeComplete,
					expect_code_err: nil,
					message:         test_message_complete,
				},
			},
			final: KeyData{
				Children: map[string]bool{
					"child1":  true,
					"child2":  true,
					"child3":  true,
					"child4":  true,
					"child5":  true,
					"child6":  true,
					"child7":  true,
					"child8":  true,
					"child9":  true,
					"child10": true,
				},
				Code:    KeyCodeComplete,
				Message: test_message_complete,
				State:   KeyStateComplete,
			},
			init_err:  nil,
			key_empty: false,
			kind:      ScanObjectTypeFile,
		},
		{
			name: "IncompleteUpdates",
			data: []testData{
				{
					children: []string{
						"child1",
						"child2",
						"child3",
						"child4",
						"child5",
						"child6",
						"child7",
						"child8",
						"child9",
						"child10",
					},
					code:            KeyCodePending,
					expect_code:     KeyCodePending,
					expect_code_err: nil,
					message:         test_message_pending,
				},
				{
					children: []string{
						"child1",
						"child2",
						"child3",
						"child4",
						"child5",
					},
					code:            KeyCodeComplete,
					expect_code:     KeyCodePending,
					expect_code_err: nil,
					message:         "partial update",
				},
			},
			final: KeyData{
				Children: map[string]bool{
					"child1":  true,
					"child2":  true,
					"child3":  true,
					"child4":  true,
					"child5":  true,
					"child6":  false,
					"child7":  false,
					"child8":  false,
					"child9":  false,
					"child10": false,
				},
				Code:    KeyCodePending,
				Message: "partial update",
				State:   KeyStatePending,
			},
			init_err:  nil,
			key_empty: false,
			kind:      ScanObjectTypeFile,
		},
	}

	for _, test_i := range tests {
		t.Run(test_i.name, func(t *testing.T) {
			tracker, err := NewKeyTracker(test_i.kind, &logger)
			if test_i.init_err == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, test_i.init_err.Error())
				return
			}

			test_key := test_i.name
			if test_i.key_empty {
				test_key = ""
			}
			for _, d := range test_i.data {
				updated_code, update_err := tracker.Update(test_key, d.code, d.message, d.children)
				if d.expect_code_err != nil {
					assert.ErrorContains(t, update_err, d.expect_code_err.Error())
					return
				}
				assert.NoError(t, update_err)
				assert.Exactly(t, d.expect_code, updated_code, "updated code mismatch")
			}

			// get the data for the key after applying all updates
			key_data, key_exists := tracker.Get(test_i.name)
			if !assert.True(t, key_exists, "key should exist in tracker map") {
				t.FailNow()
			}
			assert.Equal(t, test_i.final.Code, key_data.Code)
			assert.Equal(t, test_i.final.Message, key_data.Message)
			assert.Equal(t, test_i.final.State, key_data.State)
			assert.Equal(t, len(test_i.final.Children), len(key_data.Children), "number of children mismatch")
			for child, _ := range test_i.final.Children {
				assert.Contains(t, key_data.Children, child)
				_, child_exists := key_data.Children[child]
				assert.True(t, child_exists)
			}
		})
	}
}

// TestKeyTracker_Concurrent_Update() unit test function tests concurrent data
// access via the Update() method of the KeyTracker type.
func TestKeyTracker_Concurrent_Update(t *testing.T) {
	t.Parallel()

	logger := zerolog.New(os.Stdout)
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	type testData struct {
		children    []string
		code        int
		expect_code int
		message     string
	}

	tests := []struct {
		name      string
		data      []testData
		final     KeyData
		final_err error
		kind      string
	}{
		{
			name: "RepeatUpdatePending",
			data: []testData{
				{
					children:    []string{"child1"},
					code:        KeyCodePending,
					expect_code: KeyCodePending,
					message:     test_message_pending,
				},
				{
					children:    []string{"child2"},
					code:        KeyCodePending,
					expect_code: KeyCodePending,
					message:     test_message_pending,
				},
				{
					children:    []string{"child3"},
					code:        KeyCodePending,
					expect_code: KeyCodePending,
					message:     test_message_pending,
				},
				{
					children:    []string{"child4"},
					code:        KeyCodePending,
					expect_code: KeyCodePending,
					message:     test_message_pending,
				},
				{
					children:    []string{"child5"},
					code:        KeyCodePending,
					expect_code: KeyCodePending,
					message:     test_message_pending,
				},
				{
					children:    []string{"child6"},
					code:        KeyCodePending,
					expect_code: KeyCodePending,
					message:     test_message_pending,
				},
				{
					children:    []string{"child7"},
					code:        KeyCodePending,
					expect_code: KeyCodePending,
					message:     test_message_pending,
				},
				{
					children:    []string{"child8"},
					code:        KeyCodePending,
					expect_code: KeyCodePending,
					message:     test_message_pending,
				},
				{
					children:    []string{"child9", "child10"},
					code:        KeyCodePending,
					expect_code: KeyCodePending,
					message:     test_message_pending,
				},
				{
					children:    []string{"child11", "child12", "child13", "child14", "child15"},
					code:        KeyCodePending,
					expect_code: KeyCodePending,
					message:     test_message_pending,
				},
				{
					children:    []string{"child16"},
					code:        KeyCodePending,
					expect_code: KeyCodePending,
					message:     test_message_pending,
				},
				{
					children:    []string{"child17"},
					code:        KeyCodePending,
					expect_code: KeyCodePending,
					message:     test_message_pending,
				},
				{
					children:    []string{"child18"},
					code:        KeyCodePending,
					expect_code: KeyCodePending,
					message:     test_message_pending,
				},
				{
					children:    []string{"child19"},
					code:        KeyCodePending,
					expect_code: KeyCodePending,
					message:     test_message_pending,
				},
				{
					children:    []string{"child20"},
					code:        KeyCodePending,
					expect_code: KeyCodePending,
					message:     test_message_pending,
				},
				{
					children:    []string{"child21"},
					code:        KeyCodePending,
					expect_code: KeyCodePending,
					message:     test_message_pending,
				},
				{
					children:    []string{"child22"},
					code:        KeyCodePending,
					expect_code: KeyCodePending,
					message:     test_message_pending,
				},
				{
					children:    []string{"child23"},
					code:        KeyCodePending,
					expect_code: KeyCodePending,
					message:     test_message_pending,
				},
				{
					children:    []string{"child24"},
					code:        KeyCodePending,
					expect_code: KeyCodePending,
					message:     test_message_pending,
				},
				{
					children:    []string{"child25"},
					code:        KeyCodePending,
					expect_code: KeyCodePending,
					message:     test_message_pending,
				},
				{
					children:    []string{"child26"},
					code:        KeyCodePending,
					expect_code: KeyCodePending,
					message:     test_message_pending,
				},
				{
					children:    []string{"child27"},
					code:        KeyCodePending,
					expect_code: KeyCodePending,
					message:     test_message_pending,
				},
				{
					children:    []string{"child28"},
					code:        KeyCodePending,
					expect_code: KeyCodePending,
					message:     test_message_pending,
				},
				{
					children:    []string{"child29"},
					code:        KeyCodePending,
					expect_code: KeyCodePending,
					message:     test_message_pending,
				},
				{
					children:    []string{"child30"},
					code:        KeyCodePending,
					expect_code: KeyCodePending,
					message:     test_message_pending,
				},
			},
			final: KeyData{
				Children: map[string]bool{
					"child1":  false,
					"child2":  false,
					"child3":  false,
					"child4":  false,
					"child5":  false,
					"child6":  false,
					"child7":  false,
					"child8":  false,
					"child9":  false,
					"child10": false,
					"child11": false,
					"child12": false,
					"child13": false,
					"child14": false,
					"child15": false,
					"child16": false,
					"child17": false,
					"child18": false,
					"child19": false,
					"child20": false,
					"child21": false,
					"child22": false,
					"child23": false,
					"child24": false,
					"child25": false,
					"child26": false,
					"child27": false,
					"child28": false,
					"child29": false,
					"child30": false,
				},
				Code:    KeyCodePending,
				Message: test_message_pending,
				State:   KeyStatePending,
			},
			final_err: nil,
			kind:      ScanObjectTypeFile,
		},
		{
			name: "RepeatUpdateComplete",
			data: []testData{
				{
					children:    []string{"child1"},
					code:        KeyCodeComplete,
					expect_code: KeyCodeComplete,
					message:     test_message_complete,
				},
				{
					children:    []string{"child2"},
					code:        KeyCodeComplete,
					expect_code: KeyCodeComplete,
					message:     test_message_complete,
				},
				{
					children:    []string{"child3"},
					code:        KeyCodeComplete,
					expect_code: KeyCodeComplete,
					message:     test_message_complete,
				},
				{
					children:    []string{"child4"},
					code:        KeyCodeComplete,
					expect_code: KeyCodeComplete,
					message:     test_message_complete,
				},
				{
					children:    []string{"child5"},
					code:        KeyCodeComplete,
					expect_code: KeyCodeComplete,
					message:     test_message_complete,
				},
				{
					children:    []string{"child6"},
					code:        KeyCodeComplete,
					expect_code: KeyCodeComplete,
					message:     test_message_complete,
				},
				{
					children:    []string{"child7"},
					code:        KeyCodeComplete,
					expect_code: KeyCodeComplete,
					message:     test_message_complete,
				},
				{
					children:    []string{"child8"},
					code:        KeyCodeComplete,
					expect_code: KeyCodeComplete,
					message:     test_message_complete,
				},
				{
					children:    []string{"child9", "child10"},
					code:        KeyCodeComplete,
					expect_code: KeyCodeComplete,
					message:     test_message_complete,
				},
				{
					children:    []string{"child11", "child12", "child13", "child14", "child15"},
					code:        KeyCodeComplete,
					expect_code: KeyCodeComplete,
					message:     test_message_complete,
				},
				{
					children:    []string{"child16"},
					code:        KeyCodeComplete,
					expect_code: KeyCodeComplete,
					message:     test_message_complete,
				},
				{
					children:    []string{"child17"},
					code:        KeyCodeComplete,
					expect_code: KeyCodeComplete,
					message:     test_message_complete,
				},
				{
					children:    []string{"child18"},
					code:        KeyCodeComplete,
					expect_code: KeyCodeComplete,
					message:     test_message_complete,
				},
				{
					children:    []string{"child19"},
					code:        KeyCodeComplete,
					expect_code: KeyCodeComplete,
					message:     test_message_complete,
				},
				{
					children:    []string{"child20"},
					code:        KeyCodeComplete,
					expect_code: KeyCodeComplete,
					message:     test_message_complete,
				},
				{
					children:    []string{"child21"},
					code:        KeyCodeComplete,
					expect_code: KeyCodeComplete,
					message:     test_message_complete,
				},
				{
					children:    []string{"child22"},
					code:        KeyCodeComplete,
					expect_code: KeyCodeComplete,
					message:     test_message_complete,
				},
				{
					children:    []string{"child23"},
					code:        KeyCodeComplete,
					expect_code: KeyCodeComplete,
					message:     test_message_complete,
				},
				{
					children:    []string{"child24"},
					code:        KeyCodeComplete,
					expect_code: KeyCodeComplete,
					message:     test_message_complete,
				},
				{
					children:    []string{"child25"},
					code:        KeyCodeComplete,
					expect_code: KeyCodeComplete,
					message:     test_message_complete,
				},
				{
					children:    []string{"child26"},
					code:        KeyCodeComplete,
					expect_code: KeyCodeComplete,
					message:     test_message_complete,
				},
				{
					children:    []string{"child27"},
					code:        KeyCodeComplete,
					expect_code: KeyCodeComplete,
					message:     test_message_complete,
				},
				{
					children:    []string{"child28"},
					code:        KeyCodeComplete,
					expect_code: KeyCodeComplete,
					message:     test_message_complete,
				},
				{
					children:    []string{"child29"},
					code:        KeyCodeComplete,
					expect_code: KeyCodeComplete,
					message:     test_message_complete,
				},
				{
					children:    []string{"child30"},
					code:        KeyCodeComplete,
					expect_code: KeyCodeComplete,
					message:     test_message_complete,
				},
			},
			final: KeyData{
				Children: map[string]bool{
					"child1":  true,
					"child2":  true,
					"child3":  true,
					"child4":  true,
					"child5":  true,
					"child6":  true,
					"child7":  true,
					"child8":  true,
					"child9":  true,
					"child10": true,
					"child11": true,
					"child12": true,
					"child13": true,
					"child14": true,
					"child15": true,
					"child16": true,
					"child17": true,
					"child18": true,
					"child19": true,
					"child20": true,
					"child21": true,
					"child22": true,
					"child23": true,
					"child24": true,
					"child25": true,
					"child26": true,
					"child27": true,
					"child28": true,
					"child29": true,
					"child30": true,
				},
				Code:    KeyCodeComplete,
				Message: test_message_complete,
				State:   KeyStateComplete,
			},
			final_err: nil,
			kind:      ScanObjectTypeFile,
		},
	}

	for _, test_i := range tests {
		t.Run(test_i.name, func(t *testing.T) {
			tracker, err := NewKeyTracker(test_i.kind, &logger)
			if test_i.final_err == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, test_i.final_err)
				return
			}

			updateFunc := func(wg *sync.WaitGroup, td testData) {
				defer wg.Done()
				updated_code, update_err := tracker.Update(test_i.name, td.code, td.message, td.children)
				assert.NoError(t, update_err)
				assert.Exactly(t, td.expect_code, updated_code)
			}
			wg := &sync.WaitGroup{}
			for _, d := range test_i.data {
				wg.Add(1)
				go updateFunc(wg, d)
			}
			wg.Wait()

			// get the data for the key after applying all updates
			key_data, key_exists := tracker.Get(test_i.name)
			if !assert.True(t, key_exists) {
				t.FailNow()
			}
			assert.Equal(t, test_i.final.Code, key_data.Code)
			assert.Equal(t, test_i.final.Message, key_data.Message)
			assert.Equal(t, test_i.final.State, key_data.State)
			assert.Equal(t, len(test_i.final.Children), len(key_data.Children), "number of children mismatch")
			for child, _ := range test_i.final.Children {
				assert.Contains(t, key_data.Children, child)
				_, child_exists := key_data.Children[child]
				assert.Truef(t, child_exists, "child %s does not exist", child)
			}
		})
	}
}
