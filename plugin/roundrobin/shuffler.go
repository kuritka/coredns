package roundrobin

import (
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

type shuffler interface {
	// Shuffle runs round-robin algorithm.
	// stateless contains incoming request while *msg is response modified by other plugins
	Shuffle(req request.Request, msg *dns.Msg) []dns.RR
}
