package utils

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateSlug(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		wantPrefix  string
	}{
		{"simple", "Hello World", "hello-world-"},
		{"upper", "GO IS GREAT", "go-is-great-"},
		{"punctuation", "Foo, bar! baz?", "foo-bar-baz-"},
		{"leading trailing space", "  spaced  ", "spaced-"},
		{"unicode strips", "Tiếng Việt", "ti-ng-vi-t-"},
		{"empty title", "", "-"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := GenerateSlug(tc.title)
			assert.True(t, strings.HasPrefix(got, tc.wantPrefix), "slug %q must start with %q", got, tc.wantPrefix)
			assert.Regexp(t, regexp.MustCompile(`-\d{1,5}$`), got, "must end with random numeric suffix")
		})
	}
}

func TestGenerateSlug_UniqueSuffix(t *testing.T) {
	s1 := GenerateSlug("Same Title")
	s2 := GenerateSlug("Same Title")
	assert.NotEqual(t, s1, s2, "two calls should produce different suffixes")
}
