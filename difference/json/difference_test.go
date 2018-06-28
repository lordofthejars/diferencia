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

				result := json.CompareDocuments(documentA, documentB, core.Strict)
				Expect(result).To(Equal(true))
			})
		})

		Context("With subset mode", func() {
			It("should return that are equal", func() {
				documentA := loadFromFile("test_fixtures/document-a.json")
				documentB := loadFromFile("test_fixtures/document-a.json")

				result := json.CompareDocuments(documentA, documentB, core.Subset)
				Expect(result).To(Equal(true))
			})
		})
	})

	Describe("Compare Two Complete Different Json documents", func() {
		Context("With strict mode", func() {
			It("should return that are different", func() {
				documentA := loadFromFile("test_fixtures/document-a.json")
				documentB := loadFromFile("test_fixtures/document-b.json")

				result := json.CompareDocuments(documentB, documentA, core.Strict)
				Expect(result).To(Equal(false))
			})
		})

		Context("With subset mode", func() {
			It("should return that are different", func() {
				documentA := loadFromFile("test_fixtures/document-a.json")
				documentB := loadFromFile("test_fixtures/document-b.json")

				result := json.CompareDocuments(documentB, documentA, core.Subset)
				Expect(result).To(Equal(false))
			})
		})
	})

	Describe("Compare Two Different Json documents by value", func() {
		Context("With strict mode", func() {
			It("should return that are different", func() {
				documentA := loadFromFile("test_fixtures/document-a.json")
				documentB := loadFromFile("test_fixtures/document-a-change-date.json")

				result := json.CompareDocuments(documentB, documentA, core.Strict)
				Expect(result).To(Equal(false))
			})
		})

		Context("With subset mode", func() {
			It("should return that are different", func() {
				documentA := loadFromFile("test_fixtures/document-a.json")
				documentB := loadFromFile("test_fixtures/document-a-change-date.json")

				result := json.CompareDocuments(documentB, documentA, core.Subset)
				Expect(result).To(Equal(false))
			})
		})
	})

	Describe("Compare Two Different Json documents that are subset", func() {
		Context("With strict mode", func() {
			It("should return that are different", func() {
				documentA := loadFromFile("test_fixtures/document-c.json")
				documentB := loadFromFile("test_fixtures/document-c-update.json")

				result := json.CompareDocuments(documentB, documentA, core.Strict)
				Expect(result).To(Equal(false))
			})
		})

		Context("With subset mode", func() {
			It("should return that are equal", func() {
				documentA := loadFromFile("test_fixtures/document-c.json")
				documentB := loadFromFile("test_fixtures/document-c-update.json")

				result := json.CompareDocuments(documentB, documentA, core.Subset)
				Expect(result).To(Equal(true))
			})

			It("should return that are different when different type", func() {
				documentA := loadFromFile("test_fixtures/document-c.json")
				documentB := loadFromFile("test_fixtures/document-c-update-type.json")

				result := json.CompareDocuments(documentB, documentA, core.Subset)
				Expect(result).To(Equal(false))
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
