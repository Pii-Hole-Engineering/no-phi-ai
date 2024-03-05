package az

import "time"

const DefaultDetectionApi string = "language/:analyze-text?api-version=2022-05-01"
const DefaultLanguage string = "en"
const DocumentCharacterLimit int = 5000
const RequestDocumentLimit int = 5
const RequestTimerDuration time.Duration = time.Second * 5
const ShowStatsParam string = "&showStats=true"
