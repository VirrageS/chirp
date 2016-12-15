package cache

import (
	"strconv"
	"strings"
)

type Fields []interface{}

// Joins multiple fields to single key
func convertFieldsToKey(fields Fields) string {
	var stringFields []string
	for _, field := range fields {
		switch field.(type) {
		case string:
			stringFields = append(stringFields, field.(string))
		case int:
			stringFields = append(stringFields, strconv.FormatInt(int64(field.(int)), 10))
		case int32:
			stringFields = append(stringFields, strconv.FormatInt(int64(field.(int32)), 10))
		case int64:
			stringFields = append(stringFields, strconv.FormatInt(int64(field.(int64)), 10))
		}
	}
	return strings.Join(stringFields, "_")
}
