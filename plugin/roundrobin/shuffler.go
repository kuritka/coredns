package roundrobin

import "github.com/miekg/dns"

type shuffler interface {
	Shuffle(msg *dns.Msg) []dns.RR
}
