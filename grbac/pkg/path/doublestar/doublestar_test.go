// This file is mostly copied from Go's path/match_test.go

package doublestar

import (
    "path"
    "path/filepath"
    "runtime"
    "strings"
    "testing"
)

type MatchTest struct {
    pattern, testPath []string // a pattern and path to test the pattern on
    shouldMatch       bool     // true if the pattern should match the path
    expectedErr       error    // an expected error
    testOnDisk        bool     // true: test pattern against files in "test" directory
}

// Tests which contain escapes and symlinks will not work on Windows
var onWindows = runtime.GOOS == "windows"

var matchTests = []MatchTest{
    {[]string{"*"}, []string{""}, true, nil, false},
    {[]string{"*"}, []string{"/"}, false, nil, false},
    {[]string{"/*"}, []string{"/"}, true, nil, false},
    {[]string{"/*"}, []string{"/debug/"}, false, nil, false},
    {[]string{"abc"}, []string{"abc"}, true, nil, true},
    {[]string{"*"}, []string{"abc"}, true, nil, true},
    {[]string{"*c"}, []string{"abc"}, true, nil, true},
    {[]string{"a*"}, []string{"a"}, true, nil, true},
    {[]string{"a*"}, []string{"abc"}, true, nil, true},
    {[]string{"a*"}, []string{"ab", "c"}, false, nil, true},
    {[]string{"a*", "b"}, []string{"abc", "b"}, true, nil, true},
    {[]string{"a*", "b"}, []string{"a", "c", "b"}, false, nil, true},
    {[]string{"a*b*c*d*e*", "f"}, []string{"axbxcxdxe", "f"}, true, nil, true},
    {[]string{"a*b*c*d*e*", "f"}, []string{"axbxcxdxexxx", "f"}, true, nil, true},
    {[]string{"a*b*c*d*e*", "f"}, []string{"axbxcxdxe", "xxx", "f"}, false, nil, true},
    {[]string{"a*b*c*d*e*", "f"}, []string{"axbxcxdxexxx", "fff"}, false, nil, true},
    {[]string{"a*b?c*x"}, []string{"abxbbxdbxebxczzx"}, true, nil, true},
    {[]string{"a*b?c*x"}, []string{"abxbbxdbxebxczzy"}, false, nil, true},
    {[]string{"ab[c]"}, []string{"abc"}, true, nil, true},
    {[]string{"ab[b-d]"}, []string{"abc"}, true, nil, true},
    {[]string{"ab[e-g]"}, []string{"abc"}, false, nil, true},
    {[]string{"ab[^c]"}, []string{"abc"}, false, nil, true},
    {[]string{"ab[^b-d]"}, []string{"abc"}, false, nil, true},
    {[]string{"ab[^e-g]"}, []string{"abc"}, true, nil, true},
    {[]string{"a\\*b"}, []string{"ab"}, false, nil, true},
    {[]string{"a?b"}, []string{"a☺b"}, true, nil, true},
    {[]string{"a[^a]b"}, []string{"a☺b"}, true, nil, true},
    {[]string{"a???b"}, []string{"a☺b"}, false, nil, true},
    {[]string{"a[^a][^a][^a]b"}, []string{"a☺b"}, false, nil, true},
    {[]string{"[a-ζ]*"}, []string{"α"}, true, nil, true},
    {[]string{"*[a-ζ]"}, []string{"A"}, false, nil, true},
    {[]string{"a?b"}, []string{"a", "b"}, false, nil, true},
    {[]string{"a*b"}, []string{"a", "b"}, false, nil, true},
    {[]string{"[\\]a]"}, []string{"]"}, true, nil, !onWindows},
    {[]string{"[\\-]"}, []string{"-"}, true, nil, !onWindows},
    {[]string{"[x\\-]"}, []string{"x"}, true, nil, !onWindows},
    {[]string{"[x\\-]"}, []string{"-"}, true, nil, !onWindows},
    {[]string{"[x\\-]"}, []string{"z"}, false, nil, !onWindows},
    {[]string{"[\\-x]"}, []string{"x"}, true, nil, !onWindows},
    {[]string{"[\\-x]"}, []string{"-"}, true, nil, !onWindows},
    {[]string{"[\\-x]"}, []string{"a"}, false, nil, !onWindows},
    {[]string{"[]a]"}, []string{"]"}, false, ErrBadPattern, true},
    {[]string{"[-]"}, []string{"-"}, false, ErrBadPattern, true},
    {[]string{"[x-]"}, []string{"x"}, false, ErrBadPattern, true},
    {[]string{"[x-]"}, []string{"-"}, false, ErrBadPattern, true},
    {[]string{"[x-]"}, []string{"z"}, false, ErrBadPattern, true},
    {[]string{"[-x]"}, []string{"x"}, false, ErrBadPattern, true},
    {[]string{"[-x]"}, []string{"-"}, false, ErrBadPattern, true},
    {[]string{"[-x]"}, []string{"a"}, false, ErrBadPattern, true},
    {[]string{"\\"}, []string{"a"}, false, ErrBadPattern, !onWindows},
    {[]string{"[a-b-c]"}, []string{"a"}, false, ErrBadPattern, true},
    {[]string{"["}, []string{"a"}, false, ErrBadPattern, true},
    {[]string{"[^"}, []string{"a"}, false, ErrBadPattern, true},
    {[]string{"[^bc"}, []string{"a"}, false, ErrBadPattern, true},
    {[]string{"a["}, []string{"a"}, false, nil, false},
    {[]string{"a["}, []string{"ab"}, false, ErrBadPattern, true},
    {[]string{"*x"}, []string{"xxx"}, true, nil, true},
    {[]string{"[abc]"}, []string{"b"}, true, nil, true},
    {[]string{"a", "**"}, []string{"a"}, false, nil, true},
    {[]string{"a", "**"}, []string{"a", "b"}, true, nil, true},
    {[]string{"a", "**"}, []string{"a", "b", "c"}, true, nil, true},
    {[]string{"**", "c"}, []string{"c"}, true, nil, true},
    {[]string{"**", "c"}, []string{"b", "c"}, true, nil, true},
    {[]string{"**", "c"}, []string{"a", "b", "c"}, true, nil, true},
    {[]string{"**", "c"}, []string{"a", "b"}, false, nil, true},
    {[]string{"**", "c"}, []string{"abcd"}, false, nil, true},
    {[]string{"**", "c"}, []string{"a", "abc"}, false, nil, true},
    {[]string{"a", "**", "b"}, []string{"a", "b"}, true, nil, true},
    {[]string{"a", "**", "c"}, []string{"a", "b", "c"}, true, nil, true},
    {[]string{"a", "**", "d"}, []string{"a", "b", "c", "d"}, true, nil, true},
    {[]string{"a", "\\**"}, []string{"a", "b", "c"}, false, nil, !onWindows},
    {[]string{"a", "", "b", "c"}, []string{"a", "b", "c"}, true, nil, true},
    {[]string{"a", "b", "c"}, []string{"a", "b", "", "c"}, true, nil, true},
    {[]string{"ab{c,d}"}, []string{"abc"}, true, nil, true},
    {[]string{"ab{c,d,*}"}, []string{"abcde"}, true, nil, true},
    {[]string{"ab{c,d}["}, []string{"abcd"}, false, ErrBadPattern, true},
    {[]string{"abc", "**"}, []string{"abc", "b"}, true, nil, true},
    {[]string{"**", "abc"}, []string{"abc"}, true, nil, true},
    {[]string{"abc**"}, []string{"abc", "b"}, false, nil, true},
    {[]string{"broken-symlink"}, []string{"broken-symlink"}, true, nil, !onWindows},
    {[]string{"working-symlink", "c", "*"}, []string{"working-symlink", "c", "d"}, true, nil, !onWindows},
    {[]string{"working-sym*", "*"}, []string{"working-symlink", "c"}, true, nil, !onWindows},
    {[]string{"b", "**", "f"}, []string{"b", "symlink-dir", "f"}, true, nil, !onWindows},
}

func TestMatch(t *testing.T) {
    for idx, tt := range matchTests {
        // Since Match() always uses "/" as the separator, we
        // don't need to worry about the tt.testOnDisk flag
        testMatchWith(t, idx, tt)
    }
}

func testMatchWith(t *testing.T, idx int, tt MatchTest) {
    defer func() {
        if r := recover(); r != nil {
            t.Errorf("#%v. Match(%#q, %#q) panicked: %#v", idx, tt.pattern, tt.testPath, r)
        }
    }()

    // Match() always uses "/" as the separator
    pattern := tt.pattern[0]
    testPath := tt.testPath[0]
    if len(tt.pattern) > 1 {
        pattern = path.Join(tt.pattern...)
    }
    if len(tt.testPath) > 1 {
        testPath = path.Join(tt.testPath...)
    }

    ok, err := Match(pattern, testPath)
    if ok != tt.shouldMatch || err != tt.expectedErr {
        t.Errorf("#%v. Match(%#q, %#q) = %v, %v want %v, %v", idx, pattern, testPath, ok, err, tt.shouldMatch, tt.expectedErr)
    }

    if isStandardPattern(pattern) {
        stdOk, stdErr := path.Match(pattern, testPath)
        if ok != stdOk || !compareErrors(err, stdErr) {
            t.Errorf("#%v. Match(%#q, %#q) != path.Match(...). Got %v, %v want %v, %v", idx, pattern, testPath, ok, err, stdOk, stdErr)
        }
    }
}

func TestPathMatch(t *testing.T) {
    for idx, tt := range matchTests {
        // Even though we aren't actually matching paths on disk, we are using
        // PathMatch() which will use the system's separator. As a result, any
        // patterns that might cause problems on-disk need to also be avoided
        // here in this test.
        if tt.testOnDisk {
            testPathMatchWith(t, idx, tt)
        }
    }
}

func testPathMatchWith(t *testing.T, idx int, tt MatchTest) {
    defer func() {
        if r := recover(); r != nil {
            t.Errorf("#%v. Match(%#q, %#q) panicked: %#v", idx, tt.pattern, tt.testPath, r)
        }
    }()

    pattern := filepath.Join(tt.pattern...)
    testPath := filepath.Join(tt.testPath...)
    ok, err := PathMatch(pattern, testPath)
    if ok != tt.shouldMatch || err != tt.expectedErr {
        t.Errorf("#%v. Match(%#q, %#q) = %v, %v want %v, %v", idx, pattern, testPath, ok, err, tt.shouldMatch, tt.expectedErr)
    }

    if isStandardPattern(pattern) {
        stdOk, stdErr := filepath.Match(pattern, testPath)
        if ok != stdOk || !compareErrors(err, stdErr) {
            t.Errorf("#%v. PathMatch(%#q, %#q) != filepath.Match(...). Got %v, %v want %v, %v", idx, pattern, testPath, ok, err, stdOk, stdErr)
        }
    }
}

func isStandardPattern(pattern string) bool {
    return !strings.Contains(pattern, "**") && indexRuneWithEscaping(pattern, '{') == -1
}

func compareErrors(a, b error) bool {
    if a == nil {
        return b == nil
    }
    return b != nil
}

func inSlice(s string, a []string) bool {
    for _, i := range a {
        if i == s {
            return true
        }
    }
    return false
}

func compareSlices(a, b []string) bool {
    if len(a) != len(b) {
        return false
    }

    diff := make(map[string]int, len(a))

    for _, x := range a {
        diff[x]++
    }

    for _, y := range b {
        if _, ok := diff[y]; !ok {
            return false
        }

        diff[y]--
        if diff[y] == 0 {
            delete(diff, y)
        }
    }

    return len(diff) == 0
}
