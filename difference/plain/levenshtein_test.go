package plain_test

import (
	"github.com/lordofthejars/diferencia/difference/plain"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Levenshtein Operation", func() {

	Describe("Calculating Lavenshtein Factor", func() {
		Context("Valid strings", func() {
			It("should calculate when there are differences", func() {

				// Given
				str1 := []byte("Orange")
				str2 := []byte("Apple")

				// When
				result := plain.Levenshtein(str1, str2)

				// Then
				Expect(result).Should(Equal(5))
			})

			It("should calculate when there two equals", func() {

				// Given
				str1 := []byte("Orange")
				str2 := []byte("Orange")

				// When
				result := plain.Levenshtein(str1, str2)

				// Then
				Expect(result).Should(Equal(0))
			})
		})

		Context("Invalid strings", func() {
			It("should calculate when one is empty", func() {

				// Given
				str1 := []byte("Orange")
				str2 := []byte("")

				// When
				result := plain.Levenshtein(str1, str2)

				// Then
				Expect(result).Should(Equal(6))
			})
		})
	})

	Describe("Calculating Percentage Factor", func() {

		Context("Calculate diference", func() {
			It("should calculate diference", func() {

				// Given
				str1 := []byte("recommendation v1 from '99634814-sf4cl':")
				str2 := []byte("recommendation v2 from '7544fgh6-uftj7':")

				// When
				result := plain.CalculateSimilarity(str1, str2)

				// Then
				Expect(result).Should(Equal(0.67))

			})
		})
	})

})
