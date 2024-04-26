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

	go newTeacher(TeacherParams{w: w, s: s, ac: ac, qc: qc, dc: dc})

	for _, n := range Students {
		wg.Add(1)
		go newStudent(StudentParams{w, n, s, qc, ac, dc, &wg})
	}

	wg.Wait()
}
