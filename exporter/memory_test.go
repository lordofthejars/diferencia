package exporter_test

import (
	"time"

	"github.com/lordofthejars/diferencia/exporter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Memory Exporter", func() {

	BeforeEach(func() {
		exporter.Reset()
	})

	Describe("Store Interactions", func() {
		Context("With Error Counter", func() {
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
		Context("With Success Counter", func() {
			It("should create and increment the map with endpoint", func() {

				// Given
				primaryAverage, _ := time.ParseDuration("10ms")
				candidateAverage, _ := time.ParseDuration("20ms")
				// When
				exporter.IncrementSuccess("GET", "/", primaryAverage, candidateAverage)

				// Then
				entries := exporter.Entries()
				Expect(entries).Should(HaveLen(1))
				Expect(entries[0].Errors).Should(Equal(0))
				Expect(entries[0].Success).Should(Equal(1))
				Expect(entries[0].AveragePrimaryDuration).Should(Equal(float32(10)))
				Expect(entries[0].AverageCandidateDuration).Should(Equal(float32(20)))
			})
			It("should increment the map with endpoint", func() {

				// Given
				primaryAverage1, _ := time.ParseDuration("10ms")
				candidateAverage1, _ := time.ParseDuration("2ms")

				primaryAverage2, _ := time.ParseDuration("30ms")
				candidateAverage2, _ := time.ParseDuration("3ms")
				// When
				exporter.IncrementSuccess("GET", "/a", primaryAverage1, candidateAverage1)
				exporter.IncrementSuccess("GET", "/a", primaryAverage2, candidateAverage2)

				// Then
				entries := exporter.Entries()
				Expect(entries).Should(HaveLen(1))
				Expect(entries[0].Errors).Should(Equal(0))
				Expect(entries[0].Success).Should(Equal(2))
				Expect(entries[0].AveragePrimaryDuration).Should(Equal(float32(20)))
				Expect(entries[0].AverageCandidateDuration).Should(Equal(float32(2.5)))
			})
		})
		Context("With Success and Error Counter", func() {
			It("should create and increment the map with error and success", func() {

				// Given
				primaryAverage, _ := time.ParseDuration("10ms")
				candidateAverage, _ := time.ParseDuration("20ms")
				// When
				exporter.IncrementError("GET", "/")
				exporter.IncrementSuccess("GET", "/", primaryAverage, candidateAverage)

				// Then
				entries := exporter.Entries()
				Expect(entries).Should(HaveLen(1))
				Expect(entries[0].Errors).Should(Equal(1))
				Expect(entries[0].Success).Should(Equal(1))
				Expect(entries[0].AveragePrimaryDuration).Should(Equal(float32(10)))
				Expect(entries[0].AverageCandidateDuration).Should(Equal(float32(20)))
			})
			It("should increment the map with endpoint", func() {

				// Given
				primaryAverage1, _ := time.ParseDuration("10ms")
				candidateAverage1, _ := time.ParseDuration("2ms")

				// When
				exporter.IncrementSuccess("GET", "/a", primaryAverage1, candidateAverage1)
				exporter.IncrementError("GET", "/a")

				// Then
				entries := exporter.Entries()
				Expect(entries).Should(HaveLen(1))
				Expect(entries[0].Errors).Should(Equal(1))
				Expect(entries[0].Success).Should(Equal(1))
				Expect(entries[0].AveragePrimaryDuration).Should(Equal(float32(10)))
				Expect(entries[0].AverageCandidateDuration).Should(Equal(float32(2)))
			})
		})
	})
})
