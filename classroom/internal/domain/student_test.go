package domain_test

import (
	"bytes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	d "github.com/dannyh79/graphic-raid/classroom/internal/domain"
)

var _ = Describe("Student", func() {
	var buf bytes.Buffer
	s := d.Student{&buf, "Someone"}

	BeforeEach(func() {
		buf.Reset()
	})

	It("answers the quiz", func() {
		s.Say(d.Message{Type: d.Answer})

		Expect(buf.String()).To(MatchRegexp(`Student Someone: \d+ [\+|-|\*|/] \d+ = \d+\!\n`))
	})

	It("responds to answer", func() {
		to := "Another"
		s.Say(d.Message{Type: d.Respond, To: to})

		Expect(buf.String()).To(Equal("Student Someone: Another, you win.\n"))
	})
})
