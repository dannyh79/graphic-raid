package board_test

import (
	"strconv"

	board "github.com/dannyh79/graphic-raid/quorum/internal"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var (
	buf          *gbytes.Buffer
	ack          chan string
	readyToElect chan bool
)

var _ = Describe("NewMember", func() {
	var (
		p board.MemberParams
	)

	BeforeEach(func() {
		buf = gbytes.NewBuffer()

		ack = make(chan string)
		readyToElect = make(chan bool)
		p = board.MemberParams{
			Id:           "0",
			Writer:       buf,
			WantToLead:   func() bool { return true },
			Ack:          ack,
			ReadyToElect: readyToElect,
		}
	})

	It(`Writes "Member 0: Hi" to buffer upon starting`, func() {
		go board.NewMember(p)

		Eventually(buf).Should(gbytes.Say("Member 0: Hi\n"))
	})

	It(`Writes "Member 0: I want to be leader" to buffer`, func() {
		go board.NewMember(p)

		<-ack
		readyToElect <- true

		Eventually(buf).Should(gbytes.Say("Member 0: I want to be leader\n"))
	})

	It(`does NOT write "Member 0: I want to be leader" to buffer`, func() {
		p.WantToLead = func() bool { return false }

		go board.NewMember(p)

		<-ack
		readyToElect <- true

		Eventually(buf).Should(Not(gbytes.Say("Member 0: I want to be leader\n")))
	})

	It(`Writes to buffer in a sequential manner`, func() {
		go board.NewMember(p)

		<-ack
		readyToElect <- true

		Eventually(buf).Should(gbytes.Say("Member 0: Hi\n"))
		Eventually(buf).Should(gbytes.Say("Member 0: I want to be leader\n"))
	})
})

var _ = Describe("NewController", func() {
	var (
		n int
		c *board.Controller
	)

	It(`Sends to c.ReadyToElect`, func() {
		n = 2
		p := board.ControllerParams{
			Members:      n,
			Ack:          make(chan string),
			ReadyToElect: make(chan bool),
		}

		c = board.NewController(p)

		go func() {
			for i := 0; i < n; i++ {
				c.Ack <- strconv.Itoa(i)
			}
		}()

		Eventually(<-c.ReadyToElect).Should(BeTrue())
	})
})

var _ = Describe("Interaction between 2 BoardMembers", func() {
	var (
		p1, p2       board.MemberParams
		cp           board.ControllerParams
		ack          chan string
		readyToElect chan bool
	)

	BeforeEach(func() {
		buf = gbytes.NewBuffer()
		ack = make(chan string)
		readyToElect = make(chan bool)

		p1 = board.MemberParams{
			Id:           "0",
			Writer:       buf,
			WantToLead:   func() bool { return true },
			Ack:          ack,
			ReadyToElect: readyToElect,
		}
		p2 = board.MemberParams{
			Id:           "1",
			Writer:       buf,
			WantToLead:   func() bool { return true },
			Ack:          ack,
			ReadyToElect: readyToElect,
		}
		cp = board.ControllerParams{
			Members:      2,
			Ack:          ack,
			ReadyToElect: readyToElect,
		}
	})

	Context("when both wants to be the leader", func() {
		BeforeEach(func() {
			p1 = board.MemberParams{
				Id:           "0",
				Writer:       buf,
				WantToLead:   func() bool { return true },
				Ack:          ack,
				ReadyToElect: readyToElect,
			}
			p2 = board.MemberParams{
				Id:           "1",
				Writer:       buf,
				WantToLead:   func() bool { return true },
				Ack:          ack,
				ReadyToElect: readyToElect,
			}
		})

		It(`does NOT write "Member [01]: Accept member [01] to be leader" to buffer`, func() {
			go board.NewController(cp)
			go board.NewMember(p1)
			go board.NewMember(p2)

			Eventually(buf, 2).Should(gbytes.Say("Member [01]: Hi\n"))
			Eventually(buf, 2).Should(gbytes.Say("Member [01]: I want to be leader\n"))
			Eventually(buf).Should(Not(gbytes.Say("Member [01]: Accept member [01] to be leader\n")))
		})
	})
})
