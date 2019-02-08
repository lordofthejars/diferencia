package core

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func duplicate(request *http.Request) (dup *http.Request) {
	var bodyBytes []byte
	if request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(request.Body)
	}
	request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	dup = &http.Request{
		Method:        request.Method,
		URL:           request.URL,
		Proto:         request.Proto,
		ProtoMajor:    request.ProtoMajor,
		ProtoMinor:    request.ProtoMinor,
		Body:          ioutil.NopCloser(bytes.NewBuffer(bodyBytes)),
		Host:          request.Host,
		ContentLength: request.ContentLength,
		Close:         true,
		RequestURI:    request.RequestURI,
		TLS:           request.TLS,
	}
	copyTrailer(request, dup)
	copyTransferEncoding(request, dup)
	copyHeaders(request, dup)

	return
}

func copyTransferEncoding(original, dup *http.Request) {
	copy(dup.TransferEncoding, original.TransferEncoding)
}

func copyTrailer(original, dup *http.Request) {
	if original.Trailer != nil {
		newHeaders := http.Header{}
		for k, v := range original.Trailer {
			newHeaders[k] = v
		}
		dup.Trailer = newHeaders
	}
}

func copyHeaders(original, dup *http.Request) {
	if original.Header != nil {
		dup.Header = http.Header{}
		for header, values := range original.Header {
			for _, value := range values {
				dup.Header.Add(header, value)
			}
		}
	}
}
