package domain

import (
	"fmt"
	"io"
)

type Teacher struct {
	Writer io.Writer
}

const (
	Greet   Sentence = "greet"
	Ask     Sentence = "ask"
	Respond Sentence = "respond"
)

func (t *Teacher) getName() string { return "Teacher" }
func (t *Teacher) say(s string)    { fmt.Fprintf(t.Writer, "%s: %s\n", t.getName(), s) }

func (t *Teacher) Say(m Message) {
	switch m.Type {
	case Greet:
		t.say("Guys, are you ready?")
	case Ask:
		t.say("1 + 1 = ?")
	case Respond:
		t.say(fmt.Sprintf("%s, you are right!", m.To))
	}
}

func NewTeacher(w io.Writer) *Teacher { return &Teacher{w} }
