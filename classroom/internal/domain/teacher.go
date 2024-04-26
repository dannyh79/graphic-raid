package domain

import (
	"fmt"
	"math/rand"
)

var operators = []string{"+", "-", "*", "/"}

type quiz struct {
	term1    int
	term2    int
	operator string
}

type Teacher struct {
	quiz *quiz
}

const (
	Greet   Sentence = "greet"
	Ask     Sentence = "ask"
	Respond Sentence = "respond"
)

func (t *Teacher) getName() string     { return "Teacher" }
func (t *Teacher) say(s string) string { return fmt.Sprintf("%s: %s", t.getName(), s) }

func (t *Teacher) askQuiz() string {
	t.quiz = newQuiz()
	return fmt.Sprintf("%d %s %d = ?", t.quiz.term1, t.quiz.operator, t.quiz.term2)
}

func (t *Teacher) Say(m Message) string {
	var s string
	switch m.Type {
	case Greet:
		s = "Guys, are you ready?"
	case Ask:
		s = t.askQuiz()
	case Respond:
		s = fmt.Sprintf("%s, you are right!", m.To)
	}
	return t.say(s)
}

func NewTeacher() *Teacher { return &Teacher{} }

func newQuiz() *quiz {
	return &quiz{
		term1:    randomIntegerIn(1, 100),
		term2:    randomIntegerIn(1, 100),
		operator: randomOperator(operators),
	}
}

func randomIntegerIn(l, h int) int     { return rand.Intn(h-l) + l }
func randomOperator(o []string) string { return o[rand.Intn(len(o))] }
