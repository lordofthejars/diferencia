package core

import (
	"net/http"
)

// Client interface
type Client interface {
	MakeRequest(r *http.Request, url string) (*http.Response, error)
}

// HTTPClient implementation
type HTTPClient struct {
}

// MakeRequest to given url but maintaining r configuration
func (httpClient *HTTPClient) MakeRequest(r *http.Request, url string) (*http.Response, error) {

	newRequest, err := http.NewRequest(r.Method, url, r.Body)

	if err != nil {
		return nil, err
	}

	newRequest.Header = r.Header

	newRequest.ContentLength = r.ContentLength
	newRequest.TransferEncoding = r.TransferEncoding
	newRequest.Close = r.Close
	newRequest.Trailer = r.Trailer

	for _, c := range r.Cookies() {
		newRequest.AddCookie(c)
	}

	client := &http.Client{}
	return client.Do(newRequest)

}
