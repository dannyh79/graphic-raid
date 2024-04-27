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
	)

	BeforeEach(func() {
		buf.Reset()
	})

	It(`Writes "Member 0: Hi" to buffer upon starting`, func() {
		go board.NewMember(board.MemberParams{Id: "0", Writer: &buf})

		Eventually(o).Should(ContainSubstring("Member 0: Hi\n"))
	})
})
