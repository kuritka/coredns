package roundrobin

import (
	"github.com/coredns/coredns/request"
	"math/rand"
	"time"

	"github.com/miekg/dns"
)

type Random struct {
}

func NewRandom() *Random {
	return &Random{}
}

func (r *Random) Shuffle(req request.Request, msg *dns.Msg) []dns.RR{
	var shuffled []dns.RR
	var skipped []dns.RR
	for _, a := range  msg.Answer {
		switch a.Header().Rrtype {
		case dns.TypeA, dns.TypeAAAA:
			shuffled = append(shuffled, a)
		default:
			skipped = append(skipped, a)
		}
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })
	return append(shuffled,skipped...)
}
