package utils

import (
	"encoding/json"
	"strings"
)

func Dump(v interface{}) string {
	json, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(json)
}

func BuildSortQuery(sort string) string {
	if sort == "" {
		return ""
	}
	splitted := strings.Split(sort, ",")
	sortQueries := []string{}

	for _, s := range splitted {
		if strings.HasPrefix(s, "-") {
			sortQueries = append(sortQueries, strings.TrimLeft(s, "-")+" DESC")
		} else {
			sortQueries = append(sortQueries, s+" ASC")
		}
	}
	return strings.Join(sortQueries, ",")
}
