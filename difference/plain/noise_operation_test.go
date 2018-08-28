package plain_test

import (
	"github.com/lordofthejars/diferencia/difference/plain"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Noise Operation", func() {

	Describe("Finding for Noise between calls", func() {
		Context("Valid request", func() {
			It("should return no noise if no changes", func() {

				primary := []byte("aaaaa")
				secondary := []byte("aaaaa")

				noiseOperation := plain.NoiseOperation{}

				noiseOperation.Detect(primary, secondary)

				Expect(noiseOperation.ContainsNoise()).Should(BeFalse())

			})
			It("should return noise if are different", func() {

				primary := []byte("aaaaa")
				secondary := []byte("aaabb")

				noiseOperation := plain.NoiseOperation{}

				noiseOperation.Detect(primary, secondary)

				Expect(noiseOperation.ContainsNoise()).Should(BeTrue())

			})

			It("should return noise if primary is longer than secondary", func() {

				primary := []byte("aaaaa")
				secondary := []byte("aaa")

				noiseOperation := plain.NoiseOperation{}

				noiseOperation.Detect(primary, secondary)

				Expect(noiseOperation.ContainsNoise()).Should(BeTrue())

			})

			It("should return noise if secondary is longer than primary", func() {

				primary := []byte("aaaaa")
				secondary := []byte("aaaaaaa")

				noiseOperation := plain.NoiseOperation{}

				noiseOperation.Detect(primary, secondary)

				Expect(noiseOperation.ContainsNoise()).Should(BeTrue())

			})

		})
	})
	Describe("Removing Noise from Documents", func() {
		Context("A primary and candidate without noise", func() {
			It("should return both documents without any change", func() {
				primary := []byte("aaaaa")
				secondary := []byte("aaaaa")
				candidate := []byte("aaaaa")

				noiseOperation := plain.NoiseOperation{}

				noiseOperation.Detect(primary, secondary)
				newPrimary, newCandidate := noiseOperation.Remove(primary, candidate)

				Expect(newPrimary).Should(Equal(primary))
				Expect(newCandidate).Should(Equal(candidate))
			})
		})
		Context("A primary and candidate with noise", func() {
			It("should return both documents equal", func() {

				primary := []byte("aaaaa")
				secondary := []byte("aaabb")
				candidate := []byte("aaacc")

				noiseOperation := plain.NoiseOperation{}

				noiseOperation.Detect(primary, secondary)
				newPrimary, newCandidate := noiseOperation.Remove(primary, candidate)

				Expect(newPrimary).Should(Equal([]byte("aaa")))
				Expect(newCandidate).Should(Equal([]byte("aaa")))

			})
			It("should return untouched strings if they are totally different", func() {
				primary := []byte("abcd")
				secondary := []byte("efgh")
				candidate := []byte("ijkl")

				noiseOperation := plain.NoiseOperation{}

				noiseOperation.Detect(primary, secondary)
				newPrimary, newCandidate := noiseOperation.Remove(primary, candidate)

				Expect(newPrimary).Should(Equal(primary))
				Expect(newCandidate).Should(Equal(candidate))
			})
		})
	})
})
