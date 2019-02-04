package json_test

import (
	"fmt"
	"io/ioutil"

	"github.com/lordofthejars/diferencia/core"

	"github.com/lordofthejars/diferencia/difference/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Json Difference", func() {

	Describe("Compare Two Equal Json documents", func() {
		Context("With strict mode", func() {
			It("should return that are equal", func() {
				documentA := loadFromFile("test_fixtures/document-a.json")
				documentB := loadFromFile("test_fixtures/document-a.json")

				result, output := json.CompareDocuments(documentA, documentB, core.Strict.String())
				Expect(result).To(Equal(true))
				Expect(len(output)).To(Equal(0))
			})
		})

		Context("With subset mode", func() {
			It("should return that are equal", func() {
				documentA := loadFromFile("test_fixtures/document-a.json")
				documentB := loadFromFile("test_fixtures/document-a.json")

				result, output := json.CompareDocuments(documentA, documentB, core.Subset.String())
				Expect(result).To(Equal(true))
				Expect(len(output)).To(Equal(0))
			})
		})
	})

	Describe("Compare Two Complete Different Json documents", func() {
		Context("With strict mode", func() {
			It("should return that are different", func() {
				documentA := loadFromFile("test_fixtures/document-a.json")
				documentB := loadFromFile("test_fixtures/document-b.json")

				result, output := json.CompareDocuments(documentB, documentA, core.Strict.String())
				Expect(result).To(Equal(false))
				Expect(len(output)).Should(BeNumerically(">", 0))
			})
		})

		Context("With subset mode", func() {
			It("should return that are different", func() {
				documentA := loadFromFile("test_fixtures/document-a.json")
				documentB := loadFromFile("test_fixtures/document-b.json")

				result, output := json.CompareDocuments(documentB, documentA, core.Subset.String())
				Expect(result).To(Equal(false))
				Expect(len(output)).Should(BeNumerically(">", 0))
			})
		})
	})

	Describe("Compare Two Different Json documents by value", func() {
		Context("With strict mode", func() {
			It("should return that are different", func() {
				documentA := loadFromFile("test_fixtures/document-a.json")
				documentB := loadFromFile("test_fixtures/document-a-change-date.json")

				result, output := json.CompareDocuments(documentB, documentA, core.Strict.String())
				Expect(result).To(Equal(false))
				Expect(len(output)).Should(BeNumerically(">", 0))
			})
		})

		Context("With subset mode", func() {
			It("should return that are different", func() {
				documentA := loadFromFile("test_fixtures/document-a.json")
				documentB := loadFromFile("test_fixtures/document-a-change-date.json")

				result, output := json.CompareDocuments(documentB, documentA, core.Subset.String())
				Expect(result).To(Equal(false))
				Expect(len(output)).Should(BeNumerically(">", 0))
			})
		})
	})

	Describe("Compare Two Different Json documents that are subset", func() {
		Context("With strict mode", func() {
			It("should return that are different", func() {
				documentA := loadFromFile("test_fixtures/document-c.json")
				documentB := loadFromFile("test_fixtures/document-c-update.json")

				result, output := json.CompareDocuments(documentB, documentA, core.Strict.String())
				Expect(result).To(Equal(false))
				Expect(len(output)).Should(BeNumerically(">", 0))
			})
		})

		Context("With subset mode", func() {
			It("should return that are equal", func() {
				documentA := loadFromFile("test_fixtures/document-c.json")
				documentB := loadFromFile("test_fixtures/document-c-update.json")

				result, output := json.CompareDocuments(documentB, documentA, core.Subset.String())
				Expect(result).To(Equal(true))
				Expect(len(output)).To(Equal(0))
			})

			It("should return that are different when different type", func() {
				documentA := loadFromFile("test_fixtures/document-c.json")
				documentB := loadFromFile("test_fixtures/document-c-update-type.json")

				result, output := json.CompareDocuments(documentB, documentA, core.Subset.String())
				Expect(result).To(Equal(false))
				Expect(len(output)).Should(BeNumerically(">", 0))
			})
		})
	})

})

func loadFromFile(filePath string) []byte {
	payload, err := ioutil.ReadFile(filePath)
	if err != nil {
		Fail(fmt.Sprintf("Unable to load test fixture. Reason: %q", err))
	}
	return payload
}
