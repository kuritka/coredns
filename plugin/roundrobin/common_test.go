package roundrobin

// common_test.go contains helper for round_robin tests

import (
	"context"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)


type mid struct {
	req request.Request
	res *dns.Msg
}

func newMid() mid {
	return mid{
		req: request.Request{
			Req: &dns.Msg{},
		},
		res: &dns.Msg{},
	}
}

func (p mid) AddResponseAnswer(rr dns.RR){
	p.res.Answer = append(p.res.Answer, rr)
}

func (p mid) AddResponseExtra(rr dns.RR){
	p.res.Extra = append(p.res.Answer, rr)
}

func (p mid) AddRequestAnswer(rr dns.RR){
	p.req.Req.Answer = append(p.req.Req.Answer, rr)
}

func (p mid) AddRequestExtra(rr dns.RR){
	p.req.Req.Extra = append(p.req.Req.Extra, rr)
}


func getIPs(arr []dns.RR) (ips []string){
	for _, rr := range arr {
		switch rr.Header().Rrtype {
		case dns.TypeA:
			ips = append(ips, rr.(*dns.A).A.String())
		case dns.TypeAAAA:
			ips = append(ips, rr.(*dns.AAAA).AAAA.String())
		}
	}
	return
}

func handler() plugin.Handler {
	return plugin.HandlerFunc(func(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
		w.WriteMsg(r)
		return dns.RcodeSuccess, nil
	})
}
