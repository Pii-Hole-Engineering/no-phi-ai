package scanner

import (
	"testing"

	"github.com/Pii-Hole-Engineering/no-phi-ai/pkg/cfg"
	fixtures "github.com/go-git/go-git-fixtures/v4"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/filemode"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/stretchr/testify/assert"
)

// TestIgnoreFileObject() unit test function is used to test the
// IgnoreFileObject() function.
func TestIgnoreFileObject(t *testing.T) {
	t.Parallel()

	fileObjectFuncFooBlob := func(repo, commit, path string) *object.File {
		o := &plumbing.MemoryObject{}
		o.SetType(plumbing.BlobObject)
		o.SetSize(3)

		writer, err := o.Writer()
		assert.NoError(t, err)
		defer func() { assert.NoError(t, writer.Close()) }()

		writer.Write([]byte{'F', 'O', 'O'})

		blob := &object.Blob{}
		blob.Decode(o)
		return object.NewFile(path, filemode.Regular, blob)
	}

	fileObjectFuncFixture := func(repo, commit, path string) *object.File {
		f := fixtures.ByURL(repo).One()
		sto := filesystem.NewStorage(f.DotGit(), cache.NewObjectLRUDefault())

		h := plumbing.NewHash(commit)
		commit_object, err := object.GetCommit(sto, h)
		if !assert.NoError(t, err) {
			t.Logf("get commit error=%s : test.repo=%s : test=commit=%s : test.path=%s", err.Error(), repo, commit, path)
			return nil
		}

		file, err := commit_object.File(path)
		if !assert.NoError(t, err) {
			t.Logf("get commit file error=%s : test.repo=%s : test=commit=%s : test.path=%s", err.Error(), repo, commit, path)
			return nil
		}
		return file
	}

	tests := []struct {
		commit               string // the commit to search for the file
		extensions_ignored   []string
		extensions_supported []string
		fileObjectFunc       func(repo, commit, path string) *object.File
		ignore               bool
		lines                []string // expected lines in the file
		name                 string
		path                 string // the path of the file to find
		reason               string
		repo                 string // the repo name as in localRepos
	}{
		{
			commit:               "",
			extensions_ignored:   []string{},
			extensions_supported: cfg.DefaultScanFileExtensions,
			fileObjectFunc: func(repo, commit, path string) *object.File {
				return nil
			},
			ignore: true,
			lines:  []string{},
			name:   "FileObjectPointerNil",
			path:   "",
			reason: IgnoreReasonFileObjectPointerNil,
			repo:   "",
		},
		{
			commit:               "",
			extensions_ignored:   []string{},
			extensions_supported: cfg.DefaultScanFileExtensions,
			fileObjectFunc: func(repo, commit, path string) *object.File {
				o := &plumbing.MemoryObject{}
				o.SetType(plumbing.BlobObject)
				o.SetSize(0)

				blob := &object.Blob{}
				blob.Decode(o)
				return object.NewFile(path, filemode.Regular, blob)
			},
			ignore: true,
			lines:  []string{},
			name:   "FileObjectZeroSize",
			path:   "test.json",
			reason: IgnoreReasonFileIsEmpty,
			repo:   "",
		},
		{
			commit:               "",
			extensions_ignored:   []string{".json"},
			extensions_supported: cfg.DefaultScanFileExtensions,
			fileObjectFunc:       fileObjectFuncFooBlob,
			ignore:               true,
			lines: []string{
				"FOO",
			},
			name:   "IgnoreExtensionsByConfig",
			path:   "test.json",
			reason: IgnoreReasonFileExtensionIgnoredByConfig,
			repo:   "",
		},
		{
			commit:               "",
			extensions_ignored:   []string{},
			extensions_supported: cfg.DefaultScanFileExtensions,
			fileObjectFunc:       fileObjectFuncFooBlob,
			ignore:               true,
			lines: []string{
				"FOO",
			},
			name:   "IgnoreExtensionsByPolicy",
			path:   "test.png",
			reason: IgnoreReasonFileExtensionIgnoredByPolicy,
			repo:   "",
		},
		{
			commit:               "",
			extensions_ignored:   []string{},
			extensions_supported: cfg.DefaultScanFileExtensions,
			fileObjectFunc:       fileObjectFuncFooBlob,
			ignore:               true,
			lines: []string{
				"FOO",
			},
			name:   "IgnoreExtensionsNotIncludedInConfig",
			path:   "test.random_extension_not_included_by_default",
			reason: IgnoreReasonDefault,
			repo:   "",
		},
		{
			commit:               "6ecf0ef2c2dffb796033e5a02219af86ec6584e5",
			extensions_ignored:   []string{},
			extensions_supported: cfg.DefaultScanFileExtensions,
			fileObjectFunc:       fileObjectFuncFixture,
			ignore:               true,
			lines:                []string{},
			name:                 "FileObjectIsBinary",
			path:                 "binary.jpg",
			reason:               IgnoreReasonFileIsBinary,
			repo:                 "https://github.com/git-fixtures/basic.git",
		},
		{
			commit:               "6ecf0ef2c2dffb796033e5a02219af86ec6584e5",
			extensions_ignored:   []string{},
			extensions_supported: cfg.DefaultScanFileExtensions,
			fileObjectFunc:       fileObjectFuncFixture,
			ignore:               true,
			lines:                []string{},
			name:                 "IgnoreFilePath",
			path:                 "php/crappy.php",
			reason:               IgnoreReasonFilePath,
			repo:                 "https://github.com/git-fixtures/basic.git",
		},
		{
			commit:               "6ecf0ef2c2dffb796033e5a02219af86ec6584e5",
			extensions_ignored:   []string{},
			extensions_supported: cfg.DefaultScanFileExtensions,
			fileObjectFunc:       fileObjectFuncFixture,
			ignore:               true,
			lines: []string{
				"*.class",
				"",
				"# Mobile Tools for Java (J2ME)",
				".mtj.tmp/",
				"",
				"# Package Files #",
				"*.jar",
				"*.war",
				"*.ear",
				"",
				"# virtual machine crash logs, see http://www.java.com/en/download/help/error_hotspot.xml",
				"hs_err_pid*",
			},
			name:   "IgnoreFileName",
			path:   ".gitignore",
			reason: IgnoreReasonFileName,
			repo:   "https://github.com/git-fixtures/basic.git",
		},
		{
			commit:               "6ecf0ef2c2dffb796033e5a02219af86ec6584e5",
			extensions_ignored:   []string{},
			extensions_supported: cfg.DefaultScanFileExtensions,
			fileObjectFunc:       fileObjectFuncFixture,
			ignore:               false,
			lines: []string{
				`{`,
				`	"glossary": {`,
				`		"title": "example glossary\`,
				`		"GlossDiv": {`,
				`			"title": "S",`,
				`			"GlossList": {`,
				`				"GlossEntry": {`,
				`					"Id": "SGML",`,
				`					"SortAs": "SGML",`,
				`					"GlossTerm": "Standard Generalized Markup Language",`,
				`					"Acronym": "SGML",`,
				`					"Abbrev": "ISO 8879:1986",`,
				`					"GlossDef": {`,
				`						"para": "A meta-markup language, used to create markup languages such as DocBook.",`,
				`						"GlossSeeAlso": ["GML", "XML"]`,
				`					},`,
				`					"GlossSee": "markup"`,
				`				}`,
				`			}`,
				`		}`,
				`	}`,
				`}`,
			},
			name:   "ShortJSON",
			path:   "json/short.json",
			reason: "",
			repo:   "https://github.com/git-fixtures/basic.git",
		},
		{
			commit:               "6ecf0ef2c2dffb796033e5a02219af86ec6584e5",
			extensions_ignored:   []string{},
			extensions_supported: cfg.DefaultScanFileExtensions,
			fileObjectFunc:       fileObjectFuncFixture,
			ignore:               true,
			lines: []string{
				"package main",
				"",
				`import "fmt"`,
				"",
				"func main() {",
				`	fmt.Println("Hello, playground")`,
				"}",
			},
			name:   "VendorFooGo",
			path:   "vendor/foo.go",
			reason: IgnoreReasonDirPath,
			repo:   "https://github.com/git-fixtures/basic.git",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test_file_object := test.fileObjectFunc(test.repo, test.commit, test.path)

			// run the function under test
			ignore, reason := IgnoreFileObject(
				test_file_object,
				test.extensions_supported,
				test.extensions_ignored,
			)

			// assert the expected results
			assert.Equalf(t, test.ignore, ignore, "ignore should be %t", test.ignore)
			if test.ignore {
				assert.NotEmpty(t, reason, "reason should not be empty when ignore=true")
			}
			assert.Equal(t, test.reason, reason)
		})
	}
}

// TestIgnoreFilePath() unit test function is used to test the
// IgnoreFilePath() function.
func TestIgnoreFilePath(t *testing.T) {
	tests := []struct {
		ignore bool
		path   string
		reason string
	}{
		{
			ignore: false,
			path:   "/full/path/to/file.txt",
			reason: "",
		},
		{
			ignore: false,
			path:   "relative/path/to/file.txt",
			reason: "",
		},
		{
			ignore: true,
			path:   "LOCK",
			reason: IgnoreReasonFileName,
		},
		{
			ignore: true,
			path:   "vendor/path/to/ignored_file.txt",
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

// TestIgnorePath() unit test function is used to test the ignorePath() function.
func TestIgnorePath(t *testing.T) {
	tests := []struct {
		ignore bool
		name   string
		path   string
		reason string
	}{
		{
			ignore: false,
			name:   "Full_Path",
			path:   "/full/path/to/file.txt",
			reason: "",
		},
		{
			ignore: false,
			name:   "Relative_Path",
			path:   "relative/path/to/file.txt",
			reason: "",
		},
		{
			ignore: true,
			name:   "Ignore_dot_git",
			path:   ".git",
			reason: IgnoreReasonFilePath,
		},
		{
			ignore: true,
			name:   "Ignore_vendor",
			path:   "vendor/path/to/ignored_file.txt",
			reason: IgnoreReasonDirPath,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ignore, reason := ignorePath(test.path)
			assert.Equal(t, test.ignore, ignore, "ignore should be %v", test.ignore)
			assert.Equal(t, test.reason, reason, "reason should be %q", test.reason)
		})
	}
}
