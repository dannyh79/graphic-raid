package domain

import (
	"fmt"
)

type Teacher struct{}

const (
	Greet   Sentence = "greet"
	Ask     Sentence = "ask"
	Respond Sentence = "respond"
)

func (t *Teacher) getName() string     { return "Teacher" }
func (t *Teacher) say(s string) string { return fmt.Sprintf("%s: %s", t.getName(), s) }

func (t *Teacher) Say(m Message) string {
	var s string
	switch m.Type {
	case Greet:
		s = "Guys, are you ready?"
	case Ask:
		s = "1 + 1 = ?"
	case Respond:
		s = fmt.Sprintf("%s, you are right!", m.To)
	}
	return t.say(s)
}

func NewTeacher() *Teacher { return &Teacher{} }
