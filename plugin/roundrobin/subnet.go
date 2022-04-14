package roundrobin

import (
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

type Subnet struct {
	 state *stateful
}

func NewSubnet() *Subnet {
	return &Subnet{
		state: newStateful(),
	}
}

func (s *Subnet) Shuffle(req request.Request, res *dns.Msg) (rr []dns.RR) {
	rr, _ = s.state.handle(&req, res)
	return
}
