package main

import (
	"os"

	classroom "github.com/dannyh79/graphic-raid/classroom/internal"
)

func main() {
	w := os.Stdout
	classroom.HoldMathQuiz(w)
}
