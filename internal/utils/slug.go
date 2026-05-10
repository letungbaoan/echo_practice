package utils

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

func GenerateSlug(title string) string {
	lower := strings.ToLower(title)
	reg := regexp.MustCompile("[^a-z0-9]+")
	slug := reg.ReplaceAllString(lower, "-")
	slug = strings.Trim(slug, "-")
	suffix := fmt.Sprintf("-%d", time.Now().UnixNano()%100000)
	return slug + suffix
}
