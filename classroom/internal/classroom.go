package classroom

import (
	"io"
	"sync"
	"time"
)

type TimeSleeper interface {
	Sleep(d time.Duration)
}

type Sleeper struct{}

func (s Sleeper) Sleep(d time.Duration) {
	time.Sleep(d)
}

var Students = []string{"A", "B", "C", "D", "E"}

func HoldMathQuiz(w io.Writer, s TimeSleeper) {
	var wg sync.WaitGroup
	ac := make(chan string)
	qc := make(chan string)
	dc := make(chan string, len(Students)-1)

	go newTeacher(Params{w: w, s: s, a: ac, q: qc, dc: dc})

	for _, n := range Students {
		wg.Add(1)
		go func() {
			defer wg.Done()

			newStudent(Params{w: w, s: s, n: n, a: ac, q: qc, dc: dc})
		}()
	}

	wg.Wait()
}
