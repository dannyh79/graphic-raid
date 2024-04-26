package classroom

import (
	"fmt"
	"io"
	"math/rand"

	"time"

	d "github.com/dannyh79/graphic-raid/classroom/internal/domain"
)

type Params struct {
	writer     io.Writer
	sleeper    TimeSleeper
	name       string
	quiz       chan string
	answer     chan string
	answeredBy chan string
}

func newTeacher(tp Params) {
	e := d.NewTeacher()
	p := newPrinter(tp.writer)

	p(e.Say(d.Message{Type: d.Greet}))

	tp.sleeper.Sleep(3 * time.Second)

	q := e.Say(d.Message{Type: d.Ask})
	p(q)
	tp.quiz <- q

	s := <-tp.answer
	n := d.GetStudentName(s)

	for range cap(tp.answeredBy) {
		tp.answeredBy <- n
	}

	m := e.Say(d.Message{Type: d.Respond, To: n})
	p(m)
}

func newStudent(sp Params) {
	e := d.NewStudent(sp.name)
	p := newPrinter(sp.writer)

	select {
	case q := <-sp.quiz:
		sp.sleeper.Sleep(time.Duration(rand.Intn(3-1)+1) * time.Second)

		a := e.Say(d.Message{Type: d.Answer, Body: q})
		p(a)
		sp.answer <- a
	case n := <-sp.answeredBy:
		m := e.Say(d.Message{Type: d.Respond, To: n})
		p(m)
	}
}

func newPrinter(w io.Writer) func(s string) {
	return func(s string) {
		fmt.Fprintln(w, s)
	}
}
