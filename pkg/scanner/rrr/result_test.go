package rrr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestResult_Hash unit test function tests the Hash() method of the Result struct.
func TestResult_Hash(t *testing.T) {
	tests := []struct {
		commit_id       string
		expected_output string
		name            string
		object_id       string
		repo_id         string
		result          Result
	}{
		{
			commit_id:       "test_commit",
			expected_output: "da6544f2c9819e324b61b1e8de214dfe302fb969",
			name:            "TestResultHash1",
			object_id:       "test_object",
			repo_id:         "test_repo",
			result: Result{
				Category:        "test_category",
				ConfidenceScore: 0.99,
				Length:          1234,
				Offset:          500,
				Service:         "test_service",
				Subcategory:     "test_subcategory",
				Text:            "test_text",
			},
		},
		{
			commit_id:       "test_commit",
			expected_output: "b41daeb80879a4e20de48a4cff1be86da8818919",
			name:            "TestResultHash2",
			object_id:       "test_object",
			repo_id:         "test_repo",
			result: Result{
				Category:        "test_category",
				ConfidenceScore: 0.75,
				Length:          10,
				Offset:          20,
				Service:         "test_service",
				Subcategory:     "test_subcategory",
				Text:            "test_text",
			},
		},
		{
			commit_id:       "test_commit",
			expected_output: "ae8f05b895af05ae10bb9653a96c998217288e1f",
			name:            "TestResultHash3",
			object_id:       "test_object",
			repo_id:         "test_repo",
			result: Result{
				Category:        "test_category",
				ConfidenceScore: 0.01,
				Length:          12345678,
				Offset:          1500,
				Service:         "test_service",
				Subcategory:     "",
				Text:            "test_text",
			},
		},
		{
			commit_id:       "test_commit",
			expected_output: "e855e17cd2fdd810d2a61c26b3672cbade92899e",
			name:            "TestResultHash4",
			object_id:       "test_object",
			repo_id:         "test_repo",
			result: Result{
				Category:        "test_category",
				ConfidenceScore: 0.01,
				Length:          12345678,
				Offset:          1500,
				Service:         "test_service",
				Subcategory:     "test_subcategory",
				Text:            test_chunk_line_text_1,
			},
		},
		{
			commit_id:       "test_commit",
			expected_output: "0e063e56c2b6068bff7d6128e7d7371907594cc5",
			name:            "TestResultHash5",
			object_id:       "test_object",
			repo_id:         "test_repo",
			result: Result{
				Category:        "test_category",
				ConfidenceScore: 0.01,
				Length:          12345678,
				Offset:          1500,
				Service:         "test_service",
				Subcategory:     "test_subcategory",
				Text:            test_chunk_line_text_2,
			},
		},
		{
			commit_id:       "test_commit_alt",
			expected_output: "5c8836b87c1e15810b3bd878640e7c6c6e7528d2",
			name:            "TestResultHash6",
			object_id:       "test_object",
			repo_id:         "test_repo",
			result: Result{
				Category:        "test_category",
				ConfidenceScore: 0.01,
				Length:          12345678,
				Offset:          1500,
				Service:         "test_service",
				Subcategory:     "test_subcategory",
				Text:            test_chunk_line_text_2,
			},
		},
		{
			commit_id:       "test_commit",
			expected_output: "fc06598fecf210578e0e79257b6bbdb29ce7350e",
			name:            "TestResultHash7",
			object_id:       "test_object_alt",
			repo_id:         "test_repo",
			result: Result{
				Category:        "test_category",
				ConfidenceScore: 0.01,
				Length:          12345678,
				Offset:          1500,
				Service:         "test_service",
				Subcategory:     "test_subcategory",
				Text:            test_chunk_line_text_2,
			},
		},
		{
			commit_id:       "test_commit",
			expected_output: "492ee0e60c91fb35b498caeec3a3a81d2c8df1be",
			name:            "TestResultHash8",
			object_id:       "test_object",
			repo_id:         "test_repo_alt",
			result: Result{
				Category:        "test_category",
				ConfidenceScore: 0.01,
				Length:          12345678,
				Offset:          1500,
				Service:         "test_service",
				Subcategory:     "test_subcategory",
				Text:            test_chunk_line_text_2,
			},
		},
		{
			commit_id:       "test_commit",
			expected_output: "4d9b057d571d9547dd341791d2facfdb4c14c826",
			name:            "TestResultHash9",
			object_id:       "test_object",
			repo_id:         "test_repo",
			result: Result{
				Category:        "test_category_alt",
				ConfidenceScore: 0.01,
				Length:          12345678,
				Offset:          1500,
				Service:         "test_service",
				Subcategory:     "test_subcategory",
				Text:            test_chunk_line_text_2,
			},
		},
		{
			commit_id:       "test_commit",
			expected_output: "75d3b20b1ab7b3dc2ee16df865f6881ff400ff57",
			name:            "TestResultHash10",
			object_id:       "test_object",
			repo_id:         "test_repo",
			result: Result{
				Category:        "test_category",
				ConfidenceScore: 0.01,
				Length:          12345678,
				Offset:          1500,
				Service:         "test_service_alt",
				Subcategory:     "test_subcategory",
				Text:            test_chunk_line_text_2,
			},
		},
		{
			commit_id:       "test_commit",
			expected_output: "f47543b4a80840046657832284b1823e737e66c4",
			name:            "TestResultHash11",
			object_id:       "test_object",
			repo_id:         "test_repo",
			result: Result{
				Category:        "test_category",
				ConfidenceScore: 0.01,
				Length:          12345678,
				Offset:          1500,
				Service:         "test_service",
				Subcategory:     "test_subcategory_alt",
				Text:            test_chunk_line_text_2,
			},
		},
		{
			commit_id:       "",
			expected_output: "d20f36ec299a4d84b03e7b3c7afcc9789ca12c31",
			name:            "TestResultHashEmptyIDs",
			object_id:       "",
			repo_id:         "",
			result: Result{
				Category:        "test_category",
				ConfidenceScore: 0.75,
				Length:          10,
				Offset:          20,
				Service:         "test_service",
				Subcategory:     "test_subcategory",
				Text:            "test_text",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := test.result.Hash(test.repo_id, test.commit_id, test.object_id)

			assert.Equal(t, test.expected_output, output)
		})
	}
}

// TestResult_String unit test function tests the String() method of the Result
// struct.
func TestResult_String(t *testing.T) {
	tests := []struct {
		commit_id       string
		expected_output string
		name            string
		object_id       string
		repo_id         string
		result          Result
	}{
		{
			commit_id:       "test_commit",
			expected_output: "test_repo__test_commit__test_object__test_category__0.75__10__20__test_service__test_subcategory__test_text",
			name:            "TestResultString1",
			object_id:       "test_object",
			repo_id:         "test_repo",
			result: Result{
				Category:        "test_category",
				ConfidenceScore: 0.75,
				Length:          10,
				Offset:          20,
				Service:         "test_service",
				Subcategory:     "test_subcategory",
				Text:            "test_text",
			},
		},
		{
			commit_id:       "",
			expected_output: "###@@@###__###@@@###__###@@@###__test_category__0.75__10__20__test_service__test_subcategory__test_text",
			name:            "TestResultStringEmptyIDs",
			object_id:       "",
			repo_id:         "",
			result: Result{
				Category:        "test_category",
				ConfidenceScore: 0.75,
				Length:          10,
				Offset:          20,
				Service:         "test_service",
				Subcategory:     "test_subcategory",
				Text:            "test_text",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := test.result.String(test.repo_id, test.commit_id, test.object_id)
			assert.Equal(t, test.expected_output, output)
		})
	}
}

// TestResultRecordsFromResponse unit test function tests the
// ResultRecordsFromResponse() method of the Result struct.
func TestResultRecordsFromResponse(t *testing.T) {
	meta_req_resp := MetadataRequestResponse{
		Commit: MetadataRequestResponseCommit{
			ID: "test_commit",
		},
		Object: MetadataRequestResponseObject{
			ID: "test_object",
		},
		Repository: MetadataRequestResponseRepository{
			ID: "test_repo",
		},
	}
	resp := &Response{
		Results: []Result{
			{
				Category:        "test_category",
				ConfidenceScore: 0.99,
				Length:          1234,
				Offset:          500,
				Service:         "test_service",
				Subcategory:     "test_subcategory",
				Text:            "test_text",
			},
			{
				Category:        "test_category",
				ConfidenceScore: 0.75,
				Length:          10,
				Offset:          20,
				Service:         "test_service",
				Subcategory:     "test_subcategory",
				Text:            "test_text",
			},
		},
		MetadataRequestResponse: meta_req_resp,
	}

	expectedRecords := []ResultRecord{
		{
			Hash:                    "da6544f2c9819e324b61b1e8de214dfe302fb969",
			MetadataRequestResponse: meta_req_resp,
			Result: Result{
				Category:        "test_category",
				ConfidenceScore: 0.99,
				Length:          1234,
				Offset:          500,
				Service:         "test_service",
				Subcategory:     "test_subcategory",
				Text:            "test_text",
			},
		},
		{
			Hash:                    "b41daeb80879a4e20de48a4cff1be86da8818919",
			MetadataRequestResponse: meta_req_resp,
			Result: Result{
				Category:        "test_category",
				ConfidenceScore: 0.75,
				Length:          10,
				Offset:          20,
				Service:         "test_service",
				Subcategory:     "test_subcategory",
				Text:            "test_text",
			},
		},
	}

	records := ResultRecordsFromResponse(resp)

	for _, record := range records {
		if !assert.Contains(t, expectedRecords, record) {
			t.FailNow()
		}
		for _, expectedRecord := range expectedRecords {
			if record.Hash == expectedRecord.Hash {
				assert.Equal(t, expectedRecord.MetadataRequestResponse, record.MetadataRequestResponse)
				assert.Equal(t, expectedRecord.Result, record.Result)
			}
		}
	}
}
