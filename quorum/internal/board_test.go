package board_test

import (
	"bytes"

	board "github.com/dannyh79/graphic-raid/quorum/internal"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("HoldQuorumElection", func() {
	var (
		buf bytes.Buffer
	)

	BeforeEach(func() {
		buf.Reset()
	})

	It("prints the quorum election process of 3 members", func() {
		output := func() string { return buf.String() }

		Consistently(output).Should(Equal(""), "before the program executes")

		go board.HoldQuorumElection(&buf, 3)

		Eventually(output).Should(ContainSubstring(
			`Member 0: Hi
Member 1: Hi
Member 2: Hi
Member 0: I want to be leader
Member 2: Accept member 0 to be leader
Member 1: I want to be leader
Member 1: Accept member 0 to be leader
Member 0 voted to be leader: (2 > 3/2)
`))
	})
})
