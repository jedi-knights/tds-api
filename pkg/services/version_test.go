package services_test

import (
	"github.com/jedi-knights/tds-api/pkg/services"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("VersionService", func() {
	var service *services.VersionService

	BeforeEach(func() {
		service = services.NewVersion()
	})

	AfterEach(func() {
		service = nil
	})

	Describe("GetVersion", func() {
		It("should return the version", func() {
			// Act
			version, err := service.GetVersion()

			// Assert
			Expect(version).To(Equal("1.0.0"))
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
