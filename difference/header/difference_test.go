package header_test

import (
	"net/http"

	"github.com/lordofthejars/diferencia/difference/header"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Header Comparision", func() {

	Describe("Checking if Http Headers are equal", func() {
		Context("Headers without exclusion set", func() {
			It("should return equal if both are nil", func() {
				// Given

				// When
				result, diff := header.CompareHeaders(nil, nil)

				// Then
				Expect(result).Should(BeTrue())
				Expect(len(diff)).Should(Equal(0))
			})

			It("should return no equal if one is nil", func() {
				// Given
				mapA := http.Header{}
				mapA.Add("A", "B")
				// When
				result, diff := header.CompareHeaders(mapA, nil)

				// Then
				Expect(result).Should(BeFalse())
				Expect(len(diff)).Should(BeNumerically(">", 0))
			})

			It("should return no equal if second is nil", func() {
				// Given
				mapA := http.Header{}
				mapA.Add("A", "B")
				// When
				result, diff := header.CompareHeaders(nil, mapA)

				// Then
				Expect(result).Should(BeFalse())
				Expect(len(diff)).Should(BeNumerically(">", 0))
			})

			It("should return false if elements are not the same", func() {
				// Given
				mapA := http.Header{}
				mapA["Accept"] = []string{"text/html"}

				mapB := http.Header{}
				mapB["Accept-Charset"] = []string{"utf-8"}

				// When
				result, diff := header.CompareHeaders(mapA, mapB)

				// Then
				Expect(result).Should(BeFalse())
				Expect(len(diff)).Should(BeNumerically(">", 0))
			})

			It("should return false if values are not the same", func() {
				// Given
				mapA := http.Header{}
				mapA["Accept"] = []string{"text/html"}

				mapB := http.Header{}
				mapB["Accept"] = []string{"text/html", "text/plain"}

				// When
				result, diff := header.CompareHeaders(mapA, mapB)

				// Then
				Expect(result).Should(BeFalse())
				Expect(len(diff)).Should(BeNumerically(">", 0))
			})

			It("should return false if contains different number of headers", func() {
				// Given
				mapA := http.Header{}
				mapA["Accept"] = []string{"text/html"}

				mapB := http.Header{}
				mapB["Accept-Charset"] = []string{"utf-8"}
				mapB["Accept"] = []string{"text/html"}

				// When
				result, diff := header.CompareHeaders(mapA, mapB)

				// Then
				Expect(result).Should(BeFalse())
				Expect(len(diff)).Should(BeNumerically(">", 0))
			})

			It("should return true if values are the same", func() {
				// Given
				mapA := http.Header{}
				mapA["Accept"] = []string{"text/html"}

				mapB := http.Header{}
				mapB["Accept"] = []string{"text/html"}

				// When
				result, diff := header.CompareHeaders(mapA, mapB)

				// Then
				Expect(result).Should(BeTrue())
				Expect(len(diff)).Should(Equal(0))
			})
		})
		Context("Headers with exclusion set", func() {
			It("should return equal if both contains same keys but different values", func() {
				// Given
				mapA := http.Header{}
				mapA["Accept"] = []string{"text/html"}

				mapB := http.Header{}
				mapB["Accept"] = []string{"text/html", "text/plain"}

				// When
				result, diff := header.CompareHeaders(mapA, mapB, "Accept")

				// Then
				Expect(result).Should(BeTrue())
				Expect(len(diff)).Should(Equal(0))
			})
			It("should return false if contains different number of headers although one is excluded from its value", func() {
				// Given
				mapA := http.Header{}
				mapA["Accept"] = []string{"text/html"}

				mapB := http.Header{}
				mapB["Accept-Charset"] = []string{"utf-8"}
				mapB["Accept"] = []string{"text/html"}

				// When
				result, diff := header.CompareHeaders(mapA, mapB, "Accept-Charset")

				// Then
				Expect(result).Should(BeFalse())
				Expect(len(diff)).Should(BeNumerically(">", 0))
			})
		})
	})
})
