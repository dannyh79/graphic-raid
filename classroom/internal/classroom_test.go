package classroom_test

import (
	"bytes"
	"fmt"
	"regexp"
	"time"

	classroom "github.com/dannyh79/graphic-raid/classroom/internal"
	"github.com/dannyh79/graphic-raid/classroom/internal/mocks"
	"go.uber.org/mock/gomock"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("HoldMathQuiz", func() {
	var (
		mockCtrl    *gomock.Controller
		mockSleeper *mocks.MockTimeSleeper
		buf         bytes.Buffer
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockSleeper = mocks.NewMockTimeSleeper(mockCtrl)
		buf.Reset()
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("prints the quiz interactions to buffer in a timely manner", func() {
		output := func() string { return buf.String() }

		Consistently(output).Should(Equal(""), "before the quiz starts")

		go classroom.HoldMathQuiz(&buf, mockSleeper)

		mockSleeper.EXPECT().Sleep(3 * time.Second).Do(func(_ time.Duration) {
			Eventually(output).Should(ContainSubstring("Teacher: Guys, are you ready?\n"))
		}).Times(1)

		mockSleeper.EXPECT().Sleep(gomock.Any()).Do(func(_ time.Duration) {
			Eventually(output).Should(MatchRegexp(`Teacher: \d{1,3} [\+|-|\*|/] \d{1,3} = \?\n`))
		}).Times(1)

		var s string
		Eventually(func() bool {
			output := output()
			name := getStudentName(output)
			if name != "" {
				s = name
				return true
			}
			return false
		}).Should(BeTrue(), "someone answered to the quiz")

		Expect(output()).To(ContainSubstring(
			fmt.Sprintf("Student %s: 1 + 1 = 2!\nTeacher: %s, you are right!\n", s, s),
		))

		for _, n := range classroom.Students {
			if n != s {
				Expect(output()).To(MatchRegexp(`Student %s: %s, you win\.\n`, n, s))
			}
		}
	})
})

func getStudentName(s string) string {
	re := regexp.MustCompile(`Student (\w+):.+!`)
	m := re.FindStringSubmatch(s)
	if len(m) < 1 {
		return ""
	}
	return m[1]
}
