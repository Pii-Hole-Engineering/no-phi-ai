package rrr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewRequest() unit test function tests the NewRequest() function.
func TestNewRequest(t *testing.T) {
	t.Parallel()

	repoID := "test-repo"
	commitID := "test-commit"
	objectID := "test-object"
	text := "test-text"

	t.Run("ValidInput", func(t *testing.T) {
		request, err := NewRequest(repoID, commitID, objectID, text)

		assert.NoError(t, err)
		assert.NotEmpty(t, request.ID)
		assert.Equal(t, commitID, request.Commit.ID)
		assert.Equal(t, objectID, request.Object.ID)
		assert.Equal(t, repoID, request.Repository.ID)
		assert.NotZero(t, request.Time.Start)
		assert.Zero(t, request.Time.Stop)
		assert.Equal(t, text, request.Text)
	})

	t.Run("EmptyRepositoryID", func(t *testing.T) {
		request, err := NewRequest("", commitID, objectID, text)

		assert.Error(t, err)
		assert.Equal(t, ErrNewRequestEmptyRepositoryID, err)
		assert.Empty(t, request.ID)
	})

	t.Run("EmptyCommitID", func(t *testing.T) {
		request, err := NewRequest(repoID, "", objectID, text)

		assert.Error(t, err)
		assert.Equal(t, ErrNewRequestEmptyCommitID, err)
		assert.Empty(t, request.ID)
	})

	t.Run("EmptyObjectID", func(t *testing.T) {
		request, err := NewRequest(repoID, commitID, "", text)

		assert.Error(t, err)
		assert.Equal(t, ErrNewRequestEmptyObjectID, err)
		assert.Empty(t, request.ID)
	})

	t.Run("EmptyText", func(t *testing.T) {
		request, err := NewRequest(repoID, commitID, objectID, "")

		assert.Error(t, err)
		assert.Equal(t, ErrNewRequestEmptyText, err)
		assert.Empty(t, request.ID)
	})
}

// TestNewResponse() unit test function tests the NewResponse() function.
func TestNewResponse(t *testing.T) {
	t.Parallel()

	// create a sample request
	request := &Request{
		MetadataRequestResponse: MetadataRequestResponse{
			ID: "test-id",
		},
	}

	t.Run("EmptyResults", func(t *testing.T) {
		response := NewResponse(request)

		assert.Equal(t, request.MetadataRequestResponse, response.MetadataRequestResponse)
		assert.Empty(t, response.Results)
	})

	t.Run("NonEmptyResults", func(t *testing.T) {
		// Create a sample result
		result := Result{
			Category:        "test-category",
			ConfidenceScore: 0.5,
			Length:          10,
			Offset:          0,
			Service:         "test-service",
			Subcategory:     "test-subcategory",
			Text:            "test-text",
		}

		response := NewResponse(request)
		response.Results = append(response.Results, result)

		assert.Equal(t, request.MetadataRequestResponse, response.MetadataRequestResponse)
		assert.Equal(t, []Result{result}, response.Results)
	})
}
