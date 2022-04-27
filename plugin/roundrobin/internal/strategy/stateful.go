package strategy

import (
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

type Stateful struct {
	state *stateful
}

func NewStateful() *Stateful {
	return &Stateful{
		state: newStateful(),
	}
}

func (s *Stateful) Shuffle(req request.Request, res *dns.Msg) (rr []dns.RR) {
	rr, _ = s.state.handle(&req, res)
	return
}
