package classroom

import (
	"fmt"
	"io"
	"math/rand"
	"time"
)

func HoldMathQuiz(w io.Writer) {
	fmt.Fprintf(w, "Teacher: Guys, are you ready?\n")

	time.Sleep(3 * time.Second)

	fmt.Fprintf(w, "Teacher: 1 + 1 = ?\n")

	time.Sleep(time.Duration(rand.Intn(3-1)+1) * time.Second)

	fmt.Fprintf(w, "Student C: 1 + 1 = 2!\n")
	fmt.Fprintf(w, "Teacher: C, you are right!\n")
	fmt.Fprintf(w, "Student A: C, you win.\n")
	fmt.Fprintf(w, "Student B: C, you win.\n")
	fmt.Fprintf(w, "Student D: C, you win.\n")
	fmt.Fprintf(w, "Student E: C, you win.\n")
}
