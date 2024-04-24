package domain

import (
	"fmt"
)

type Student struct {
	Name string
}

const (
	Answer Sentence = "answer"
)

func (s *Student) getName() string       { return fmt.Sprintf("Student %s", s.Name) }
func (s *Student) say(str string) string { return fmt.Sprintf("%s: %s", s.getName(), str) }

func (s *Student) Say(m Message) string {
	var str string
	switch m.Type {
	case Answer:
		str = "1 + 1 = 2!"
	case Respond:
		str = fmt.Sprintf("%s, you win.", m.To)
	}
	return s.say(str)
}

func NewStudent(n string) *Student { return &Student{n} }