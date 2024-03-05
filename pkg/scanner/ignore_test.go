package scanner

import (
	"testing"
)

func TestIgnoreFilePath(t *testing.T) {
	tests := []struct {
		path   string
		ignore bool
		reason string
	}{
		{
			path:   "/full/path/to/file.txt",
			ignore: false,
			reason: "",
		},
		{
			path:   "relative/path/to/file.txt",
			ignore: false,
			reason: "",
		},
		{
			path:   "LOCK",
			ignore: true,
			reason: IgnoreReasonFileName,
		},
		{
			path:   "vendor/path/to/ignored_file.txt",
			ignore: true,
			reason: IgnoreReasonDirPath,
		},
	}

	for _, test := range tests {
		ignore, reason := IgnoreFilePath(test.path)
		if ignore != test.ignore {
			t.Errorf("IgnoreFilePath(%q) returned ignore=%v, want %v", test.path, ignore, test.ignore)
		}
		if reason != test.reason {
			t.Errorf("IgnoreFilePath(%q) returned reason=%q, want %q", test.path, reason, test.reason)
		}
	}
}
