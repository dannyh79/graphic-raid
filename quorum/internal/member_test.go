package board_test

import (
	"bytes"

	board "github.com/dannyh79/graphic-raid/quorum/internal"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewMember", func() {
	var (
		buf bytes.Buffer
		o   func() string = func() string { return buf.String() }
		p   board.MemberParams
	)

	BeforeEach(func() {
		buf.Reset()

		p = board.MemberParams{
			Id:         "0",
			Writer:     &buf,
			WantToLead: func() bool { return true },
		}
	})

	It(`Writes "Member 0: Hi" to buffer upon starting`, func() {
		go board.NewMember(p)

		Eventually(o).Should(ContainSubstring("Member 0: Hi\n"))
	})

	It(`Writes "Member 0: I want to be leader" to buffer`, func() {
		go board.NewMember(p)

		Eventually(o).Should(ContainSubstring("Member 0: I want to be leader\n"))
	})

	It(`does NOT write "Member 0: I want to be leader" to buffer`, func() {
		p.WantToLead = func() bool { return false }

		go board.NewMember(p)

		Eventually(o).Should(Not(ContainSubstring("Member 0: I want to be leader\n")))
	})

	It(`Writes to buffer in a sequential manner`, func() {
		go board.NewMember(p)

		Eventually(o).Should(ContainSubstring(
			"Member 0: Hi\n" +
				"Member 0: I want to be leader\n",
		))
	})
})
