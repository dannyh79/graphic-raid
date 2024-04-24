package classroom

import (
	"io"
	"math/rand"
	"time"

	d "github.com/dannyh79/graphic-raid/classroom/internal/domain"
)

type TimeSleeper interface {
	Sleep(d time.Duration)
}

type Sleeper struct{}

func (s Sleeper) Sleep(d time.Duration) {
	time.Sleep(d)
}

func HoldMathQuiz(w io.Writer, s TimeSleeper) {
	t := d.NewTeacher(w)
	ss := map[string]*d.Student{}
	for _, n := range []string{"A", "B", "C", "D", "E"} {
		ss[n] = d.NewStudent(w, n)
	}

	t.Say(d.Message{Type: d.Greet})

	s.Sleep(3 * time.Second)

	t.Say(d.Message{Type: d.Ask})

	s.Sleep(time.Duration(rand.Intn(3-1)+1) * time.Second)

	ss["C"].Say(d.Message{Type: d.Answer})
	t.Say(d.Message{Type: d.Respond, To: "C"})
	ss["A"].Say(d.Message{Type: d.Respond, To: "C"})
	ss["B"].Say(d.Message{Type: d.Respond, To: "C"})
	ss["D"].Say(d.Message{Type: d.Respond, To: "C"})
	ss["E"].Say(d.Message{Type: d.Respond, To: "C"})
}
