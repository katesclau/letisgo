package helpers

import (
	"context"
	"errors"
)

var ErrKeyNotFound = errors.New("Key not found in context")
var ErrInvalidType = errors.New("Invalid type")

// Get Typed Value from context
// Ref: https://www.willem.dev/articles/how-to-add-values-to-context/
func Get[T any](ctx context.Context, key string) (T, error) {
	var ret T
	v := ctx.Value(key)
	if v == nil {
		return ret, ErrKeyNotFound
	}

	ret, ok := v.(T)
	if !ok {
		return ret, ErrInvalidType
	}
	return ret, nil
}
