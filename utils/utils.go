package utils

import (
	"fmt"
	"net/http"
	"net/url"
	"proxyScanner/dataType"
	"regexp"
	"strings"
)

func ValuesToJSONMap(keyValues url.Values) dataType.JSONMap {
	result := make(dataType.JSONMap)

	for key, values := range keyValues {
		if len(values) == 1 {
			result[key] = values[0]
		} else {
			result[key] = values
		}
	}

	return result
}

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

	r, err := regexp.Compile(result)

	if err != nil {
		fmt.Errorf("Somthing Wrong with config file ScopeToDomainRegex Error: %v", err)
		panic(err)
	}

	return r
}
