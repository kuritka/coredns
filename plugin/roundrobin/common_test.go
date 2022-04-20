package roundrobin

// common_test.go contains helper for round_robin tests

import (
	"context"
	"encoding/json"
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

func (p mid) AddRequestOpt(rr ...dns.RR){
	json, _ := json.Marshal(rr)
	opt := new(dns.OPT)
	opt.Hdr.Name = "."
	opt.Hdr.Rrtype = dns.TypeOPT
	e := new(dns.EDNS0_LOCAL)
	e.Code = dns.EDNS0LOCALSTART
	e.Data = append([]byte("_rr_state="),json...)
	opt.Option = append(opt.Option, e)
	p.req.Req.Extra = append(rr, opt)
}

func (p mid) AddRequestOptRaw(data string){
	opt := new(dns.OPT)
	opt.Hdr.Name = "."
	opt.Hdr.Rrtype = dns.TypeOPT
	e := new(dns.EDNS0_LOCAL)
	e.Code = dns.EDNS0LOCALSTART
	e.Data = append([]byte("_rr_state="),data)
	opt.Option = append(opt.Option, e)
	p.req.Req.Extra = append(rr, opt)
}


func getIPs(arr []dns.RR) (ips []string){
	ips = []string{}
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
