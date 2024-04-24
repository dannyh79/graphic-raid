package classroom_test

import (
	"bytes"
	"time"

	classroom "github.com/dannyh79/graphic-raid/classroom/internal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("HoldMathQuiz", func() {
	var buf bytes.Buffer
	output := func() string { return buf.String() }

	BeforeEach(func() {
		buf.Reset()
	})

	It(`prints the quiz interactions in a timely manner`, func() {
		Consistently(output).Should(Equal(""), "before the quiz starts")

		go classroom.HoldMathQuiz(&buf)

		time.Sleep(3 * time.Second)

		Eventually(output).Should(ContainSubstring("Teacher: Guys, are you ready?\n"), "after 3 seconds")

		time.Sleep(3 * time.Second)

		Eventually(output).Should(ContainSubstring("Teacher: 1 + 1 = ?\n"), "after 6 seconds")

		time.Sleep(2 * time.Second)

		Eventually(output).Should(ContainSubstring("Student C: 1 + 1 = 2!\nTeacher: C, you are right!\n"), "after 8 seconds")

		time.Sleep(1 * time.Second)

		Eventually(output).Should(ContainSubstring("Student A: C, you win.\n"), "after 9 seconds")
		Eventually(output).Should(ContainSubstring("Student B: C, you win.\n"), "after 9 seconds")
		Eventually(output).Should(ContainSubstring("Student D: C, you win.\n"), "after 9 seconds")
		Eventually(output).Should(ContainSubstring("Student E: C, you win.\n"), "after 9 seconds")
	})
})
