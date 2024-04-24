package classroom

import (
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
	ss := map[string]*Student{}
	for _, n := range []string{"A", "B", "C", "D", "E"} {
		ss[n] = &Student{w, n}
	}

	t.Say(Sentence{Type: Greet})

	s.Sleep(3 * time.Second)

	t.Say(Sentence{Type: Ask})

	s.Sleep(time.Duration(rand.Intn(3-1)+1) * time.Second)

	ss["C"].Say(Sentence{Type: Answer})
	t.Say(Sentence{Type: Respond, To: "C"})
	ss["A"].Say(Sentence{Type: Respond, To: "C"})
	ss["B"].Say(Sentence{Type: Respond, To: "C"})
	ss["D"].Say(Sentence{Type: Respond, To: "C"})
	ss["E"].Say(Sentence{Type: Respond, To: "C"})
}
