package scanner

import (
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/object"
)

// IgnorePaths is a map of paths that should not be scanned.
var IgnorePaths = map[string]bool{
	".git":   true,
	"vendor": true,
}

// IgnoreFilenames is a map of file names that should not be scanned.
var IgnoreFilenames = map[string]bool{
	"":           true,
	".gitignore": true,
	"LOCK":       true,
}

// IgnoreExtensions is a map of file extensions that should not be scanned.
var IgnoreExtensions = map[string]bool{
	".":    true,
	".jpg": true,
	".png": true,
}

// IgnoreFileObject() function can be used to check whether a file should be
// ignored (i.e. not scanned) for any reason, such as:
//   - file extension must be allowed by user-provided config (i.e. default ignore)
//   - file extension cannot be forbidden by app policy
//   - binary files are always ignored
//   - empty files are always ignored
//
// Returns a boolean to indicate whether the file object should be ignored, along
// with a reason (string) if ignore=true.
func IgnoreFileObject(file *object.File, supported_extensions, ignored_extensions []string) (ignore bool, reason string) {
	// ignore binary files
	if is_binary, _ := file.IsBinary(); is_binary {
		ignore = true
		reason = IgnoreReasonFileIsBinary
		return
	}

	// ignore empty files
	if file.Size == 0 {
		ignore = true
		reason = IgnoreReasonFileIsEmpty
		return
	}

	// check if the file path should be ignored by (app) policy
	ignore, reason = IgnoreFilePath(file.Name)
	if ignore && reason != "" {
		return
	}

	// get the file extension from the path
	file_extension := filepath.Ext(file.Name)
	// check against each file extension that is explicitly ignored by (user) config
	for _, ignored_extension := range ignored_extensions {
		if file_extension == ignored_extension {
			ignore = true
			reason = IgnoreReasonFileExtensionIgnoredByConfig
			return
		}
	}

	for _, ext := range supported_extensions {
		if file_extension == ext {
			ignore = false
			reason = ""
			return
		}
	}

	// ignore by default in order to avoid generating false positives for
	// files that are not well supported by existing models and/or highly
	// unlikely to contain PHI/PII data.
	ignore = true
	reason = IgnoreReasonDefault

	return
}

// IgnoreFilePath() function checks if a file path should be ignored
// (i.e. not scanned) based on the file extansion and name. Returns
// ignore = true if the file should be ignored, and ignore = false
// if the file should be scanned. Also returns a reason string that
// explains the justification for ignoring the file, if applicable.
func IgnoreFilePath(path string) (ignore bool, reason string) {
	// explicitly set defaults for return values
	ignore = false
	reason = ""

	// get the file name from the path
	file_name := filepath.Base(path)
	// check if the filename is in the IgnoreFilenames map
	if ignore_file_name, exists := IgnoreFilenames[file_name]; exists && ignore_file_name {
		// ignore the file / do not scan
		ignore = true
		reason = IgnoreReasonFileName
		return
	}

	// get the file extension from the path
	file_extension := filepath.Ext(path)
	if ignore_file_extension, exists := IgnoreExtensions[file_extension]; exists && ignore_file_extension {
		// ignore the file / do not scan
		ignore = true
		reason = IgnoreReasonFileExtensionIgnoredByPolicy
		return
	}

	return ignorePath(path)
}

// ignorePath() function checks if a (file or directory) path should be ignored
// (i.e. not scanned) based on the path itself, the directory (parent) of the
// input path, or the top-level directory of the input path. Returns ignore=true
// if the path should be ignored, and ignore=false if the path should be scanned.
func ignorePath(path string) (ignore bool, reason string) {
	// explicitly set defaults for return values
	ignore = false
	reason = ""

	// check if the path itself is in the IgnorePaths map
	if ignore_file_path, exists := IgnorePaths[path]; exists && ignore_file_path {
		ignore = true
		reason = IgnoreReasonFilePath
		return
	}

	// trim the "/" prefix from the input path as needed
	path = strings.TrimPrefix(path, "/")

	// iterate through each directory parent of the original path
	for dir_path := filepath.Dir(path); dir_path != "." && dir_path != "/"; dir_path = filepath.Dir(dir_path) {
		// check if the directory is an explicitly ignored path
		if ignore_dir, exists := IgnorePaths[dir_path]; exists && ignore_dir {
			ignore = true
			reason = IgnoreReasonDirPath
			return
		}
	}

	return
}
