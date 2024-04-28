package board_test

import (
	board "github.com/dannyh79/graphic-raid/quorum/internal"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var (
	buf *gbytes.Buffer
)

var _ = Describe("NewMember", func() {
	var (
		n       int
		p       board.MemberParams
		mailbox chan board.Message
		allcast chan board.Message
	)

	BeforeEach(func() {
		buf = gbytes.NewBuffer()

		n = 2
		mailbox = make(chan board.Message, 2)
		allcast = make(chan board.Message, 2)
		p = board.MemberParams{
			Id:         "0",
			Writer:     buf,
			WantToLead: func() bool { return true },
			Mailbox:    mailbox,
			Allcast:    allcast,
			BoardSize:  n,
		}
	})

	It(`Writes "Member 0: Hi" to buffer upon starting`, func() {
		go board.NewMember(p)

		Eventually(buf).Should(gbytes.Say("Member 0: Hi\n"))
	})

	It(`Writes "Member 0: I want to be leader" to buffer`, func() {
		go board.NewMember(p)

		mailbox <- board.Message{board.Ack, "1", ""}

		Eventually(buf).Should(gbytes.Say("Member 0: I want to be leader\n"))
	})

	It(`Does NOT write "Member 0: I want to be leader" to buffer`, func() {
		p.WantToLead = func() bool { return false }

		go board.NewMember(p)

		mailbox <- board.Message{board.Ack, "1", ""}

		Eventually(buf).ShouldNot(gbytes.Say("Member 0: I want to be leader\n"))
	})

	It(`Writes "Member 0: Accept member 1 to be leader" to buffer`, func() {
		p.WantToLead = func() bool { return false }

		go board.NewMember(p)

		go func() {
			mailbox <- board.Message{board.Ack, "1", ""}
			mailbox <- board.Message{board.WantToLead, "1", ""}
		}()

		Eventually(buf).Should(gbytes.Say("Member 0: Accept member 1 to be leader\n"))
	})

	It(`Writes to buffer in a sequential manner`, func() {
		go board.NewMember(p)

		go func() {
			mailbox <- board.Message{board.Ack, "1", ""}
			mailbox <- board.Message{board.WantToLead, "1", ""}
		}()

		Eventually(buf).Should(gbytes.Say("Member 0: Hi\n"))
		Eventually(buf).Should(gbytes.Say("Member 0: I want to be leader\n"))
		Eventually(buf).Should(gbytes.Say("Member 0: Accept member 1 to be leader\n"))
	})
})

var _ = Describe("NewController", func() {
	var (
		n         int
		allcast   chan board.Message
		mailboxes map[string]chan board.Message
		p         board.ControllerParams
	)

	BeforeEach(func() {
		n = 3

		allcast = make(chan board.Message, n)
		mailboxes = make(map[string]chan board.Message)

		for _, id := range []string{"0", "1", "2"} {
			mailboxes[id] = make(chan board.Message, 1)
		}

		p = board.ControllerParams{
			Members:   n,
			Mailboxes: mailboxes,
			Allcast:   allcast,
		}
	})

	It(`Broadcasts message to other mailboxes`, func() {
		go board.NewController(p)

		go func() {
			allcast <- board.Message{board.Ack, "0", ""}
		}()

		Eventually(mailboxes["1"]).Should(Receive(Equal(board.Message{board.Ack, "0", ""})))
		Eventually(mailboxes["2"]).Should(Receive(Equal(board.Message{board.Ack, "0", ""})))
	})
})

// var _ = Describe("Interaction between 2 BoardMembers", func() {
// 	var (
// 		p1, p2 board.MemberParams
// 		cp     board.ControllerParams
// 	)
//
// 	BeforeEach(func() {
// 		buf = gbytes.NewBuffer()
// 		ack = make(chan string)
// 		readyToElect = make(chan bool)
// 		candidate = make(chan string)
// 		ackCandidate = make(chan string)
//
// 		cp = board.ControllerParams{
// 			Members:      2,
// 			Ack:          ack,
// 			ReadyToElect: readyToElect,
// 		}
// 	})
//
// 	Context("when both wants to be the leader", func() {
// 		BeforeEach(func() {
// 			p1 = board.MemberParams{
// 				Id:           "0",
// 				Writer:       buf,
// 				WantToLead:   func() bool { return true },
// 				Ack:          ack,
// 				ReadyToElect: readyToElect,
// 				Candidate:    candidate,
// 				AckCandidate: ackCandidate,
// 			}
// 			p2 = board.MemberParams{
// 				Id:           "1",
// 				Writer:       buf,
// 				WantToLead:   func() bool { return true },
// 				Ack:          ack,
// 				ReadyToElect: readyToElect,
// 				Candidate:    candidate,
// 				AckCandidate: ackCandidate,
// 			}
// 		})
//
// 		It(`Does NOT write "Member [01]: Accept member [01] to be leader" to buffer`, func() {
// 			go board.NewController(cp)
// 			go board.NewMember(p1)
// 			go board.NewMember(p2)
//
// 			Eventually(buf, 2).Should(gbytes.Say("Member [01]: Hi\n"))
// 			Eventually(buf, 2).Should(gbytes.Say("Member [01]: I want to be leader\n"))
// 			Eventually(buf, 2).Should(Not(gbytes.Say("Member [01]: Accept member [01] to be leader\n")))
// 		})
// 	})

// 	Context("when ONLY one wants to be the leader", func() {
// 		BeforeEach(func() {
// 			p1 = board.MemberParams{
// 				Id:           "0",
// 				Writer:       buf,
// 				WantToLead:   func() bool { return true },
// 				Ack:          ack,
// 				ReadyToElect: readyToElect,
// 				Candidate:    candidate,
// 				AckCandidate: ackCandidate,
// 			}
// 			p2 = board.MemberParams{
// 				Id:           "1",
// 				Writer:       buf,
// 				WantToLead:   func() bool { return false },
// 				Ack:          ack,
// 				ReadyToElect: readyToElect,
// 				Candidate:    candidate,
// 				AckCandidate: ackCandidate,
// 			}
// 		})
//
// 		It(`Writes "Member 1: Accept member 0 to be leader" to buffer`, func() {
// 			go board.NewController(cp)
// 			go board.NewMember(p1)
// 			go board.NewMember(p2)
//
// 			Eventually(buf, 2).Should(gbytes.Say("Member [01]: Hi\n"))
// 			Eventually(buf).Should(gbytes.Say("Member 0: I want to be leader\n"))
// 			Eventually(buf).Should(gbytes.Say("Member 1: Accept member 0 to be leader\n"))
// 		})
// 	})
//
// 	Context("when no one wants to be the leader", func() {
// 		BeforeEach(func() {
// 			p1 = board.MemberParams{
// 				Id:           "0",
// 				Writer:       buf,
// 				WantToLead:   func() bool { return false },
// 				Ack:          ack,
// 				ReadyToElect: readyToElect,
// 				Candidate:    candidate,
// 				AckCandidate: ackCandidate,
// 			}
// 			p2 = board.MemberParams{
// 				Id:           "1",
// 				Writer:       buf,
// 				WantToLead:   func() bool { return false },
// 				Ack:          ack,
// 				ReadyToElect: readyToElect,
// 				Candidate:    candidate,
// 				AckCandidate: ackCandidate,
// 			}
// 		})
//
// 		It(`Does NOT write "Member [01]: Accept member 0 to be leader" to buffer`, func() {
// 			go board.NewController(cp)
// 			go board.NewMember(p1)
// 			go board.NewMember(p2)
//
// 			Eventually(buf, 2).Should(gbytes.Say("Member [01]: Hi\n"))
// 			Eventually(buf).Should(Not(gbytes.Say("Member [01]: I want to be leader\n")))
// 			Eventually(buf, 2).Should(Not(gbytes.Say("Member [01]: Accept member [01] to be leader\n")))
// 		})
// 	})
// })
//
// var _ = Describe("Interaction between 3 BoardMembers", func() {
// 	var (
// 		p1, p2, p3 board.MemberParams
// 		cp         board.ControllerParams
// 	)
//
// 	BeforeEach(func() {
// 		buf = gbytes.NewBuffer()
// 		ack = make(chan string)
// 		readyToElect = make(chan bool)
// 		candidate = make(chan string)
// 		ackCandidate = make(chan string)
//
// 		p1 = board.MemberParams{
// 			Id:           "0",
// 			Writer:       buf,
// 			WantToLead:   func() bool { return true },
// 			Ack:          ack,
// 			ReadyToElect: readyToElect,
// 			Candidate:    candidate,
// 			AckCandidate: ackCandidate,
// 		}
// 		p2 = board.MemberParams{
// 			Id:           "1",
// 			Writer:       buf,
// 			WantToLead:   func() bool { return true },
// 			Ack:          ack,
// 			ReadyToElect: readyToElect,
// 			Candidate:    candidate,
// 			AckCandidate: ackCandidate,
// 		}
// 		p3 = board.MemberParams{
// 			Id:           "2",
// 			Writer:       buf,
// 			WantToLead:   func() bool { return false },
// 			Ack:          ack,
// 			ReadyToElect: readyToElect,
// 			Candidate:    candidate,
// 			AckCandidate: ackCandidate,
// 		}
// 		cp = board.ControllerParams{
// 			Members:      3,
// 			Ack:          ack,
// 			ReadyToElect: readyToElect,
// 		}
// 	})
//
// 	Context("when 2 BoardMembers wanting to be the leader", func() {
// 		It(`Writes "Member [012]: Accept member [01] to be leader" to buffer`, func() {
// 			go board.NewController(cp)
//
// 			for _, p := range []board.MemberParams{p1, p2, p3} {
// 				go board.NewMember(p)
// 			}
//
// 			Eventually(buf, 3).Should(gbytes.Say("Member [012]: Hi\n"))
// 			Eventually(buf, 2).Should(gbytes.Say("Member [01]: I want to be leader\n"))
// 			Eventually(buf, 2).Should(gbytes.Say("Member [012]: Accept member [01] to be leader\n"))
// 		})
// 	})
// })
