package header

import (
	"net/http"
	"reflect"
)

// CompareHeaders comparing two headers and returns true or false
func CompareHeaders(candidate, original http.Header) bool {
	return reflect.DeepEqual(candidate, original)
}
