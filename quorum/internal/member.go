package board

import (
	"fmt"
	"io"
	"maps"
	"math/rand"
	"sync"
)

type MessageType int
type Body string

type Message struct {
	T    MessageType
	From string
	Body string
}

const (
	Ack            MessageType = 0
	ReadyToElect   MessageType = 1
	WantToLead     MessageType = 2
	AckCandidate   MessageType = 3
	PromoteLeader  MessageType = 4
	AckLeader      MessageType = 5
	KeepAliveStart MessageType = 6
	KeepAlive      MessageType = 7
	KeepAliveFail  MessageType = 8
	Kill           MessageType = 9
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
	leader            string
}

type MemberParams struct {
	Id                string
	WaitGroup         *sync.WaitGroup
	Writer            io.Writer
	WantToLead        func() bool
	Mailbox           <-chan Message
	Allcast           chan Message
	ControllerMailbox chan Message
	BoardSize         int
}

func NewMember(p MemberParams) {
	defer p.WaitGroup.Done()

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
		leader:            "",
	}

	pl(prefix + "Hi")
	m.Allcast <- Message{Ack, m.Id, ""}

	for {
		select {
		case msg, ok := <-m.Mailbox:
			if !ok {
				return
			}

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
				if m.leader != "" {
					continue
				}

				pl(prefix + fmt.Sprintf("Accept member %s to be leader", msg.From))
				m.Allcast <- Message{AckCandidate, m.Id, msg.From}
			case PromoteLeader:
				m.ControllerMailbox <- Message{Ack, m.Id, ""}
				m.Allcast <- Message{AckLeader, m.Id, ""}
				m.leader = m.Id
				m.Allcast <- Message{KeepAliveStart, m.Id, ""}
			case AckLeader:
				m.leader = msg.From
			case KeepAliveStart:
				m.Allcast <- Message{KeepAlive, m.Id, ""}
			case KeepAlive:
				m.Allcast <- Message{KeepAlive, m.Id, ""}
			case KeepAliveFail:
				pl(prefix + fmt.Sprintf("failed heartbeat with Member %s", msg.Body))
			case Kill:
				return
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
	m          *sync.Mutex
	members    map[string]role
	quorumRule func(id string, polls map[string]int, members int) (string, bool)
	mailbox    chan Message
	mailboxes  map[string]chan Message
	allcast    chan Message
	ballot     map[string]int
	pl         func(string)
}

func (c *Controller) recordPolls(msg Message) {
	switch msg.T {
	case AckCandidate:
		cId := msg.Body

		c.m.Lock()

		c.maybePromoteLeader(cId)

		c.m.Unlock()
	}
}

func (c *Controller) maybePromoteLeader(id string) {
	if c.hasLeader() {
		return
	}

	pendingPolls := maps.Clone(c.ballot)
	pendingPolls[id]++
	reason, ok := c.quorumRule(id, pendingPolls, len(c.members))
	if !ok {
		c.ballot[id]++
	}

	c.mailboxes[id] <- Message{PromoteLeader, c.id, reason}
	<-c.mailbox
	c.members[id] = leader
	c.ballot = make(map[string]int)
	c.pl(reason)
}

func (c *Controller) hasLeader() bool {
	for _, role := range c.members {
		if role == leader {
			return true
		}
	}
	return false
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
	var m sync.Mutex

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
		ballot:     make(map[string]int),
		pl:         func(s string) { fmt.Fprintln(p.Writer, s) },
		m:          &m,
	}

	for {
		select {
		case msg := <-c.allcast:
			for id, mbx := range c.mailboxes {
				if id != msg.From {
					mbx <- msg
					if len(mbx) >= cap(mbx) {
						c.mailboxes[msg.From] <- Message{KeepAliveFail, c.id, id}
					}
				}
			}

			c.recordPolls(msg)
		}
	}
}

func PromoteFirstReachedHalfVotes(id string, ballot map[string]int, members int) (reason string, ok bool) {
	if ballot[id] >= members/2 {
		return fmt.Sprintf("Member %s voted to be leader", id), true
	}
	return fmt.Sprintf("Quorum failed: (1 < %v/2)", members), false
}
