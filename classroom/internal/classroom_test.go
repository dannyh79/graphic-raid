package classroom_test

import (
	"bytes"
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

		gomock.InOrder(
			mockSleeper.EXPECT().Sleep(3*time.Second).Do(func(_ time.Duration) {
				Expect(buf.String()).To(Equal("Teacher: Guys, are you ready?\n"))
			}),
			mockSleeper.EXPECT().Sleep(gomock.Any()).Do(func(_ time.Duration) {
				Expect(buf.String()).To(ContainSubstring("Teacher: 1 + 1 = ?\n"))
			}),
		)

		classroom.HoldMathQuiz(&buf, mockSleeper)

		Eventually(output).Should(ContainSubstring(
			"Student C: 1 + 1 = 2!\nTeacher: C, you are right!\n"), "teacher reponding after student's correct answer",
		)
		Eventually(output).Should(ContainSubstring(
			"Student A: C, you win.\n"), "student responding to correct answer",
		)
		Eventually(output).Should(ContainSubstring(
			"Student B: C, you win.\n"), "student responding to correct answer",
		)
		Eventually(output).Should(ContainSubstring(
			"Student D: C, you win.\n"), "student responding to correct answer",
		)
		Eventually(output).Should(ContainSubstring(
			"Student E: C, you win.\n"), "student responding to correct answer",
		)
	})
})
