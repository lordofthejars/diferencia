package json_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDiferenciaCore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Diferencia JSON Suite")
}
