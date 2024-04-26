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
	q := make(chan string)
	ans := make(chan string)
	ansBy := make(chan string, len(Students)-1)

	go newTeacher(Params{
		writer:     w,
		sleeper:    s,
		answer:     ans,
		quiz:       q,
		answeredBy: ansBy,
	})

	for _, n := range Students {
		wg.Add(1)
		go func() {
			defer wg.Done()

			newStudent(Params{
				writer:     w,
				sleeper:    s,
				name:       n,
				answer:     ans,
				quiz:       q,
				answeredBy: ansBy,
			})
		}()
	}

	wg.Wait()
}
