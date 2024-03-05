package rrr

import (
	"bufio"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/pkg/errors"
)

// ChunkFileInput struct contains the input parameters required for
// the ChunkFileToRequests() function.
type ChunkFileInput struct {
	CommitID     string
	File         *object.File
	MaxChunkSize int
	RepoID       string
}

// ChunkFileToRequests() function reads the input object.File and
// generates a slice of requests, where the text in each request is
// limited to MaxChunkSize characters.
func ChunkFileToRequests(in ChunkFileInput) (requests []Request, e error) {
	if in.File == nil {
		e = ErrChunkFileToRequestsInFileNil
		return
	}

	requests = make([]Request, 0)
	file_reader, err := in.File.Reader()
	if err != nil {
		e = errors.Wrapf(err, ErrMsgScanFileRequestsGenerate, in.File.Hash.String())
		return
	}

	file_scanner := bufio.NewScanner(file_reader)
	file_scanner.Split(bufio.ScanLines)

	var current_offset int
	var current_text string

	for file_scanner.Scan() {
		line := file_scanner.Text()
		next_length := len(line) + len(current_text) + len("\n")

		if next_length < in.MaxChunkSize {
			current_text += "\n" + line
			continue
		}

		// create a request from the current_text
		request, err := NewRequest(
			in.RepoID,
			in.CommitID,
			in.File.ID().String(),
			current_text, // create Request from current_text
		)
		if err != nil {
			e = errors.Wrap(err, ErrMsgScanFileRequestsGenerate)
			return
		}
		requests = append(requests, request)
		// increment the current_offset by the length of the current_text
		current_offset += len(current_text)
		// reset the current_text variable to prepare for the next
		// chunk of the file
		current_text = ""

		if len(line) < in.MaxChunkSize {
			current_text = line
			continue
		}
		// chunk the line into smaller pieces of MaxChunkSize
		var line_requests []Request
		var req_err error
		// pass the current_offset to the ChunkLineToRequests() function
		// and get back the updated current_offset, a slice of requests
		// generated just for this line, and any error in generating the
		// requests
		current_offset, line_requests, req_err = ChunkLineToRequests(ChunkLineInput{
			CommitID:     in.CommitID,
			Line:         line,
			MaxChunkSize: in.MaxChunkSize,
			ObjectID:     in.File.ID().String(),
			Offset:       current_offset,
			RepoID:       in.RepoID,
		})
		if req_err != nil {
			e = errors.Wrap(req_err, ErrMsgScanFileRequestsGenerate)
			return
		}
		requests = append(requests, line_requests...)
		continue
	}

	// ensure that the last chunk of the file is included in the requests
	if current_text != "" {
		// create a new request for the remaining text
		request, err := NewRequest(
			in.RepoID,
			in.CommitID,
			in.File.ID().String(),
			current_text,
		)
		if err != nil {
			e = err
			return
		}
		requests = append(requests, request)
	}

	// validate that the chunking process produced requests if the file
	// has a size greater than 0
	if in.File.Size > 0 && len(requests) == 0 {
		e = ErrChunkFileToRequestsFailed
		return
	}

	return
}

// ChunkLineInput struct contains the input parameters required for the
// ChunkLineToRequests() function.
type ChunkLineInput struct {
	CommitID     string
	Line         string
	MaxChunkSize int
	ObjectID     string
	Offset       int
	RepoID       string
}

// ChunkLineToRequests() function chunks the input line (string) of text
// into smaller pieces of MaxChunkSize and sends requests to the channel
// for processing.
func ChunkLineToRequests(in ChunkLineInput) (offset int, requests []Request, e error) {
	offset = in.Offset

	if in.Line == "" {
		return
	}
	if in.MaxChunkSize <= 0 {
		e = ErrMaxChunkSizeInvalid
		return
	}

	line_reader := strings.NewReader(in.Line)
	// scan the words of the line
	line_scanner := bufio.NewScanner(line_reader)
	line_scanner.Split(bufio.ScanWords)

	var current_text string

	for line_scanner.Scan() {
		word := line_scanner.Text()
		next_length := len(current_text) + len(" ") + len(word)

		if next_length < in.MaxChunkSize {
			if current_text == "" {
				current_text = word
			} else {
				current_text += (" " + word)
			}
			continue
		}

		// create a new request when current text is within a word of
		// the MaxChunkSize
		request, err := NewRequest(
			in.RepoID,
			in.CommitID,
			in.ObjectID,
			current_text,
		)
		if err != nil {
			e = err
			return
		}
		requests = append(requests, request)
		// increment the offset by the length of the current_text
		offset += len(current_text) + len(" ")
		// reset the current_text to the value of the current word
		current_text = word
	}

	if current_text != "" {
		// create a new request for the remaining text
		request, err := NewRequest(
			in.RepoID,
			in.CommitID,
			in.ObjectID,
			current_text,
		)
		if err != nil {
			e = err
			return
		}
		requests = append(requests, request)
		// increment the offset by the length of the current_text
		offset += len(current_text)
	}

	return
}
