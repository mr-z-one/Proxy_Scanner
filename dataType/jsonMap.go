package dataType

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
)

type JSONMap map[string]interface{}

func (j *JSONMap) JsonMapFlat() string {
	result := make([]string, 0)
	for s := range *j {
		value, _ := (*j)[s]
		if v, ok := value.(string); ok {
			result = append(result, v)
		}
	}
	return strings.Join(result, ";")
}

func (j JSONMap) GetKeyValue(key string) string {
	if key == "" {
		return key
	}
	if j[key] == nil {
		return ""
	}

	if c, ok := j[key]; ok {
		cc, ok := c.(string)
		if !ok {
			return ""
		}
		return cc
	}

	return ""
}

func (j JSONMap) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan JSONMap: %v", value)
	}

	return json.Unmarshal(bytes, j)
}
