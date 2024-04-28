package board

import (
	"fmt"
	"io"
	"math/rand"
)

type MessageType int
type Body string

type Message struct {
	T    MessageType
	From string
	Body string
}

const (
	Ack          MessageType = 0
	ReadyToElect MessageType = 1
	WantToLead   MessageType = 2
	AckCandidate MessageType = 3
)

type BoardMember struct {
	Id         string
	Writer     io.Writer
	WantToLead func() bool
	Mailbox    <-chan Message
	Allcast    chan Message
	boardSize  int
	peers      map[string]bool
}

type MemberParams struct {
	Id         string
	Writer     io.Writer
	WantToLead func() bool
	Mailbox    <-chan Message
	Allcast    chan Message
	BoardSize  int
}

func NewMember(p MemberParams) {
	prefix := fmt.Sprintf("Member %s: ", p.Id)
	pl := func(s string) { fmt.Fprintln(p.Writer, s) }

	m := BoardMember{
		Id:         p.Id,
		WantToLead: p.WantToLead,
		Mailbox:    p.Mailbox,
		Allcast:    p.Allcast,
		boardSize:  p.BoardSize,
		peers:      make(map[string]bool),
	}

	pl(prefix + "Hi")
	m.Allcast <- Message{Ack, m.Id, ""}

	for {
		select {
		case msg := <-m.Mailbox:
			switch msg.T {
			case Ack:
				m.peers[msg.From] = true
			}
		}
		if len(m.peers) == m.boardSize-1 {
			break
		}
	}

	if m.WantToLead() {
		pl(prefix + "I want to be leader")
		m.Allcast <- Message{WantToLead, m.Id, ""}
	}

	for {
		select {
		case msg := <-m.Mailbox:
			switch msg.T {
			case WantToLead:
				pl(prefix + fmt.Sprintf("Accept member %s to be leader", msg.From))
				m.Allcast <- Message{AckCandidate, m.Id, ""}
			}
		}
	}
}

func FiftyFifty() bool { return rand.Intn(2) == 0 }

type Controller struct {
	members   int
	mailboxes map[string]chan Message
	allcast   chan Message
}

type ControllerParams struct {
	Members   int
	Mailboxes map[string]chan Message
	Allcast   chan Message
}

func NewController(p ControllerParams) {
	c := Controller{
		members:   p.Members,
		mailboxes: p.Mailboxes,
		allcast:   p.Allcast,
	}

	for {
		select {
		case msg := <-c.allcast:
			for id, mbx := range c.mailboxes {
				if id != msg.From {
					mbx <- msg
				}
			}
		}
	}
}
