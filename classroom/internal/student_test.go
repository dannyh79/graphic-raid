package classroom_test

import (
	"bytes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	c "github.com/dannyh79/graphic-raid/classroom/internal"
)

var _ = Describe("Student", func() {
	var buf bytes.Buffer
	s := c.Student{&buf, "Someone"}

	BeforeEach(func() {
		buf.Reset()
	})

	It("answers the quiz", func() {
		s.Say(c.Message{Type: c.Answer})

		Expect(buf.String()).To(MatchRegexp(`Student Someone: \d+ [\+|-|\*|/] \d+ = \d+\!\n`))
	})

	It("responds to answer", func() {
		to := "Another"
		s.Say(c.Message{Type: c.Respond, To: to})

		Expect(buf.String()).To(Equal("Student Someone: Another, you win.\n"))
	})
})
