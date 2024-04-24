package main

import (
	"os"

	classroom "github.com/dannyh79/graphic-raid/classroom/internal"
)

func main() {
	w := os.Stdout
	s := classroom.Sleeper{}
	classroom.HoldMathQuiz(w, s)
}
