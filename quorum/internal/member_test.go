package board_test

import (
	"sync"

	b "github.com/dannyh79/graphic-raid/quorum/internal"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var (
	buf *gbytes.Buffer
	wg  sync.WaitGroup
	m   sync.Mutex
)

var _ = Describe("NewMember", func() {
	var (
		n       int
		p       b.MemberParams
		mailbox chan b.Message
		allcast chan b.Message
		cMbx    chan b.Message
	)

	BeforeEach(func() {
		buf = gbytes.NewBuffer()

		n = 2
		mailbox = make(chan b.Message, 2)
		allcast = make(chan b.Message, 2)
		cMbx = make(chan b.Message, 1)
		p = b.MemberParams{
			Id:                "0",
			WaitGroup:         &wg,
			Writer:            buf,
			Mutex:             &m,
			WantToLead:        func() bool { return true },
			Mailbox:           mailbox,
			ControllerMailbox: cMbx,
			Allcast:           allcast,
			BoardSize:         n,
		}
	})

	It(`Writes "Member 0: Hi" to buffer upon starting`, func() {
		go b.NewMember(p)

		Eventually(buf).Should(gbytes.Say("Member 0: Hi\n"))
	})

	It(`Writes "Member 0: I want to be leader" to buffer`, func() {
		go b.NewMember(p)

		mailbox <- b.Message{b.Ack, "1", ""}

		Eventually(buf).Should(gbytes.Say("Member 0: I want to be leader\n"))
	})

	It(`Does NOT write "Member 0: I want to be leader" to buffer`, func() {
		p.WantToLead = func() bool { return false }

		go b.NewMember(p)

		mailbox <- b.Message{b.Ack, "1", ""}

		Eventually(buf).ShouldNot(gbytes.Say("Member 0: I want to be leader\n"))
	})

	It(`Writes "Member 0: Accept member 1 to be leader" to buffer`, func() {
		p.WantToLead = func() bool { return false }

		go b.NewMember(p)

		mailbox <- b.Message{b.Ack, "1", ""}
		mailbox <- b.Message{b.WantToLead, "1", ""}

		Eventually(buf).Should(gbytes.Say("Member 0: Accept member 1 to be leader\n"))
	})

	It(`Writes "Member 0: failed heartbeat with Member 1"`, func() {
		go b.NewMember(p)

		mailbox <- b.Message{b.Ack, "1", ""}
		mailbox <- b.Message{b.KeepAliveFail, b.ControllerId, "1"}

		Eventually(buf).Should(gbytes.Say("Member 0: failed heartbeat with Member 1\n"))
	})

	It(`Sends KeepAliveStart message to others upon promoted`, func() {
		go b.NewMember(p)

		mailbox <- b.Message{b.Ack, "1", ""}
		mailbox <- b.Message{b.PromoteLeader, "1", ""}

		Eventually(allcast).Should(Receive(Equal(b.Message{b.KeepAliveStart, p.Id, ""})))
	})

	It(`Sends KeepAlive message to others upon receiving KeepAliveStart`, func() {
		go b.NewMember(p)

		mailbox <- b.Message{b.Ack, "1", ""}
		mailbox <- b.Message{b.KeepAliveStart, "1", ""}

		Eventually(allcast).Should(Receive(Equal(b.Message{b.KeepAlive, p.Id, ""})))
	})

	It(`Sends KeepAlive message to others upon receiving KeepAlive`, func() {
		p.WantToLead = func() bool { return false }
		go b.NewMember(p)

		mailbox <- b.Message{b.Ack, "1", ""}
		mailbox <- b.Message{b.KeepAlive, "1", ""}

		Eventually(allcast).Should(Receive(Equal(b.Message{b.KeepAlive, p.Id, ""})))
	})

	It(`Does NOT send KeepAlive after receiving Kill message`, func() {
		wg.Add(1)
		go b.NewMember(p)

		mailbox <- b.Message{b.Ack, "1", ""}
		mailbox <- b.Message{b.KeepAliveStart, "1", ""}

		Eventually(allcast).Should(Receive(Equal(b.Message{b.KeepAlive, p.Id, ""})))

		mailbox <- b.Message{b.Kill, b.ControllerId, ""}

		Consistently(allcast).ShouldNot(Receive(Equal(b.Message{b.KeepAlive, p.Id, ""})))

		wg.Wait()
	})

	It(`Writes to buffer in a sequential manner`, func() {
		go b.NewMember(p)

		mailbox <- b.Message{b.Ack, "1", ""}
		mailbox <- b.Message{b.WantToLead, "1", ""}

		Eventually(buf).Should(gbytes.Say("Member 0: Hi\n"))
		Eventually(buf).Should(gbytes.Say("Member 0: I want to be leader\n"))
		Eventually(buf).Should(gbytes.Say("Member 0: Accept member 1 to be leader\n"))
	})
})

var _ = Describe("NewController", func() {
	var (
		allcast   chan b.Message
		mailboxes map[string]chan b.Message
		cMbx      chan b.Message
		p         b.ControllerParams
	)

	ids := []string{"0", "1", "2"}
	n := len(ids)

	BeforeEach(func() {
		buf = gbytes.NewBuffer()

		allcast = make(chan b.Message, n)
		mailboxes = make(map[string]chan b.Message)
		cMbx = make(chan b.Message)

		for _, id := range []string{"0", "1", "2"} {
			mailboxes[id] = make(chan b.Message, n)
		}

		p = b.ControllerParams{
			Writer:     buf,
			Members:    ids,
			QuorumRule: b.AgreeOnHalfVotes,
			Mailboxes:  mailboxes,
			Allcast:    allcast,
			Mailbox:    cMbx,
		}
	})

	It(`Broadcasts message to other mailboxes`, func() {
		go b.NewController(p)

		allcast <- b.Message{b.Ack, "0", ""}

		Eventually(mailboxes["1"]).Should(Receive(Equal(b.Message{b.Ack, "0", ""})))
		Eventually(mailboxes["2"]).Should(Receive(Equal(b.Message{b.Ack, "0", ""})))
	})

	It(`Sends KeepAliveFail message to sender`, func() {
		go b.NewController(p)

		mailboxes["0"] <- b.Message{b.Kill, b.ControllerId, ""}
		for range n + 1 {
			allcast <- b.Message{b.KeepAlive, "1", ""}
		}

		getMessageType := func(m b.Message) b.MessageType {
			return m.T
		}
		Eventually(mailboxes["1"]).Should(Receive(WithTransform(getMessageType, Equal(b.KeepAliveFail))))
	})
})

var _ = Describe("Interaction between 3 BoardMembers", func() {
	var (
		allcast    chan b.Message
		mailboxes  map[string]chan b.Message
		p1, p2, p3 b.MemberParams
		cp         b.ControllerParams
		cMbx       chan b.Message
	)

	ids := []string{"0", "1", "2"}
	n := len(ids)

	BeforeEach(func() {
		buf = gbytes.NewBuffer()

		allcast = make(chan b.Message, n)
		mailboxes = make(map[string]chan b.Message)
		for _, id := range ids {
			mailboxes[id] = make(chan b.Message, 1)
		}
		cMbx = make(chan b.Message, 1)

		p1 = b.MemberParams{
			Id:                ids[0],
			WaitGroup:         &wg,
			Writer:            buf,
			Mutex:             &m,
			WantToLead:        func() bool { return true },
			Mailbox:           mailboxes[ids[0]],
			Allcast:           allcast,
			ControllerMailbox: cMbx,
			BoardSize:         n,
		}
		p2 = b.MemberParams{
			Id:                ids[1],
			WaitGroup:         &wg,
			Writer:            buf,
			Mutex:             &m,
			WantToLead:        func() bool { return true },
			Mailbox:           mailboxes[ids[1]],
			Allcast:           allcast,
			ControllerMailbox: cMbx,
			BoardSize:         n,
		}
		p3 = b.MemberParams{
			Id:                ids[2],
			WaitGroup:         &wg,
			Writer:            buf,
			Mutex:             &m,
			WantToLead:        func() bool { return false },
			Mailbox:           mailboxes[ids[2]],
			Allcast:           allcast,
			ControllerMailbox: cMbx,
			BoardSize:         n,
		}

		cp = b.ControllerParams{
			Writer:     buf,
			Mutex:      &m,
			Members:    ids,
			QuorumRule: b.AgreeOnHalfVotes,
			Mailboxes:  mailboxes,
			Allcast:    allcast,
			Mailbox:    cMbx,
		}
	})

	It(`Writes to buffer in a sequential manner`, func() {
		go b.NewController(cp)

		for _, p := range []b.MemberParams{p1, p2, p3} {
			go b.NewMember(p)
		}

		Eventually(buf).MustPassRepeatedly(3).Should(gbytes.Say("Member [012]: Hi\n"))
		Eventually(buf).MustPassRepeatedly(2).Should(gbytes.Say("Member [01]: I want to be leader\n"))
		Eventually(buf).MustPassRepeatedly(3).Should(gbytes.Say("Member [012]: Accept member [01] to be leader\n"))
	})
})
