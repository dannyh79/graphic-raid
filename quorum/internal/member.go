package board

import (
	"fmt"
	"io"
	"math/rand"
)

type BoardMember struct {
	Id           string
	Writer       io.Writer
	WantToLead   func() bool
	Ack          chan string
	ReadyToElect chan bool
}

type MemberParams struct {
	Id           string
	Writer       io.Writer
	WantToLead   func() bool
	Ack          chan string
	ReadyToElect chan bool
}

func NewMember(p MemberParams) {
	prefix := fmt.Sprintf("Member %s: ", p.Id)
	pl := func(s string) { fmt.Fprintln(p.Writer, s) }

	m := BoardMember{
		WantToLead:   p.WantToLead,
		Ack:          p.Ack,
		ReadyToElect: p.ReadyToElect,
	}

	pl(prefix + "Hi")
	m.Ack <- p.Id

	select {
	case <-m.ReadyToElect:
		if m.WantToLead() {
			pl(prefix + "I want to be leader")
		}
	}
}

func FiftyFifty() bool { return rand.Intn(2) == 0 }

type Controller struct {
	members      int
	Ack          chan string
	ReadyToElect chan bool
}

type ControllerParams struct {
	Members      int
	Ack          chan string
	ReadyToElect chan bool
}

func (c *Controller) allReadyToElect() {
	c.ReadyToElect <- true
}

func (c *Controller) start() {
	i := 0
	for range c.Ack {
		i++
		if i == c.members {
			c.allReadyToElect()
		}
	}
}

func NewController(p ControllerParams) *Controller {
	c := Controller{
		members:      p.Members,
		ReadyToElect: p.ReadyToElect,
		Ack:          p.Ack,
	}

	go c.start()

	return &c
}
