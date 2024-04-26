package domain_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	d "github.com/dannyh79/graphic-raid/classroom/internal/domain"
)

var _ = Describe("Teacher", func() {
	t := d.Teacher{}

	It("greets students", func() {
		r := t.Say(d.Message{Type: d.Greet})

		Expect(r).To(Equal("Teacher: Guys, are you ready?"))
	})

	It("gives a quiz", func() {
		r := t.Say(d.Message{Type: d.Ask})

		Expect(r).To(MatchRegexp(`Teacher: \d+ \+|-|\*|/ \d+ = \?`))
	})

	It("responds to answer", func() {
		to := "C"
		r := t.Say(d.Message{Type: d.Respond, To: to})

		Expect(r).To(Equal(fmt.Sprintf("Teacher: %s, you are right!", to)))
	})
})
