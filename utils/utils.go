package utils

import (
	"regexp"
	"strings"
)

func ScopeToDomainRegex(scopeDomain []string) *regexp.Regexp {

	s := strings.Join(scopeDomain, "|")
	result := "(" + s + ")"

	return regexp.MustCompile(result)
}
