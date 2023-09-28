package pkg_test

import (
	"github.com/jedi-knights/tds-api/pkg"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Gender", func() {
	Describe("GenderToString", func() {
		It("should return 'female' for GenderFemale", func() {
			// Act
			result := pkg.GenderToString(pkg.GenderFemale)

			// Assert
			Expect(result).To(Equal("female"))
		})

		It("should return 'male' for GenderMale", func() {
			// Act
			result := pkg.GenderToString(pkg.GenderMale)

			// Assert
			Expect(result).To(Equal("male"))
		})

		It("should return 'both' for GenderBoth", func() {
			// Act
			result := pkg.GenderToString(pkg.GenderBoth)

			// Assert
			Expect(result).To(Equal("both"))
		})

		It("should return 'both' for GenderUnknown", func() {
			// Act
			result := pkg.GenderToString(pkg.GenderUnknown)

			// Assert
			Expect(result).To(Equal("both"))
		})
	})

	Describe("StringToGender", func() {
		It("should return GenderBoth for 'both'", func() {
			// Act
			var result = pkg.StringToGender("both")

			// Assert
			Expect(result).To(Equal(pkg.Gender(pkg.GenderBoth)))
		})

		It("should return GenderMale for 'male'", func() {
			// Act
			var result = pkg.StringToGender("male")

			// Assert
			Expect(result).To(Equal(pkg.Gender(pkg.GenderMale)))
		})

		It("should return GenderFemale for 'female'", func() {
			// Act
			var result = pkg.StringToGender("female")

			// Assert
			Expect(result).To(Equal(pkg.Gender(pkg.GenderFemale)))
		})

		It("should return GenderUnknown for 'foo'", func() {
			// Act
			var result = pkg.StringToGender("foo")

			// Assert
			Expect(result).To(Equal(pkg.Gender(pkg.GenderUnknown)))
		})
	})
})
