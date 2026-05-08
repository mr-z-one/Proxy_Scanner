package utils

import (
	"net/http"
	"proxyScanner/dataType"
	"regexp"
	"strings"
)

func HeaderToJSONMap(headers http.Header) dataType.JSONMap {
	result := make(dataType.JSONMap)

	for key, values := range headers {
		if len(values) == 1 {
			result[key] = values[0]
		} else {
			result[key] = values
		}
	}

	return result
}

func ScopeToDomainRegex(scopeDomain []string) *regexp.Regexp {

	s := strings.Join(scopeDomain, "|")
	result := "(" + s + ")"

	return regexp.MustCompile(result)
}
