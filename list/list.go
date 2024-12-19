package list

import "strings"

func IsSliceContainsPrefix(str string, arr []string) bool {
	if str == "" {
		return false
	}

	for _, v := range arr {
		if strings.HasPrefix(v, str) {
			return true
		}
	}

	return false
}
