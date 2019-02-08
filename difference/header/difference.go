package header

import (
	"bytes"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

// CompareHeaders comparing two headers and returns true or false
func CompareHeaders(candidate, original http.Header, excluseValuesFromKeys ...string) (bool, string) {

	raw := make(map[string]interface{})

	if candidate == nil && original == nil {
		return true, toString(raw)
	}

	if candidate == nil && original != nil {
		raw = copy(original, "%v =>  ")
		return false, toString(raw)
	}

	if candidate != nil && original == nil {
		raw = copy(candidate, "  => %v")
		return false, toString(raw)
	}

	if len(original) > len(candidate) {
		for key, value := range original {
			if _, ok := candidate[key]; !ok {
				// Candidate does not contain the original header
				raw[key] = fmt.Sprintf("%v =>  ", value)
			}
		}

		return false, toString(raw)
	}

	if len(candidate) > len(original) {
		for key, value := range candidate {
			if _, ok := original[key]; !ok {
				// Candidate does not contain the original header
				raw[key] = fmt.Sprintf("  => %v", value)
			}
		}

		return false, toString(raw)
	}

	for key, value := range original {

		if val, ok := candidate[key]; ok {
			if contains(excluseValuesFromKeys, key) {
				// We do not care about value since it is in exclusions list
				continue
			}

			if !reflect.DeepEqual(value, val) {
				raw[key] = fmt.Sprintf("%v => %v", val, value)
			}

		} else {
			// Candidate does not contain the key
			raw[key] = fmt.Sprintf("%v =>  ", value)
		}
	}

	// In case of Candidate hasheaders not present in original
	for key, value := range candidate {
		if _, ok := original[key]; !ok {
			raw[key] = fmt.Sprintf("  => %v", value)
		}
	}
	return len(raw) == 0, toString(raw)
}

func copy(headers http.Header, expr string) map[string]interface{} {
	raw := make(map[string]interface{})
	for key, value := range headers {
		raw[key] = fmt.Sprintf(expr, value)
	}
	return raw
}

func toString(raw map[string]interface{}) string {
	var b bytes.Buffer

	for key, value := range raw {
		b.WriteString(key)
		b.WriteString(":")
		b.WriteString(fmt.Sprintf("%v", value))
		b.WriteString("\n")
	}

	return strings.Trim(b.String(), " \n")
}

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
