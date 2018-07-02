package header

import (
	"net/http"
	"reflect"
)

// CompareHeaders comparing two headers and returns true or false
func CompareHeaders(candidate, original http.Header, excluseValuesFromKeys ...string) bool {

	if candidate == nil && original == nil {
		return true
	}

	if (candidate == nil && original != nil) || (candidate != nil && original == nil) {
		return false
	}

	if len(candidate) != len(original) {
		return false
	}

	for key, value := range candidate {

		if val, ok := original[key]; ok {
			if contains(excluseValuesFromKeys, key) {
				// We do not care about value since it is in exclusions list
				return true
			}

			return reflect.DeepEqual(value, val)

		}

		// Original does not contain the key
		return false

	}

	return reflect.DeepEqual(candidate, original)
}

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
