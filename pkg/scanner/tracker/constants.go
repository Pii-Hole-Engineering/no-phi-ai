package tracker

const KeyCodeComplete = 2
const KeyCodeIgnore = 0
const KeyCodeInit = -2
const KeyCodeError = -1
const KeyCodePending = 1

const KeyStateComplete = "complete"
const KeyStateIgnore = "ignore"
const KeyStateInit = "init"
const KeyStateError = "error"
const KeyStatePending = "pending"

const ScanObjectTypeCommit string = "commit"
const ScanObjectTypeFile string = "file"
const ScanObjectTypeRequestResponse string = "request_and_response"
