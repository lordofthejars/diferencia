package plain

import (
	"bytes"
)

// CompareDocuments comparing two plain strings documents and returns true or false according to configured difference
func CompareDocuments(candidate, original []byte, difference string) bool {
	return bytes.Equal(candidate, original)
}
