package classroom

import (
	"fmt"
	"io"
)

type Student struct {
	Writer io.Writer
	Name   string
}

const (
	Answer SentenceType = "answer"
)

const studentPrefix = "Student"

func (s *Student) Say(sentence Sentence) {
	say := func(str string) { fmt.Fprintf(s.Writer, "%s %s: %s\n", studentPrefix, s.Name, str) }

	switch sentence.Type {
	case Answer:
		say(fmt.Sprintf("1 + 1 = 2!"))
	case Respond:
		say(fmt.Sprintf("%s, you win.", sentence.To))
	}
}
