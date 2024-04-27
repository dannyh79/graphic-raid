package board_test

import (
	"strconv"

	board "github.com/dannyh79/graphic-raid/quorum/internal"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("NewMember", func() {
	var (
		buf          *gbytes.Buffer
		p            board.MemberParams
		ack          chan string
		readyToElect chan bool
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
