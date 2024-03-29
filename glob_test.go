package glob

import (
	"errors"
	"testing"

	"github.com/halimath/expect"
	"github.com/halimath/expect/is"
	"github.com/halimath/fsx"
	"github.com/halimath/fsx/memfs"
)

type test struct {
	pattern, f string
	match      bool
	err        error
}

var tests = []test{
	// Test cases not covered by path.Match
	{"main.go", "main.go", true, nil},
	{"main_test.go", "main_test.go", true, nil},
	{"foo/foo_test.go", "foo/foo_test.go", true, nil},
	{"?.go", "m.go", true, nil},
	{"*.go", "main.go", true, nil},
	{"**/*.go", "main.go", true, nil},
	{"*.go", "*.go", true, nil},

	{"//", "", false, ErrBadPattern},
	{"foo//", "", false, ErrBadPattern},
	{"*?.go", "", false, ErrBadPattern},
	{"?*.go", "", false, ErrBadPattern},
	{"**?.go", "", false, ErrBadPattern},
	{"**f", "", false, ErrBadPattern},
	{"[a-", "", false, ErrBadPattern},
	{"[a-\\", "", false, ErrBadPattern},
	{"[\\", "", false, ErrBadPattern},

	{"**/m.go", "foo.go", false, nil},
	{"**/m.go", "foo/a.go", false, nil},
	{"**/m.go", "m.go", true, nil},
	{"**/m.go", "foo/m.go", true, nil},
	{"**/m.go", "bar/m.go", true, nil},
	{"**/m.go", "foo/bar/m.go", true, nil},

	{"ab[cde]", "abc", true, nil},
	{"ab[cde]", "abd", true, nil},
	{"ab[cde]", "abe", true, nil},
	{"ab[+-\\-]", "ab-", true, nil},
	{"ab[\\--a]", "ab-", true, nil},

	{"[a-fA-F]", "a", true, nil},
	{"[a-fA-F]", "f", true, nil},
	{"[a-fA-F]", "A", true, nil},
	{"[a-fA-F]", "F", true, nil},

	// The following test cases are taken from
	// https://github.com/golang/go/blob/master/src/path/match_test.go and are
	// provided here to test compatebility of the match implementation with the
	// test cases from the golang standard lib.
	{"abc", "abc", true, nil},
	{"*", "abc", true, nil},
	{"*c", "abc", true, nil},
	{"a*", "a", true, nil},
	{"a*", "abc", true, nil},
	{"a*", "ab/c", false, nil},
	{"a*/b", "abc/b", true, nil},
	{"a*/b", "a/c/b", false, nil},
	{"a*b*c*d*e*/f", "axbxcxdxe/f", true, nil},
	{"a*b*c*d*e*/f", "axbxcxdxexxx/f", true, nil},
	{"a*b*c*d*e*/f", "axbxcxdxe/xxx/f", false, nil},
	{"a*b*c*d*e*/f", "axbxcxdxexxx/fff", false, nil},
	{"a*b?c*x", "abxbbxdbxebxczzx", true, nil},
	{"a*b?c*x", "abxbbxdbxebxczzy", false, nil},
	{"ab[c]", "abc", true, nil},
	{"ab[b-d]", "abc", true, nil},
	{"ab[e-g]", "abc", false, nil},
	{"ab[^c]", "abc", false, nil},
	{"ab[^b-d]", "abc", false, nil},
	{"ab[^e-g]", "abc", true, nil},
	{"a\\*b", "a*b", true, nil},
	{"a\\*b", "ab", false, nil},
	{"a?b", "a☺b", true, nil},
	{"a[^a]b", "a☺b", true, nil},
	{"a???b", "a☺b", false, nil},
	{"a[^a][^a][^a]b", "a☺b", false, nil},
	{"[a-ζ]*", "α", true, nil},
	{"*[a-ζ]", "A", false, nil},
	{"a?b", "a/b", false, nil},
	{"a*b", "a/b", false, nil},
	{"[\\]a]", "]", true, nil},
	{"[\\-]", "-", true, nil},
	{"[x\\-]", "x", true, nil},
	{"[x\\-]", "-", true, nil},
	{"[x\\-]", "z", false, nil},
	{"[\\-x]", "x", true, nil},
	{"[\\-x]", "-", true, nil},
	{"[\\-x]", "a", false, nil},
	{"[]a]", "]", false, ErrBadPattern},
	{"[-]", "-", false, ErrBadPattern},
	{"[x-]", "x", false, ErrBadPattern},
	{"[x-]", "-", false, ErrBadPattern},
	{"[x-]", "z", false, ErrBadPattern},
	{"[-x]", "x", false, ErrBadPattern},
	{"[-x]", "-", false, ErrBadPattern},
	{"[-x]", "a", false, ErrBadPattern},
	{"\\", "a", false, ErrBadPattern},
	{"[a-b-c]", "a", false, ErrBadPattern},
	{"[", "a", false, ErrBadPattern},
	{"[^", "a", false, ErrBadPattern},
	{"[^bc", "a", false, ErrBadPattern},
	{"a[", "a", false, ErrBadPattern},
	{"a[", "ab", false, ErrBadPattern},
	{"a[", "x", false, ErrBadPattern},
	{"a/b[", "x", false, ErrBadPattern},
	{"*x", "xxx", true, nil},
}

func TestPattern_Match(t *testing.T) {
	for _, tt := range tests {
		pat, err := New(tt.pattern)
		if err != tt.err && !errors.Is(err, tt.err) {
			t.Errorf("New(%#q): wanted error %v but got %v", tt.pattern, tt.err, err)
		}

		if pat != nil {
			match := pat.Match(tt.f)
			if match != tt.match {
				t.Errorf("New(%#q).Match(%#q): wanted match %v but got %v", tt.pattern, tt.f, tt.match, match)
			}
		}
	}
}

func TestPattern_MatchPrefix(t *testing.T) {
	tests = []test{
		{"**/*.go", "foo/", true, nil},
		{"**/*.go", "foo", true, nil},
		{"**/*.go", "foo/bar/", true, nil},
		{"**/*.go", "foo/bar", true, nil},
		{"*/*.go", "foo", true, nil},
	}

	for _, test := range tests {
		pat, err := New(test.pattern)
		expct := expect.WithMessage(t, "%q", test.pattern)
		expct.That(is.NoError(err))
		got := pat.MatchPrefix(test.f)
		expct.That(is.EqualTo(got, test.match))
	}
}

func TestPattern_GlobFS(t *testing.T) {
	fsys := memfs.New()

	fsx.WriteFile(fsys, "go.mod", []byte{}, 0644)
	fsx.WriteFile(fsys, "go.sum", []byte{}, 0644)
	fsx.MkdirAll(fsys, "cmd", 0755)
	fsx.WriteFile(fsys, "cmd/main.go", []byte{}, 0644)
	fsx.WriteFile(fsys, "cmd/main_test.go", []byte{}, 0644)
	fsx.MkdirAll(fsys, "internal/tool", 0755)
	fsx.WriteFile(fsys, "internal/tool/tool.go", []byte{}, 0644)
	fsx.WriteFile(fsys, "internal/tool/tool_test.go", []byte{}, 0644)
	fsx.MkdirAll(fsys, "internal/cli", 0755)
	fsx.WriteFile(fsys, "internal/cli/cli.go", []byte{}, 0644)
	fsx.WriteFile(fsys, "internal/cli/cli_test.go", []byte{}, 0644)

	pat, err := New("**/*_test.go")
	if err != nil {
		t.Fatal(err)
	}

	files, err := pat.GlobFS(fsys, "")
	if err != nil {
		t.Fatal(err)
	}

	expect.That(t, is.DeepEqualTo(files, []string{
		"cmd/main_test.go",
		"internal/cli/cli_test.go",
		"internal/tool/tool_test.go",
	}))
}
