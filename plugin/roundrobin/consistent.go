package roundrobin

import "github.com/miekg/dns"

type Consistent struct {
}

func NewConsistent() *Consistent {
	return &Consistent{}
}

func (r *Consistent) Shuffle(answer []dns.RR) []dns.RR{
	return []dns.RR{}
}