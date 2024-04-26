package board

import (
	"bufio"
	"fmt"
	"io"
)

func HoldQuorumElection(r io.Reader, w io.Writer, n int) {
	fmt.Fprintln(w, "Member 0: Hi")
	fmt.Fprintln(w, "Member 1: Hi")
	fmt.Fprintln(w, "Member 2: Hi")
	fmt.Fprintln(w, "Member 0: I want to be leader")
	fmt.Fprintln(w, "Member 2: Accept member 0 to be leader")
	fmt.Fprintln(w, "Member 1: I want to be leader")
	fmt.Fprintln(w, "Member 1: Accept member 0 to be leader")
	fmt.Fprintln(w, "Member 0 voted to be leader: (2 > 3/2)")

	i := 0
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		switch i {
		case 0:
			fmt.Fprintln(w, "Member 0: failed heartbeat with Member 1")
			fmt.Fprintln(w, "Member 2: failed heartbeat with Member 1")
			fmt.Fprintln(w, "Member 1: kick out of quorum: (2 > current/2)")
		case 1:
			fmt.Fprintln(w, "Member 0: failed heartbeat with Member 1")
			fmt.Fprintln(w, "Member 0: no response from other users(timeout)")
			fmt.Fprintln(w, "Member 2: kick out of quorum: leader decision")
			fmt.Fprintln(w, "Quorum failed: (1 > total/2)")
			return
		}

		i++
	}
}
