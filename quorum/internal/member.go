package board

import (
	"fmt"
	"io"
	"math/rand"
)

type BoardMember struct {
	Id         string
	Writer     io.Writer
	WantToLead func() bool
}

type MemberParams struct {
	Id         string
	Writer     io.Writer
	WantToLead func() bool
}

func NewMember(p MemberParams) {
	prefix := fmt.Sprintf("Member %s: ", p.Id)
	pl := func(s string) { fmt.Fprintln(p.Writer, s) }

	m := BoardMember{
		WantToLead: p.WantToLead,
	}

	pl(prefix + "Hi")

	if m.WantToLead() {
		pl(prefix + "I want to be leader")
	}
}

func FiftyFifty() bool { return rand.Intn(2) == 0 }
