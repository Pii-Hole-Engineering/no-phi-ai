package nogit

import (
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

// ParseOrgNameFromURL() function is used to parse the name of the owner
// from a GitHub repository URL.
func ParseOrgNameFromURL(url_in string) (string, error) {
	// get the path elements from the URL
	path_elements, err := parseElementsFromURL(url_in)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse org name")
	}

	// the owner name should be the first element in the path
	owner_name := path_elements[0]
	if owner_name == "" {
		return "", errors.New("failed to parse org name : invalid path in URL = " + url_in)
	}
	if strings.HasSuffix(owner_name, ".git") {
		return "", errors.New("failed to parse org name : invalid owner name in URL = " + url_in)
	}

	return owner_name, nil
}

// ParseRepoNameFromURL() function is used to parse the repository name
// from a GitHub repository URL.
func ParseRepoNameFromURL(url_in string) (string, error) {
	// get the path elements from the URL
	path_elements, err := parseElementsFromURL(url_in)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse repo name")
	}
	if len(path_elements) < 2 {
		return "", errors.New("failed to parse repo name : invalid path in URL = " + url_in)
	}

	// the repo name should be the second element in the path
	repo_name := path_elements[1]
	// trim any ".git" suffix from the repo name
	repo_name = strings.TrimSuffix(repo_name, ".git")

	return repo_name, nil
}

func convertGitToURL(in string) (out string) {
	trim_string := "git@github.com:"
	if strings.HasPrefix(in, trim_string) {
		out = "https://github.com/" + strings.TrimPrefix(in, trim_string)
	} else {
		out = in
	}

	return
}

// parseElementsFromURL() function is used to parse the elements from the
// input URL into a slice of strings. Returns a non-nil error if the URL
// is invalid (for the purposes of this app) or the path is empty.
func parseElementsFromURL(url_in string) ([]string, error) {
	parsed_url, parse_err := url.Parse(convertGitToURL(url_in))
	if parse_err != nil {
		return []string{}, parse_err
	}

	if parsed_url.Path == "" {
		return []string{}, errors.New("failed parsing due to empty path in URL = " + url_in)
	}

	// trim any leading or trailing slashes and split the path into elements
	path_elements := strings.Split(strings.Trim(parsed_url.Path, "/"), "/")
	if len(path_elements) < 1 {
		return []string{}, errors.New("failed to parse elements from URL = " + url_in)
	}

	return path_elements, nil
}
