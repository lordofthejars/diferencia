package json

import (
	"github.com/lordofthejars/jsondiff"
)

// CompareDocuments comparing two JSON documents and returns true or false according to configured difference
func CompareDocuments(candidate, original []byte, difference string) bool {
	options := jsondiff.DefaultConsoleOptions()

	result, _ := jsondiff.Compare(candidate, original, &options)

	switch difference {
	case "Strict":
		return result == jsondiff.FullMatch
	case "Subset":
		return result == jsondiff.FullMatch || result == jsondiff.SupersetMatch
	}

	return false

}
