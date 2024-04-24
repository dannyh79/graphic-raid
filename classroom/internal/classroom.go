package classroom

import (
	"fmt"
	"io"
	"math/rand"
	"time"
)

type TimeSleeper interface {
	Sleep(d time.Duration)
}

type Sleeper struct{}

func (s Sleeper) Sleep(d time.Duration) {
	time.Sleep(d)
}

func HoldMathQuiz(w io.Writer, s TimeSleeper) {
	t := Teacher{w}

	t.Say(Sentence{Type: Greet})

	s.Sleep(3 * time.Second)

	t.Say(Sentence{Type: Ask})

	s.Sleep(time.Duration(rand.Intn(3-1)+1) * time.Second)

	fmt.Fprintf(w, "Student C: 1 + 1 = 2!\n")
	t.Say(Sentence{Type: Respond, To: "C"})
	fmt.Fprintf(w, "Student A: C, you win.\n")
	fmt.Fprintf(w, "Student B: C, you win.\n")
	fmt.Fprintf(w, "Student D: C, you win.\n")
	fmt.Fprintf(w, "Student E: C, you win.\n")
}
