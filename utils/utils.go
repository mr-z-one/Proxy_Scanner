package utils

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"proxyScanner/dataType"
	"regexp"
	"strings"
)

func ReadeFile(path string) string {
	f, err := os.Open(path)

	if err != nil {
		fmt.Println("[-]", err)
		return ""
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		fmt.Println("[-]", err)
		return ""
	}
	return string(data)
}

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
