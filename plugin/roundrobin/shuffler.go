package roundrobin

import "github.com/miekg/dns"

type shuffler interface {
	Shuffle(answer []dns.RR) []dns.RR
}
