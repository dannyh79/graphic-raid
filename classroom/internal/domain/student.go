package domain

import (
	"fmt"
	"regexp"
	"strconv"
)

type Student struct {
	Name string
}

const (
	Answer        Sentence = "answer"
	questionRegex string   = `(\d+) (\+|-|\*|/) (\d+)`
)

func (s *Student) getName() string       { return fmt.Sprintf("Student %s", s.Name) }
func (s *Student) say(str string) string { return fmt.Sprintf("%s: %s", s.getName(), str) }

func (t *Student) answerQuiz(q string) string {
	re := regexp.MustCompile(questionRegex)
	m := re.FindStringSubmatch(q)
	if len(m) < 4 {
		panic("undefined spec: incorrect quiz format")
	}

	t1, op, t2 := parseFloat(m[1]), m[2], parseFloat(m[3])

	var a float64
	switch op {
	case "+":
		a = t1 + t2
	case "-":
		a = t1 - t2
	case "*":
		a = t1 * t2
	case "/":
		a = t1 / t2
	}
	return fmt.Sprintf("%v %s %v = %v!", t1, op, t2, a)
}

func (s *Student) Say(m Message) string {
	var str string
	switch m.Type {
	case Answer:
		str = s.answerQuiz(m.Body)
	case Respond:
		str = fmt.Sprintf("%s, you win.", m.To)
	}
	return s.say(str)
}

func NewStudent(n string) *Student { return &Student{n} }

func GetStudentName(s string) string {
	re := regexp.MustCompile(`Student (\w+):`)
	m := re.FindStringSubmatch(s)

	if len(m) > 1 {
		return m[1]
	}
	return ""
}

func parseFloat(n string) float64 {
	i, err := strconv.ParseFloat(n, 10)
	if err != nil {
		panic("undefined spec: handling parser failure")
	}

	return i
}
