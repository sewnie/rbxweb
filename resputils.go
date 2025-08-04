package rbxweb

import (
	"fmt"
)

func formatSlice[T any](values []T) []string {
	if len(values) == 0 {
		return nil
	}

	s := make([]string, len(values))
	for i, v := range values {
		s[i] = fmt.Sprintf("%v", v)
	}
	return s
}

func getList[T any](v []T, err error) (*T, error) {
	if err != nil {
		return nil, err
	}
	if len(v) == 0 {
		return nil, nil
	}
	return &v[0], nil
}
