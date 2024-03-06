package rrr

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"strings"
)

// MetadataRequestResponse struct contains metadata specific to a request and
// its associated response, where this metadata is typically copied from the
// request to the response.
type MetadataRequestResponse struct {
	// ID is the unique identifier of the request.
	ID string `json:"id"`
	// Commit struct contains information about the associated commit.
	Commit MetadataRequestResponseCommit `json:"commit"`
	// Object struct contains information about the associated object (e.g. file).
	Object MetadataRequestResponseObject `json:"object"`
	// Repository struct contains information about the associated repository.
	Repository MetadataRequestResponseRepository `json:"repository"`
	// Time struct contains timestamps set during request processing.
	Time MetadataRequestResponseTime `json:"time"`
}

type MetadataRequestResponseCommit struct {
	ID string `json:"id"`
}

type MetadataRequestResponseObject struct {
	// ID is the string version of the file's SHA1 hash, which is unique
	// to the file's content and context (e.g. repository, commit, etc.)
	ID string `json:"id"`
	// Length is the number of characters in the source text.
	Length int `json:"length"`
	// Offset is the starting character position of the source text within its
	// original context (e.g. offset fromn start of file).
	Offset int `json:"offset"`
}

type MetadataRequestResponseRepository struct {
	// ID is the unique identifier created by this app for the purpose of
	// tracking the repository.
	ID string `json:"id"`
	// URL is the unique URL used to interact with the repository.
	URL string `json:"url"`
}

type MetadataRequestResponseTime struct {
	// Start is the time the request was received.
	Start int64 `json:"start"`
	// Stop is the time the request was completed.
	Stop int64 `json:"stop"`
}

// Request struct contains all the information needed to process a request to
// detect PHI/PII data in some source (e.g. file) object and to identify the
// offending data within the source.
type Request struct {
	// embed the MetadataRequestResponse struct
	MetadataRequestResponse

	// Text is the source text to be scanned for PHI/PII data and is only
	// included in the Request (not the Response) object in order to limit the
	// size of the response and the exposure of the source text.
	Text string `json:"text"`
}

// NewRequestInput struct contains the input parameters required for the
// NewRequest() function.
type NewRequestInput struct {
	CommitID string
	Length   int
	ObjectID string
	Offset   int
	RepoID   string
	Text     string
}

// NewRequest() function initializes a new Request object.
func NewRequest(in NewRequestInput) (Request, error) {
	if in.RepoID == "" {
		return Request{}, ErrNewRequestEmptyRepositoryID
	}
	if in.ObjectID == "" {
		return Request{}, ErrNewRequestEmptyObjectID
	}
	if in.CommitID == "" {
		return Request{}, ErrNewRequestEmptyCommitID
	}
	if in.Text == "" {
		return Request{}, ErrNewRequestEmptyText
	}

	// create a unique ID by taking the SHA1 hash of the input strings
	// separated by the ResultSeparatorUID
	elements := []string{in.RepoID, in.CommitID, in.ObjectID, in.Text}
	sum := sha1.Sum([]byte(strings.Join(elements, ResultSeparatorUID)))

	return Request{
		MetadataRequestResponse: MetadataRequestResponse{
			ID: hex.EncodeToString(sum[:]),
			Commit: MetadataRequestResponseCommit{
				ID: in.CommitID,
			},
			Object: MetadataRequestResponseObject{
				ID:     in.ObjectID,
				Length: in.Length,
				Offset: in.Offset,
			},
			Repository: MetadataRequestResponseRepository{
				ID: in.RepoID,
			},
			Time: MetadataRequestResponseTime{
				Start: TimestampNow(),
				Stop:  0,
			},
		},
		Text: in.Text,
	}, nil
}

// Response struct embeds the Request struct and adds fields and methods
// specific to the response, such as the results from the detection services.
// A single Response may contain many results for the same scanned text.
type Response struct {
	// embed the MetadataRequestResponse struct
	MetadataRequestResponse
	// Results is a slice of detection results from the detection services.
	Results []Result `json:"results"`
}

// NewResponse() function initializes a new Response object from the provided
// Request.
func NewResponse(request *Request) Response {
	return Response{
		MetadataRequestResponse: request.MetadataRequestResponse,
		Results:                 make([]Result, 0),
	}
}

// RequestResponsePhiDetector interface defines the inputs of the
// Run() method, which is used to process Requests sent from
// the Scanner and to send Responses back to the Scanner.
type RequestResponsePhiDetector interface {
	Run(
		ctx context.Context,
		requests <-chan Request,
		responses chan<- Response,
	)
}
