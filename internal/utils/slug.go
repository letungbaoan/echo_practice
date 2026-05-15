package utils

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
)

func GenerateSlug(title string) string {
	lower := strings.ToLower(title)
	reg := regexp.MustCompile("[^a-z0-9]+")
	slug := reg.ReplaceAllString(lower, "-")
	slug = strings.Trim(slug, "-")
	suffix := fmt.Sprintf("-%d", rand.Intn(100000))
	return slug + suffix
}
