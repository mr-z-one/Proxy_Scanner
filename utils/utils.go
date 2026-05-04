package utils

import (
	"regexp"
	"strings"
)

func ScopeToDomain(scopeDomain []string) *regexp.Regexp {

	s := strings.Join(scopeDomain, "|")
	result := "(" + s + ")"

	return regexp.MustCompile(result)
}
