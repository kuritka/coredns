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

func (s *Subnet) Shuffle(req request.Request, msg *dns.Msg) []dns.RR {
	_ = s.state.handle(&req)
	return nil
}
