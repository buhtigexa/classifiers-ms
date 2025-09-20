package models

import (
	"strings"
)

// Preallocate buffer for crear cache keys
var keyBuilder strings.Builder

// makeCacheKey constructs cache keys efficiently
func makeCacheKey(parts ...string) string {
	keyBuilder.Reset()
	keyBuilder.Grow(64) // Preallocate typical key size

	for i, part := range parts {
		if i > 0 {
			keyBuilder.WriteByte(':')
		}
		keyBuilder.WriteString(part)
	}

	return keyBuilder.String()
}