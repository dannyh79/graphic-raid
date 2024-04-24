package classroom

import (
	"fmt"
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
	t := d.NewTeacher()
	ss := map[string]*d.Student{}
	for _, n := range []string{"A", "B", "C", "D", "E"} {
		ss[n] = d.NewStudent(n)
	}

	fmt.Fprintln(w, t.Say(d.Message{Type: d.Greet}))

	s.Sleep(3 * time.Second)

	fmt.Fprintln(w, t.Say(d.Message{Type: d.Ask}))

	s.Sleep(time.Duration(rand.Intn(3-1)+1) * time.Second)

	fmt.Fprintln(w, ss["C"].Say(d.Message{Type: d.Answer}))

	fmt.Fprintln(w, t.Say(d.Message{Type: d.Respond, To: "C"}))
	fmt.Fprintln(w, ss["A"].Say(d.Message{Type: d.Respond, To: "C"}))
	fmt.Fprintln(w, ss["B"].Say(d.Message{Type: d.Respond, To: "C"}))
	fmt.Fprintln(w, ss["D"].Say(d.Message{Type: d.Respond, To: "C"}))
	fmt.Fprintln(w, ss["E"].Say(d.Message{Type: d.Respond, To: "C"}))
}
