package core_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/lordofthejars/diferencia/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type StubHttpClient struct {
	content []string
	status  []int
	index   int
}

func (httpClient *StubHttpClient) MakeRequest(r *http.Request, url string) (*http.Response, error) {
	response := &http.Response{}
	buff := ioutil.NopCloser(strings.NewReader(httpClient.content[httpClient.index]))
	response.Body = buff
	response.StatusCode = httpClient.status[httpClient.index]
	httpClient.index += 1
	return response, nil
}

var _ = Describe("Proxy", func() {

	Describe("Run Diferencia", func() {
		Context("With no noise reduction", func() {
			It("should return true if both documents are equal", func() {

				// Given
				var httpClient = &StubHttpClient{}
				// Record Http Client responses
				recordContent(httpClient, "test_fixtures/document-a.json", "test_fixtures/document-a.json")
				recordStatus(httpClient, 200, 200)
				core.HttpClient = httpClient

				// Prepare Configuration object
				conf := &core.DiferenciaConfiguration{
					Port:                  8080,
					Primary:               "http://now.httpbin.org/",
					Candidate:             "http://now.httpbin.org/",
					StoreResults:          "",
					DifferenceMode:        core.Strict,
					NoiseDetection:        false,
					AllowUnsafeOperations: false,
				}
				core.Config = conf

				// Create stubbed http.Request object
				url, _ := url.Parse("http://localhost:8080")
				request := createRequest(http.MethodGet, url)

				// When

				result, err := core.Diferencia(&request)

				//Then

				Expect(result).Should(Equal(true))
				Expect(err).Should(Succeed())
			})

			It("should return false if both documents are different", func() {

				// Given
				var httpClient = &StubHttpClient{}
				// Record Http Client responses
				recordContent(httpClient, "test_fixtures/document-a.json", "test_fixtures/document-a-change-date.json")
				recordStatus(httpClient, 200, 200)
				core.HttpClient = httpClient

				// Prepare Configuration object
				conf := &core.DiferenciaConfiguration{
					Port:                  8080,
					Primary:               "http://now.httpbin.org/",
					Candidate:             "http://now.httpbin.org/",
					StoreResults:          "",
					DifferenceMode:        core.Strict,
					NoiseDetection:        false,
					AllowUnsafeOperations: false,
				}
				core.Config = conf

				// Create stubbed http.Request object
				url, _ := url.Parse("http://localhost:8080")
				request := createRequest(http.MethodGet, url)

				// When

				result, err := core.Diferencia(&request)

				//Then

				Expect(result).Should(Equal(false))
				Expect(err).Should(Succeed())
			})
		})

		Context("With noise reduction", func() {
			It("should return true if both documents are same but with different values", func() {

				// Given
				var httpClient = &StubHttpClient{}
				// Record Http Client responses
				recordContent(httpClient, "test_fixtures/document-a.json", "test_fixtures/document-a-change-date.json", "test_fixtures/document-a-change-date.json")
				recordStatus(httpClient, 200, 200, 200)
				core.HttpClient = httpClient

				// Prepare Configuration object
				conf := &core.DiferenciaConfiguration{
					Port:                  8080,
					Primary:               "http://now.httpbin.org/",
					Secondary:             "http://now.httpbin.org/",
					Candidate:             "http://now.httpbin.org/",
					StoreResults:          "",
					DifferenceMode:        core.Strict,
					NoiseDetection:        true,
					AllowUnsafeOperations: false,
				}
				core.Config = conf

				// Create stubbed http.Request object
				url, _ := url.Parse("http://localhost:8080")
				request := createRequest(http.MethodGet, url)

				// When

				result, err := core.Diferencia(&request)

				//Then

				Expect(result).Should(Equal(true))
				Expect(err).Should(Succeed())
			})
		})

		Context("With incorrect configuration", func() {
			It("should return error if safe enabled and unsafe operation", func() {

				// Given
				var httpClient = &StubHttpClient{}
				// Record Http Client responses
				recordContent(httpClient, "test_fixtures/document-a.json", "test_fixtures/document-a.json")
				recordStatus(httpClient, 200, 200)
				core.HttpClient = httpClient

				// Prepare Configuration object
				conf := &core.DiferenciaConfiguration{
					Port:                  8080,
					Primary:               "http://now.httpbin.org/",
					Candidate:             "http://now.httpbin.org/",
					StoreResults:          "",
					DifferenceMode:        core.Strict,
					NoiseDetection:        false,
					AllowUnsafeOperations: false,
				}
				core.Config = conf

				// Create stubbed http.Request object
				url, _ := url.Parse("http://localhost:8080")
				request := createRequest(http.MethodPost, url)

				// When

				result, err := core.Diferencia(&request)

				//Then

				Expect(result).Should(Equal(false))
				Expect(err).Should(HaveOccurred())
			})
		})
	})
})

func createRequest(method string, url *url.URL) http.Request {
	request := http.Request{}
	request.URL = url
	request.Method = method

	return request
}

func recordStatus(httpClient *StubHttpClient, statusCode ...int) {
	var status []int
	for _, v := range statusCode {
		status = append(status, v)
	}
	httpClient.status = status
}

func recordContent(httpClient *StubHttpClient, contentFiles ...string) {
	var content []string
	for _, v := range contentFiles {
		content = append(content, loadFromFile(v))
	}
	httpClient.content = content
}

func loadFromFile(filePath string) string {
	payload, err := ioutil.ReadFile(filePath)
	if err != nil {
		Fail(fmt.Sprintf("Unable to load test fixture. Reason: %q", err))
	}
	return string(payload)
}