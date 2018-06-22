package exporter_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/lordofthejars/diferencia/core"
	"github.com/lordofthejars/diferencia/exporter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("File Exporter", func() {

	var (
		primary   exporter.Interaction
		secondary exporter.Interaction
		candidate exporter.Interaction
	)

	BeforeEach(func() {
		str := `{"page": 1, "fruits": ["apple", "peach"]}`
		primary = exporter.Interaction{
			URL:        "http://localhost:8080",
			Content:    str,
			StatusCode: 200,
		}
		secondary = exporter.Interaction{
			URL:        "http://localhost:8081",
			Content:    str,
			StatusCode: 200,
		}

		candidate = exporter.Interaction{
			URL:        "http://localhost:8082",
			Content:    str,
			StatusCode: 200,
		}
	})

	Describe("Store Interactions", func() {
		Context("With Full interactions", func() {
			It("should write file", func() {
				tmpfile, err := ioutil.TempFile("", "log.json")

				if err != nil {
					Fail(fmt.Sprintf("Unable to create temporal file. Reason: %q", err))
				}
				defer os.Remove(tmpfile.Name())

				interactions := exporter.Interactions{
					Primary:        primary,
					Secondary:      &secondary,
					Candidate:      candidate,
					DifferenceMode: core.Strict.String(),
					Result:         true,
					Processed:      time.Now(),
				}
				err = exporter.ExportToFile(tmpfile.Name(), interactions)
				if err != nil {
					Fail(fmt.Sprintf("Unable to export results. Reason: %q", err))
				}

				byt, err := ioutil.ReadFile(tmpfile.Name())
				expectedInteractions := &exporter.Interactions{}
				json.Unmarshal(byt, expectedInteractions)

				Expect(expectedInteractions.Primary).Should(Equal(primary))
				Expect(*expectedInteractions.Secondary).Should(Equal(secondary))
				Expect(expectedInteractions.Candidate).Should(Equal(candidate))
				Expect(expectedInteractions.DifferenceMode).Should(Equal(interactions.DifferenceMode))
				Expect(expectedInteractions.Result).Should(Equal(interactions.Result))
			})
		})

		Context("With Simple interactions", func() {
			It("should write file", func() {
				tmpfile, err := ioutil.TempFile("", "log.json")

				if err != nil {
					Fail(fmt.Sprintf("Unable to create temporal file. Reason: %q", err))
				}
				defer os.Remove(tmpfile.Name())

				interactions := exporter.Interactions{
					Primary:        primary,
					Candidate:      candidate,
					DifferenceMode: core.Strict.String(),
					Result:         true,
					Processed:      time.Now(),
				}
				err = exporter.ExportToFile(tmpfile.Name(), interactions)
				if err != nil {
					Fail(fmt.Sprintf("Unable to export results. Reason: %q", err))
				}

				byt, err := ioutil.ReadFile(tmpfile.Name())
				expectedInteractions := &exporter.Interactions{}
				json.Unmarshal(byt, expectedInteractions)

				Expect(expectedInteractions.Primary).Should(Equal(primary))
				Expect(expectedInteractions.Candidate).Should(Equal(candidate))
				Expect(expectedInteractions.DifferenceMode).Should(Equal(interactions.DifferenceMode))
				Expect(expectedInteractions.Result).Should(Equal(interactions.Result))
			})
		})
	})
})
