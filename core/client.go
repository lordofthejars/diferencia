package core

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
)

// Client interface
type Client interface {
	MakeRequest(r *http.Request, url string) (*http.Response, error)
}

// HTTPClient implementation
type HTTPClient struct {
	config *DiferenciaConfiguration
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

	// To avoid any nil problem if caller does not set the configuration object (tests)
	if httpClient.config != nil {
		if httpClient.config.InsecureSkipVerify || httpClient.config.AreHttpsClientParamsSet() {
			config := tls.Config{}

			if httpClient.config.InsecureSkipVerify {
				config.InsecureSkipVerify = true
			}

			if httpClient.config.AreHttpsClientParamsSet() {
				caCert, err := ioutil.ReadFile(httpClient.config.CaCert)
				if err != nil {
					return nil, err
				}
				caCertPool := x509.NewCertPool()
				caCertPool.AppendCertsFromPEM(caCert)

				cert, err := tls.LoadX509KeyPair(httpClient.config.ClientCert, httpClient.config.ClientKey)
				if err != nil {
					return nil, err
				}

				config.RootCAs = caCertPool
				config.Certificates = []tls.Certificate{cert}

			}

			tr := &http.Transport{
				TLSClientConfig: &config,
			}

			client.Transport = tr
		}

	}

	return client.Do(newRequest)

}
