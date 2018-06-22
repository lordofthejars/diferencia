package json_test

import (
	"fmt"

	"github.com/lordofthejars/diferencia/json"
	"github.com/mattbaird/jsonpatch"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Noise Operation", func() {

	Describe("Finding for Noise between calls", func() {
		Context("Valid request", func() {
			It("should return no operations if no changes", func() {
				documentA := loadFromFile("test_fixtures/document-a.json")
				documentB := loadFromFile("test_fixtures/document-a.json")

				noiseOperation := json.NoiseOperation{}

				error := noiseOperation.Detect(documentA, documentB)
				Expect(error).Should(Succeed())
				Expect(noiseOperation.Patch).Should(HaveLen(0))
			})

			It("should return noise operations with remove instead of replace", func() {
				documentA := loadFromFile("test_fixtures/document-a.json")
				documentB := loadFromFile("test_fixtures/document-a-change-date.json")

				noiseOperation := json.NoiseOperation{}

				error := noiseOperation.Detect(documentA, documentB)

				Expect(error).Should(Succeed())
				Expect(noiseOperation.Patch).Should(HaveLen(4))

				Expect(noiseOperation.Patch).Should(ContainElement(jsonpatch.NewPatch("remove", "/now/epoch", nil)))
				Expect(noiseOperation.Patch).Should(ContainElement(jsonpatch.NewPatch("remove", "/now/iso8601", nil)))
				Expect(noiseOperation.Patch).Should(ContainElement(jsonpatch.NewPatch("remove", "/now/rfc2822", nil)))
				Expect(noiseOperation.Patch).Should(ContainElement(jsonpatch.NewPatch("remove", "/now/rfc3339", nil)))
			})
		})

		Context("Invalid request", func() {
			It("should return error if new attribute added", func() {
				documentA := loadFromFile("test_fixtures/document-a.json")
				documentB := loadFromFile("test_fixtures/document-a-different-key.json")

				noiseOperation := json.NoiseOperation{}

				error := noiseOperation.Detect(documentA, documentB)
				Expect(error).Should(HaveOccurred())
			})
		})
	})

	Describe("Removing Noise from Documents", func() {
		Context("A primary and candidate without noise", func() {
			It("should return both documents without any change", func() {

				documentA := loadFromFile("test_fixtures/document-a.json")
				documentB := loadFromFile("test_fixtures/document-a.json")

				noiseOperation := json.NoiseOperation{}
				noiseOperation.Patch = make([]jsonpatch.JsonPatchOperation, 0, 0)

				primary, candidate, err := noiseOperation.Remove(documentA, documentB)

				if err != nil {
					Fail(fmt.Sprintf("Failing removing noise. Reason: %q", err))
				}

				result := json.CompareDocuments(candidate, primary, "Strict")

				Expect(result).Should(Equal(true))

			})
		})

		Context("A primary and candidate with noise", func() {
			It("should return both documents equal", func() {

				documentA := loadFromFile("test_fixtures/document-a.json")
				documentB := loadFromFile("test_fixtures/document-a-change-date.json")

				noiseOperation := json.NoiseOperation{}
				noiseOperation.Patch = make([]jsonpatch.JsonPatchOperation, 4)
				noiseOperation.Patch[0] = jsonpatch.NewPatch("remove", "/now/epoch", nil)
				noiseOperation.Patch[1] = jsonpatch.NewPatch("remove", "/now/iso8601", nil)
				noiseOperation.Patch[2] = jsonpatch.NewPatch("remove", "/now/rfc2822", nil)
				noiseOperation.Patch[3] = jsonpatch.NewPatch("remove", "/now/rfc3339", nil)

				primary, candidate, err := noiseOperation.Remove(documentA, documentB)

				if err != nil {
					Fail(fmt.Sprintf("Failing removing noise. Reason: %q", err))
				}

				result := json.CompareDocuments(candidate, primary, "Strict")

				Expect(result).Should(Equal(true))

			})
		})

		Context("A primary and candidate with noise", func() {
			It("should return both documents not equal if not all noise is removed", func() {

				documentA := loadFromFile("test_fixtures/document-a.json")
				documentB := loadFromFile("test_fixtures/document-a-change-date.json")

				noiseOperation := json.NoiseOperation{}
				noiseOperation.Patch = make([]jsonpatch.JsonPatchOperation, 3)
				noiseOperation.Patch[0] = jsonpatch.NewPatch("remove", "/now/epoch", nil)
				noiseOperation.Patch[1] = jsonpatch.NewPatch("remove", "/now/iso8601", nil)
				noiseOperation.Patch[2] = jsonpatch.NewPatch("remove", "/now/rfc2822", nil)

				primary, candidate, err := noiseOperation.Remove(documentA, documentB)

				if err != nil {
					Fail(fmt.Sprintf("Failing removing noise. Reason: %q", err))
				}

				result := json.CompareDocuments(candidate, primary, "Strict")

				Expect(result).Should(Equal(false))

			})
		})
	})
})
