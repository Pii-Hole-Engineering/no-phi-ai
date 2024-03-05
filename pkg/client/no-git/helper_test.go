package nogit

import (
	"errors"
	"testing"
)

func TestParseOrgNameFromURL(t *testing.T) {
	tests := []struct {
		url       string
		orgName   string
		expectErr bool
	}{
		{
			url:       "https://github.com/example-org/repo?some-random-query?true",
			orgName:   "example-org",
			expectErr: false,
		},
		{
			url:       "https://github.com/example-org/repo.git",
			orgName:   "example-org",
			expectErr: false,
		},
		{
			url:       "https://github.com/example-org/repo",
			orgName:   "example-org",
			expectErr: false,
		},
		{
			url:       "https://github.com/example-org/",
			orgName:   "example-org",
			expectErr: false,
		},
		{
			url:       "https://github.com/repo.git",
			orgName:   "",
			expectErr: true,
		},
		{
			url:       "https://github.com/",
			orgName:   "",
			expectErr: true,
		},
	}

	for _, test := range tests {
		orgName, err := ParseOrgNameFromURL(test.url)
		if test.expectErr && err == nil {
			t.Errorf("Expected error for url '%s', but got no error", test.url)
		}
		if !test.expectErr && err != nil {
			t.Errorf("Unexpected error for url '%s': %s", test.url, err)
		}
		if orgName != test.orgName {
			t.Errorf("Incorrect org name for url '%s'. Expected '%s', but got '%s'", test.url, test.orgName, orgName)
		}
	}
}

func TestParseRepoNameFromURL(t *testing.T) {
	tests := []struct {
		url      string
		expected string
		err      error
	}{
		{
			url:      "git@github.com:example-org/repo.git",
			expected: "repo",
			err:      nil,
		},
		{
			url:      "https://github.com/example-org/repo?some-random-query?true",
			expected: "repo",
			err:      nil,
		},
		{
			url:      "https://github.com/username/repo.git",
			expected: "repo",
			err:      nil,
		},
		{
			url:      "https://github.com/username/repo",
			expected: "repo",
			err:      nil,
		},
		{
			url:      "https://github.com/repo.git",
			expected: "",
			err:      errors.New("failed to parse repo name : invalid path in URL = https://github.com/repo.git"),
		},
		{
			url:      "https://github.com/repo",
			expected: "",
			err:      errors.New("failed to parse repo name : invalid path in URL = https://github.com/repo"),
		},
		{
			url:      "https://github.com/",
			expected: "",
			err:      errors.New("failed to parse repo name : invalid path in URL = https://github.com/"),
		},
		{
			url:      "https://github.com/repo/",
			expected: "",
			err:      errors.New("failed to parse repo name : invalid path in URL = https://github.com/repo/"),
		},
	}

	for _, test := range tests {
		actual, err := ParseRepoNameFromURL(test.url)
		if actual != test.expected {
			t.Errorf("ParseRepoNameFromURL(%s) = %s, expected %s", test.url, actual, test.expected)
		}
		if (err == nil && test.err != nil) || (err != nil && test.err == nil) || (err != nil && test.err != nil && err.Error() != test.err.Error()) {
			t.Errorf("ParseRepoNameFromURL(%s) error = %v, expected %v", test.url, err, test.err)
		}
	}
}
