package core

import (
	"bytes"
	"io"
	"net/http"
)

// MirrorResponse to response
func MirrorResponse(primaryCommunication communicationcontent, w http.ResponseWriter) {
	copyHeader(w.Header(), primaryCommunication.header)
	w.WriteHeader(primaryCommunication.statusCode)
	r := bytes.NewReader(primaryCommunication.content)
	io.Copy(w, r)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
