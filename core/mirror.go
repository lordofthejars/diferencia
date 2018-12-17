package core

import (
	"bytes"
	"io"
	"net/http"
)

// MirrorResponse to response
func MirrorResponse(primaryCommunication Communicationcontent, w http.ResponseWriter) {
	copyHeader(w.Header(), primaryCommunication.Header)
	w.WriteHeader(primaryCommunication.StatusCode)
	setCookies(w, primaryCommunication.Cookies)
	r := bytes.NewReader(primaryCommunication.Content)
	io.Copy(w, r)
}

func setCookies(w http.ResponseWriter, cookies []*http.Cookie) {
	for _, cookie := range cookies {
		http.SetCookie(w, cookie)
	}
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
