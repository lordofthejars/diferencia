package core_test

import (
	"net/url"

	"github.com/lordofthejars/diferencia/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Url Creator", func() {

	Describe("Get File part from a URL", func() {
		Context("URL without file", func() {
			It("should return index if the host does not contain any slash at the end", func() {
				url, _ := url.Parse("http://www.google.com:433")

				result := core.ExtractFile(*url)

				Expect(result).To(Equal("index.html"))
			})
			It("should return index if there is a slash at the end but no file", func() {
				url, _ := url.Parse("http://www.google.com:433/site/")

				result := core.ExtractFile(*url)

				Expect(result).To(Equal("index.html"))
			})
			It("should return index if there is no slash at the end but subpath", func() {
				url, _ := url.Parse("http://www.google.com:433/site")

				result := core.ExtractFile(*url)

				Expect(result).To(Equal("index.html"))
			})
		})
		Context("URL with file", func() {
			It("should return the file when it is root", func() {
				url, _ := url.Parse("http://www.google.com:433/script.js")

				result := core.ExtractFile(*url)

				Expect(result).To(Equal("script.js"))
			})
			It("should return the file when it is un subpath", func() {
				url, _ := url.Parse("http://www.google.com:433/site/script.js")

				result := core.ExtractFile(*url)

				Expect(result).To(Equal("script.js"))
			})
		})
	})

	Describe("Transform Host Url to Provided one", func() {
		Context("With host without query path", func() {
			It("should replace host", func() {
				url, _ := url.Parse("http://www.google.com:433")
				replacement := "http://localhost:8080"

				result := core.CreateUrl(*url, replacement)

				Expect(result).To(Equal("http://localhost:8080/"))
			})
		})

		Context("With host and final slash", func() {
			It("should replace host", func() {
				url, _ := url.Parse("http://www.google.com:443/")
				replacement := "http://localhost:8080/"

				result := core.CreateUrl(*url, replacement)

				Expect(result).To(Equal("http://localhost:8080/"))
			})
		})

		Context("With host with path", func() {
			It("should replace only host and append path", func() {
				url, _ := url.Parse("http://www.google.com:433/a/b/c")
				replacement := "http://localhost:8080"

				result := core.CreateUrl(*url, replacement)

				Expect(result).To(Equal("http://localhost:8080/a/b/c"))
			})
		})

		Context("With host with path and query path", func() {
			It("should replace only host and append path and query path", func() {
				url, _ := url.Parse("http://www.google.com:433/a/b/c?q=aaa")
				replacement := "http://localhost:8080"

				result := core.CreateUrl(*url, replacement)

				Expect(result).To(Equal("http://localhost:8080/a/b/c?q=aaa"))
			})
		})

	})

})
