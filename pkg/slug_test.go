package pkg_test

import (
	"github.com/jedi-knights/tds-api/pkg"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Slug", func() {
	Describe("CreateSlug", func() {
		It("should return a slug", func() {
			// Arrange
			input := "Hello World"

			// Act
			output := pkg.CreateSlug(input)

			// Assert
			Expect(output).To(Equal("hello-world"))
		})

		It("should return a slug with a number", func() {
			// Arrange
			input := "Hello World 123"

			// Act
			output := pkg.CreateSlug(input)

			// Assert
			Expect(output).To(Equal("hello-world-123"))
		})

		It("should return a slug with a number and a special character", func() {
			// Arrange
			input := "Hello World 123 !"

			// Act
			output := pkg.CreateSlug(input)

			// Assert
			Expect(output).To(Equal("hello-world-123"))
		})

		It("should return an empty string when given an empty string", func() {
			// Arrange
			input := ""

			// Act
			output := pkg.CreateSlug(input)

			// Assert
			Expect(output).To(Equal(""))
		})
	})
})
