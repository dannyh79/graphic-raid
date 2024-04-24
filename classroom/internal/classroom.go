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

	t.Say(Message{Type: Greet})

	s.Sleep(3 * time.Second)

	t.Say(Message{Type: Ask})

	s.Sleep(time.Duration(rand.Intn(3-1)+1) * time.Second)

	ss["C"].Say(Message{Type: Answer})
	t.Say(Message{Type: Respond, To: "C"})
	ss["A"].Say(Message{Type: Respond, To: "C"})
	ss["B"].Say(Message{Type: Respond, To: "C"})
	ss["D"].Say(Message{Type: Respond, To: "C"})
	ss["E"].Say(Message{Type: Respond, To: "C"})
}
