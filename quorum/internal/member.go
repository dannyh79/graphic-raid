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
	Ack           MessageType = 0
	ReadyToElect  MessageType = 1
	WantToLead    MessageType = 2
	AckCandidate  MessageType = 3
	PromoteLeader MessageType = 4
)

type BoardMember struct {
	Id                string
	Writer            io.Writer
	WantToLead        func() bool
	Mailbox           <-chan Message
	Allcast           chan Message
	ControllerMailbox chan Message
	boardSize         int
	peers             map[string]bool
	isLeader          bool
}

type MemberParams struct {
	Id                string
	Writer            io.Writer
	WantToLead        func() bool
	Mailbox           <-chan Message
	Allcast           chan Message
	ControllerMailbox chan Message
	BoardSize         int
}

func NewMember(p MemberParams) {
	prefix := fmt.Sprintf("Member %s: ", p.Id)
	pl := func(s string) { fmt.Fprintln(p.Writer, s) }

	m := BoardMember{
		Id:                p.Id,
		WantToLead:        p.WantToLead,
		Mailbox:           p.Mailbox,
		Allcast:           p.Allcast,
		ControllerMailbox: p.ControllerMailbox,
		boardSize:         p.BoardSize,
		peers:             make(map[string]bool),
		isLeader:          false,
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
				m.Allcast <- Message{AckCandidate, m.Id, msg.From}
			case PromoteLeader:
				m.ControllerMailbox <- Message{Ack, m.Id, ""}
				m.isLeader = true
			}
		}
	}
}

func FiftyFifty() bool { return rand.Intn(2) == 0 }

const ControllerId = "controller"

type role int

const (
	follower role = 0
	leader   role = 1
)

type Controller struct {
	id         string
	members    map[string]role
	quorumRule func(id string, polls map[string]int, members int) (string, bool)
	mailbox    chan Message
	mailboxes  map[string]chan Message
	allcast    chan Message
	candidates map[string]int
	pl         func(string)
}

func (c *Controller) recordPolls(msg Message) {
	switch msg.T {
	case AckCandidate:
		cId := msg.Body
		c.candidates[cId]++
		c.maybePromoteLeader(cId)
	}
}

func (c *Controller) maybePromoteLeader(id string) {
	if reason, ok := c.quorumRule(id, c.candidates, len(c.members)); ok {
		c.mailboxes[id] <- Message{PromoteLeader, c.id, reason}
		<-c.mailbox
		c.members[id] = leader
		c.candidates = make(map[string]int)
		c.pl(reason)
	}
}

type ControllerParams struct {
	Writer     io.Writer
	Members    []string
	QuorumRule func(id string, polls map[string]int, members int) (string, bool)
	Mailbox    chan Message
	Mailboxes  map[string]chan Message
	Allcast    chan Message
}

func NewController(p ControllerParams) {
	members := make(map[string]role)
	for _, id := range p.Members {
		members[id] = follower
	}

	c := Controller{
		id:         ControllerId,
		members:    members,
		quorumRule: p.QuorumRule,
		mailbox:    p.Mailbox,
		mailboxes:  p.Mailboxes,
		allcast:    p.Allcast,
		candidates: make(map[string]int),
		pl:         func(s string) { fmt.Fprintln(p.Writer, s) },
	}

	for {
		select {
		case msg := <-c.allcast:
			for id, mbx := range c.mailboxes {
				if id != msg.From {
					mbx <- msg
				}
			}

			c.recordPolls(msg)
		}
	}
}

func PromoteFirstReachedHalfVotes(id string, polls map[string]int, members int) (reason string, ok bool) {
	if polls[id] >= members/2 {
		return fmt.Sprintf("Member %s voted to be leader", id), true
	}
	return fmt.Sprintf("Quorum failed: (1 < %v/2)", members), false
}
