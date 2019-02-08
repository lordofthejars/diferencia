package json

import (
	"github.com/lordofthejars/jsondiff"
)

// CompareDocuments comparing two JSON documents and returns true or false according to configured difference
func CompareDocuments(candidate, original []byte, difference string) (bool, string) {
	options := defaultJsonOptions()

	result, output := jsondiff.Compare(candidate, original, &options)

	finalResult := false
	finalOutput := ""

	switch difference {
	case "Strict":
		finalResult = result == jsondiff.FullMatch
	case "Subset":
		finalResult = result == jsondiff.FullMatch || result == jsondiff.SupersetMatch
	}

	if !finalResult {
		finalOutput = output
	}

	return finalResult, finalOutput

}

func defaultJsonOptions() jsondiff.Options {
	return jsondiff.Options{
		Added:   jsondiff.Tag{Begin: "", End: ""},
		Removed: jsondiff.Tag{Begin: "", End: ""},
		Changed: jsondiff.Tag{Begin: "", End: ""},
		Indent:  "    ",
	}
}
