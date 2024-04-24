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
	Answer Sentence = "answer"
)

func (s *Student) getName() string { return fmt.Sprintf("Student %s", s.Name) }
func (s *Student) say(str string)  { fmt.Fprintf(s.Writer, "%s: %s\n", s.getName(), str) }

func (s *Student) Say(m Message) {
	switch m.Type {
	case Answer:
		s.say(fmt.Sprintf("1 + 1 = 2!"))
	case Respond:
		s.say(fmt.Sprintf("%s, you win.", m.To))
	}
}
