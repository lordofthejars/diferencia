package json_test

import (
	"fmt"

	"github.com/lordofthejars/diferencia/difference/json"
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

			It("should return noise operations with replace", func() {
				documentA := loadFromFile("test_fixtures/document-a.json")
				documentB := loadFromFile("test_fixtures/document-a-change-date.json")

				noiseOperation := json.NoiseOperation{}

				error := noiseOperation.Detect(documentA, documentB)

				Expect(error).Should(Succeed())
				Expect(noiseOperation.Patch).Should(HaveLen(4))

				Expect(noiseOperation.Patch).Should(ContainElement(jsonpatch.NewPatch("replace", "/now/epoch", 0)))
				Expect(noiseOperation.Patch).Should(ContainElement(jsonpatch.NewPatch("replace", "/now/iso8601", 0)))
				Expect(noiseOperation.Patch).Should(ContainElement(jsonpatch.NewPatch("replace", "/now/rfc2822", 0)))
				Expect(noiseOperation.Patch).Should(ContainElement(jsonpatch.NewPatch("replace", "/now/rfc3339", 0)))
			})

			It("should return noise operations with replace in arrays", func() {
				documentA := loadFromFile("test_fixtures/document-a.json")
				documentB := loadFromFile("test_fixtures/document-a-different-array.json")

				noiseOperation := json.NoiseOperation{}

				error := noiseOperation.Detect(documentA, documentB)

				Expect(error).Should(Succeed())
				Expect(noiseOperation.Patch).Should(HaveLen(1))

				Expect(noiseOperation.Patch).Should(ContainElement(jsonpatch.NewPatch("replace", "/urls/1", 0)))
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

		Context("Initial noise", func() {
			It("should append manual noise to automatic noise", func() {
				documentA := loadFromFile("test_fixtures/document-a.json")
				documentB := loadFromFile("test_fixtures/document-a-change-date.json")

				noiseOperation := json.NoiseOperation{}
				noiseOperation.Initialize([]string{"/now/slang_time"})
				error := noiseOperation.Detect(documentA, documentB)

				Expect(error).Should(Succeed())
				Expect(noiseOperation.Patch).Should(HaveLen(5))

				Expect(noiseOperation.Patch).Should(ContainElement(jsonpatch.NewPatch("replace", "/now/epoch", 0)))
				Expect(noiseOperation.Patch).Should(ContainElement(jsonpatch.NewPatch("replace", "/now/iso8601", 0)))
				Expect(noiseOperation.Patch).Should(ContainElement(jsonpatch.NewPatch("replace", "/now/rfc2822", 0)))
				Expect(noiseOperation.Patch).Should(ContainElement(jsonpatch.NewPatch("replace", "/now/rfc3339", 0)))
				Expect(noiseOperation.Patch).Should(ContainElement(jsonpatch.NewPatch("replace", "/now/slang_time", 0)))
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
				noiseOperation.Detect(documentA, documentB)

				primary, candidate, err := noiseOperation.Remove(documentA, documentB)

				if err != nil {
					Fail(fmt.Sprintf("Failing removing noise. Reason: %q", err))
				}

				result := json.CompareDocuments(candidate, primary, "Strict")

				Expect(result).Should(Equal(true))

			})
			It("should return both documents equal with array changes", func() {

				documentA := loadFromFile("test_fixtures/document-a.json")
				documentB := loadFromFile("test_fixtures/document-a-different-array.json")

				noiseOperation := json.NoiseOperation{}
				noiseOperation.Detect(documentA, documentB)

				primary, candidate, err := noiseOperation.Remove(documentA, documentB)

				if err != nil {
					Fail(fmt.Sprintf("Failing removing noise. Reason: %q", err))
				}

				result := json.CompareDocuments(candidate, primary, "Strict")

				Expect(result).Should(Equal(true))

			})
			It("should return both documents equal with json object changes", func() {

				documentA := loadFromFile("test_fixtures/document-c.json")
				documentB := loadFromFile("test_fixtures/document-c-simple.json")

				noiseOperation := json.NoiseOperation{}
				noiseOperation.Detect(documentA, documentB)

				primary, candidate, err := noiseOperation.Remove(documentA, documentB)

				if err != nil {
					Fail(fmt.Sprintf("Failing removing noise. Reason: %q", err))
				}

				result := json.CompareDocuments(candidate, primary, "Strict")

				Expect(result).Should(Equal(true))

			})
			It("should return both documents not equal if not all noise is removed", func() {

				documentA := loadFromFile("test_fixtures/document-a.json")
				documentB := loadFromFile("test_fixtures/document-a-change-date.json")

				noiseOperation := json.NoiseOperation{}
				noiseOperation.Patch = make([]jsonpatch.JsonPatchOperation, 3)
				noiseOperation.Patch[0] = jsonpatch.NewPatch("replace", "/now/epoch", 0)
				noiseOperation.Patch[1] = jsonpatch.NewPatch("replace", "/now/iso8601", 0)
				noiseOperation.Patch[2] = jsonpatch.NewPatch("replace", "/now/rfc2822", 0)

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
