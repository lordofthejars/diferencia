package header_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDiferenciaCore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Diferencia Header Suite")
}
