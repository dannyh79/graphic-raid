package classroom

import (
	"fmt"
	"io"
	"math/rand"
	"sync"

	"time"

	d "github.com/dannyh79/graphic-raid/classroom/internal/domain"
)

type TeacherParams struct {
	w  io.Writer
	s  TimeSleeper
	qc chan string
	ac chan string
	dc chan string
}

func NewTeacher(tp TeacherParams) {
	e := d.NewTeacher()
	p := newPrinter(tp.w)

	p(e.Say(d.Message{Type: d.Greet}))

	tp.s.Sleep(3 * time.Second)

	q := e.Say(d.Message{Type: d.Ask})
	p(q)
	tp.qc <- q

	s := <-tp.ac
	n := d.GetStudentName(s)

	for range cap(tp.dc) {
		tp.dc <- n
	}

	m := e.Say(d.Message{Type: d.Respond, To: n})
	p(m)
}

type StudentParams struct {
	w  io.Writer
	n  string
	s  TimeSleeper
	qc chan string
	ac chan string
	dc chan string
	wg *sync.WaitGroup
}

func NewStudent(sp StudentParams) {
	defer sp.wg.Done()

	e := d.NewStudent(sp.n)
	p := newPrinter(sp.w)

	select {
	case q := <-sp.qc:
		sp.s.Sleep(time.Duration(rand.Intn(3-1)+1) * time.Second)

		a := e.Say(d.Message{Type: d.Answer, Body: q})
		p(a)
		sp.ac <- a
	case n := <-sp.dc:
		m := e.Say(d.Message{Type: d.Respond, To: n})
		p(m)
	}
}

func newPrinter(w io.Writer) func(s string) {
	return func(s string) {
		fmt.Fprintln(w, s)
	}
}
