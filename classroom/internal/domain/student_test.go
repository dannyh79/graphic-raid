package domain_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	d "github.com/dannyh79/graphic-raid/classroom/internal/domain"
)

var _ = Describe("Student", func() {
	s := d.Student{"Someone"}

	It(`answers to the quiz "1 / 2 = ?"`, func() {
		r := s.Say(d.Message{Type: d.Answer, Body: "1 / 2 = ?"})

		Expect(r).To(MatchRegexp(`Student Someone: \d+ \+|-|\*|/ \d+ = -?\d+\.?\d*!`))
	})

	It(`panics to the quiz "a + 1 = ?"`, func() {

		Expect(func() {
			s.Say(d.Message{Type: d.Answer, Body: "a + 1 = ?"})
		}).To(PanicWith("undefined spec: incorrect quiz format"))
	})

	It("responds to answer", func() {
		to := "Another"
		r := s.Say(d.Message{Type: d.Respond, To: to})

		Expect(r).To(Equal("Student Someone: Another, you win."))
	})
})

var _ = Describe("GetStudentName", func() {
	It(`returns "Someone"`, func() {
		s := "Student Someone: Yes!"
		r := d.GetStudentName(s)

		Expect(r).To(Equal("Someone"))
	})

	It(`returns ""`, func() {
		s := "Student : Yes!"
		r := d.GetStudentName(s)

		Expect(r).To(Equal(""))
	})
})
