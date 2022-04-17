package roundrobin

import "github.com/miekg/dns"

// parseAnswerSection converts []dns.RR into map of A or AAAA records and slice containing all except A or AAAA
// todo: reuse in random / stateless
func parseAnswerSection(arr []dns.RR) (ip map[string]dns.RR, noip []dns.RR) {
	ip = make(map[string]dns.RR)
	noip = make([]dns.RR,0)
	for _, r := range arr {
		switch r.Header().Rrtype {
		case dns.TypeA:
			ip[r.(*dns.A).A.String()] = r
		case dns.TypeAAAA:
			ip[r.(*dns.AAAA).AAAA.String()] = r
		default:
			noip = append(noip, r)
		}
	}
	return
}

// ipsToSet converts list of IPs into set of IP's
func ipsToSet(ips []string) (m map[string]bool) {
	m = make(map[string]bool)
	for _, ip := range ips {
		m[ip] = true
	}
	return
}

func rotate(slice []string) (r []string){
	l := len(slice)
	for i := range slice {
		r = append(r, slice[(i+1) % l])
	}
	return
}

//func rotate(slice []string) (r []string){
//	return append(slice[(len(slice)):],slice[0:len(slice)]...)
//}