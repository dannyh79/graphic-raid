package board

import (
	"fmt"
	"io"
)

type BoardMember struct {
	Id     string
	Writer io.Writer
}

type MemberParams struct {
	Id     string
	Writer io.Writer
}

func NewMember(m MemberParams) {
	p := func(s string) { fmt.Fprintln(m.Writer, s) }

	p(fmt.Sprintf("Member %s: Hi", m.Id))
}
