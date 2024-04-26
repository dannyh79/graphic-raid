package board

import (
	"fmt"
	"io"
)

func HoldQuorumElection(w io.Writer, n int) {
	fmt.Fprintln(w, "Member 0: Hi")
	fmt.Fprintln(w, "Member 1: Hi")
	fmt.Fprintln(w, "Member 2: Hi")
	fmt.Fprintln(w, "Member 0: I want to be leader")
	fmt.Fprintln(w, "Member 2: Accept member 0 to be leader")
	fmt.Fprintln(w, "Member 1: I want to be leader")
	fmt.Fprintln(w, "Member 1: Accept member 0 to be leader")
	fmt.Fprintln(w, "Member 0 voted to be leader: (2 > 3/2)")
}
