package path

import (
	"strings"

	"github.com/sohaha/zlsgo/zstring"
)

// HasWildcardPrefix is​used to determine whether an expression is a wildcard at the beginning
func HasWildcardPrefix(pattern string) bool {
	if len(pattern) == 0 {
		return false
	}
	switch pattern[0] {
	case '?', '*', '{':
		return true
	}
	return false
}

// TrimWildcard is used to intercept the pattern before the first wildcard
func TrimWildcard(pattern string) (trimmed string, hasWildcard bool) {
	var chars []byte
Pattern:
	for i := 0; i < len(pattern); i++ {
		switch pattern[i] {
		case '\\':
			if i == len(pattern)-1 {
				break Pattern
			}
			i++
		case '?', '*', '{':
			hasWildcard = true
			break Pattern
		}
		chars = append(chars, pattern[i])
	}
	return string(chars), hasWildcard
}

// Match returns true if name matches the shell file name pattern
func Match(pattern string, s string) bool {
	switch pattern {
	case "**":
		return true
	case "*":
		if strings.Contains(s, "/") {
			return false
		}
		return true
	}
	return zstring.Match(s, pattern)
}
