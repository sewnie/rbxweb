package stringutils

import (
	"fmt"
)

// FormatSlice returns a slice string with the elements of values as strings.
func FormatSlice[T any](values []T) []string {
	if len(values) == 0 {
		return nil
	}

	s := make([]string, len(values))
	for i, v := range values {
		s[i] = fmt.Sprintf("%v", v)
	}
	return s
}
