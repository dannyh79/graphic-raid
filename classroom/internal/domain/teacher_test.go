package domain_test

import (
	"bytes"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	d "github.com/dannyh79/graphic-raid/classroom/internal/domain"
)

var _ = Describe("Teacher", func() {
	var buf bytes.Buffer
	t := d.Teacher{&buf}

	BeforeEach(func() {
		buf.Reset()
	})

	It("greets students", func() {
		t.Say(d.Message{Type: d.Greet})

		Expect(buf.String()).To(Equal("Teacher: Guys, are you ready?\n"))
	})

	It("gives a quiz", func() {
		t.Say(d.Message{Type: d.Ask})

		Expect(buf.String()).To(MatchRegexp(`Teacher: \d+ [\+|-|\*|/] \d+ = \?\n`))
	})

	It("responds to answer", func() {
		to := "C"
		t.Say(d.Message{Type: d.Respond, To: to})

		Expect(buf.String()).To(Equal(fmt.Sprintf("Teacher: %s, you are right!\n", to)))
	})
})
