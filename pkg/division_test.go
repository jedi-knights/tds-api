package pkg_test

import (
	"github.com/jedi-knights/tds-api/pkg"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Division", func() {
	Describe("DivisionToString", func() {
		It("should return 'di' for DivisionDI", func() {
			// Arrange
			var division = pkg.Division(pkg.DivisionDI)

			// Act
			result := pkg.DivisionToString(division)

			// Assert
			Expect(result).To(Equal("di"))
		})

		It("should return 'dii' for DivisionDII", func() {
			// Arrange
			var division = pkg.Division(pkg.DivisionDII)

			// Act
			result := pkg.DivisionToString(division)

			// Assert
			Expect(result).To(Equal("dii"))
		})

		It("should return 'diii' for DivisionDIII", func() {
			// Arrange
			var division = pkg.Division(pkg.DivisionDIII)

			// Act
			result := pkg.DivisionToString(division)

			// Assert
			Expect(result).To(Equal("diii"))
		})

		It("should return 'naia' for DivisionNAIA", func() {
			// Arrange
			var division = pkg.Division(pkg.DivisionNAIA)

			// Act
			result := pkg.DivisionToString(division)

			// Assert
			Expect(result).To(Equal("naia"))
		})

		It("should return 'njcaa' for DivisionNJCAA", func() {
			// Arrange
			var division = pkg.Division(pkg.DivisionNJCAA)

			// Act
			result := pkg.DivisionToString(division)

			// Assert
			Expect(result).To(Equal("njcaa"))
		})

		It("should return 'all' for DivisionAll", func() {
			// Arrange
			var division = pkg.Division(pkg.DivisionAll)

			// Act
			result := pkg.DivisionToString(division)

			// Assert
			Expect(result).To(Equal("all"))
		})
	})

	Describe("StringToDivision", func() {
		It("should return DivisionDI for 'di'", func() {
			// Act
			result := pkg.StringToDivision("di")

			// Assert
			Expect(result).To(Equal(pkg.Division(pkg.DivisionDI)))
		})

		It("should return DivisionDII for 'dii'", func() {
			// Act
			result := pkg.StringToDivision("dii")

			// Assert
			Expect(result).To(Equal(pkg.Division(pkg.DivisionDII)))
		})

		It("should return DivisionDIII for 'diii'", func() {
			// Act
			result := pkg.StringToDivision("diii")

			// Assert
			Expect(result).To(Equal(pkg.Division(pkg.DivisionDIII)))
		})

		It("should return DivisionNAIA for 'naia'", func() {
			// Act
			result := pkg.StringToDivision("naia")

			// Assert
			Expect(result).To(Equal(pkg.Division(pkg.DivisionNAIA)))
		})

		It("should return DivisionNJCAA for 'njcaa'", func() {
			// Act
			result := pkg.StringToDivision("njcaa")

			// Assert
			Expect(result).To(Equal(pkg.Division(pkg.DivisionNJCAA)))
		})

		It("should return DivisionUnknown for 'foo'", func() {
			// Act
			var result = pkg.StringToDivision("foo")

			// Assert
			Expect(result).To(Equal(pkg.Division(pkg.DivisionUnknown)))
		})

		It("should return DivisionAll for an empty string", func() {
			// Act
			var result = pkg.StringToDivision("")

			// Assert
			Expect(result).To(Equal(pkg.Division(pkg.DivisionAll)))
		})
	})
})
