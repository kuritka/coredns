package roundrobin

import (
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

type Consistent struct {
}

func NewConsistent() *Consistent {
	return &Consistent{}
}

func (r *Consistent) Shuffle(req request.Request, msg *dns.Msg) []dns.RR{
	return newStateless(req.Req, msg).updateState().rotate().getAnswers()
}

