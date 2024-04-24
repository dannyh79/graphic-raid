package classroom

import (
	"fmt"
	"io"
)

type Teacher struct {
	Writer io.Writer
}

type SentenceType string

const (
	Greet   SentenceType = "greet"
	Ask     SentenceType = "ask"
	Respond SentenceType = "respond"
)

const teacherPrefix string = "Teacher"

type Sentence struct {
	Type SentenceType
	To   string
}

func (t *Teacher) Say(s Sentence) {
	say := func(s string) { fmt.Fprintf(t.Writer, "%s: %s\n", teacherPrefix, s) }

	switch s.Type {
	case Greet:
		say("Guys, are you ready?")
	case Ask:
		say("1 + 1 = ?")
	case Respond:
		say(fmt.Sprintf("%s, you are right!", s.To))
	}
}
