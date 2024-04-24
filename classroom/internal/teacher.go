package classroom

import (
	"fmt"
	"io"
)

type Teacher struct {
	Writer io.Writer
}

type Sentence string

const (
	Greet   Sentence = "greet"
	Ask     Sentence = "ask"
	Respond Sentence = "respond"
)

type Message struct {
	Type Sentence
	To   string
}

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
