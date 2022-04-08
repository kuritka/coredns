package roundrobin

import (
	"github.com/miekg/dns"
)

type Consistent struct {
}

func NewConsistent() *Consistent {
	return &Consistent{}
}

func (r *Consistent) Shuffle(msg *dns.Msg) []dns.RR{
	return msg.Answer
}
