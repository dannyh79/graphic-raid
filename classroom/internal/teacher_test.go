package classroom_test

import (
	"bytes"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	c "github.com/dannyh79/graphic-raid/classroom/internal"
)

var _ = Describe("Teacher", func() {
	var buf bytes.Buffer
	t := c.Teacher{&buf}

	BeforeEach(func() {
		buf.Reset()
	})

	It("greets students", func() {
		t.Say(c.Message{Type: c.Greet})

		Expect(buf.String()).To(Equal("Teacher: Guys, are you ready?\n"))
	})

	It("gives a quiz", func() {
		t.Say(c.Message{Type: c.Ask})

		Expect(buf.String()).To(MatchRegexp(`Teacher: \d+ [\+|-|\*|/] \d+ = \?\n`))
	})

	It("responds to answer", func() {
		to := "C"
		t.Say(c.Message{Type: c.Respond, To: to})

		Expect(buf.String()).To(Equal(fmt.Sprintf("Teacher: %s, you are right!\n", to)))
	})
})
