package domain

type Sentence string

type Message struct {
	Type Sentence
	To   string
}
