package board_test

import (
	"bytes"
	"io"

	board "github.com/dannyh79/graphic-raid/quorum/internal"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Quorum Election", func() {
	var (
		buf *bytes.Buffer
		w   *io.PipeWriter
		r   *io.PipeReader
		o   func() string = func() string { return buf.String() }
	)

	BeforeEach(func() {
		buf = new(bytes.Buffer)
		r, w = io.Pipe()
	})

	AfterEach(func() {
		w.Close()
		r.Close()
	})

	It("prints the quorum election process of 3 members", func() {
		go board.HoldQuorumElection(r, buf, 3)

		_, err := w.Write([]byte("kill 1\n"))
		Expect(err).NotTo(HaveOccurred())

		Eventually(o).Should(ContainSubstring(
			"Member 0: failed heartbeat with Member 1\n"+
				"Member 2: failed heartbeat with Member 1\n"+
				"Member 1: kick out of quorum: (2 > current/2)\n",
		), `after command "kill 1"`)

		_, err = w.Write([]byte("kill 2\n"))
		Expect(err).NotTo(HaveOccurred())

		Eventually(o).Should(ContainSubstring(
			"Member 0: no response from other users(timeout)\n"+
				"Member 2: kick out of quorum: leader decision\n"+
				"Quorum failed: (1 > total/2)\n",
		), `followed by command "kill 2"`)

		Eventually(o).Should(ContainSubstring(
			"Member 0: Hi\n" +
				"Member 1: Hi\n" +
				"Member 2: Hi\n" +
				"Member 0: I want to be leader\n" +
				"Member 2: Accept member 0 to be leader\n" +
				"Member 1: I want to be leader\n" +
				"Member 1: Accept member 0 to be leader\n" +
				"Member 0 voted to be leader: (2 > 3/2)\n" +
				"Member 0: failed heartbeat with Member 1\n" +
				"Member 2: failed heartbeat with Member 1\n" +
				"Member 1: kick out of quorum: (2 > current/2)\n" +
				"Member 0: failed heartbeat with Member 1\n" +
				"Member 0: no response from other users(timeout)\n" +
				"Member 2: kick out of quorum: leader decision\n" +
				"Quorum failed: (1 > total/2)\n",
		))
	})
})
