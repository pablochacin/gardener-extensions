package generator

import (
	"github.com/gardener/gardener-extensions/controllers/os-common/pkg/generator/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var baseUnitsPath = "/etc/systemd/system/"

var _ = Describe("My Generator Test", func() {
	generator, err := NewCloudInitGenerator()
	It("should not error", func() {
		Expect(err).NotTo(HaveOccurred())
	})
	Describe("Conformance Tests", test.DescribeTest(generator, box))
})
