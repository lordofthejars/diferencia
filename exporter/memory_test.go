package exporter_test

import (
	"github.com/lordofthejars/diferencia/exporter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Memory Exporter", func() {

	BeforeEach(func() {
		exporter.Reset()
	})

	Describe("Store Interactions", func() {
		Context("With Counter", func() {
			It("should create and increment the map with endpoint", func() {

				// Given

				// When
				exporter.IncrementError("GET", "/")

				// Then
				entries := exporter.Entries()
				Expect(entries).Should(HaveLen(1))
				Expect(entries[0].Errors).Should(Equal(1))
			})
			It("should increment the map with endpoint", func() {

				// Given

				// When
				exporter.IncrementError("GET", "/a")
				exporter.IncrementError("GET", "/a")

				// Then
				entries := exporter.Entries()
				Expect(entries).Should(HaveLen(1))
				Expect(entries[0].Errors).Should(Equal(2))
			})
		})
	})

})
