package roundrobin

import (
	"fmt"
	"github.com/coredns/coredns/plugin/test"
	"github.com/miekg/dns"
	"testing"
)

var (
	arotations =[]string{"[10.240.0.2 10.240.0.3 10.240.0.4 10.240.0.1]","[10.240.0.3 10.240.0.4 10.240.0.1 10.240.0.2]",
		"[10.240.0.4 10.240.0.1 10.240.0.2 10.240.0.3]", "[10.240.0.1 10.240.0.2 10.240.0.3 10.240.0.4]",
		"[10.240.0.1 10.240.0.2 10.240.0.3 10.240.0.4]", "[10.240.0.2 10.240.0.3 10.240.0.4 10.240.0.1]"}
)

func TestRoundRobinStatelessWithNilRequestState(t *testing.T) {
	// rr_state=nil was provided at the beginning of transaction
	// The CoreDNS doesn't know anything about client's state.
	var clientState []dns.RR
	expected := "[10.240.0.2 10.240.0.3 10.240.0.4 10.240.0.1]"
	m := newMid()

	m.AddResponseAnswer(test.CNAME("alpha.cloud.example.com.	300	IN	CNAME		beta.cloud.example.com."))
	m.AddResponseAnswer(test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.1"))
	m.AddResponseAnswer(test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.2"))
	m.AddResponseAnswer(test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.3"))
	m.AddResponseAnswer(test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.4"))
	m.AddResponseAnswer(test.MX("alpha.cloud.example.com.			300	IN	MX		1	mxa-alpha.cloud.example.com."))
	m.AddRequestOpt(clientState...)

	clientState = NewStateless().Shuffle(m.req, m.res)
	if fmt.Sprintf("%v", getIPs(clientState)) != expected {
		t.Errorf("The Stateless shuffle is not working as expected. Expecting %v but got %v.", expected, getIPs(clientState))
	}
}

func TestRoundRobinStatelessWithoutRequestState(t *testing.T) {
	// rr_state was not provided at the beginning of transaction
	// The CoreDNS doesn't know anything about client's state.
	var clientState []dns.RR
	expected := "[10.240.0.2 10.240.0.3 10.240.0.4 10.240.0.1]"
	m := newMid()

	m.AddResponseAnswer(test.CNAME("alpha.cloud.example.com.	300	IN	CNAME		beta.cloud.example.com."))
	m.AddResponseAnswer(test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.1"))
	m.AddResponseAnswer(test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.2"))
	m.AddResponseAnswer(test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.3"))
	m.AddResponseAnswer(test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.4"))
	m.AddResponseAnswer(test.MX("alpha.cloud.example.com.			300	IN	MX		1	mxa-alpha.cloud.example.com."))

	clientState = NewStateless().Shuffle(m.req, m.res)
	if fmt.Sprintf("%v", getIPs(clientState)) != expected {
		t.Errorf("The Stateless shuffle is not working as expected. Expecting %v but got %v.", expected, getIPs(clientState))
	}
}

func TestRoundRobinStatelessWithEmptyRequestState(t *testing.T) {

}

func TestRoundRobinStatelessWithInvalidRequestState(t *testing.T) {
	// rr_state=invalid was provided at the beginning of transaction
	// The CoreDNS doesn't know anything about client's state.
	var clientState []dns.RR
	expected := "[10.240.0.2 10.240.0.3 10.240.0.4 10.240.0.1]"
	m := newMid()

	m.AddResponseAnswer(test.CNAME("alpha.cloud.example.com.	300	IN	CNAME		beta.cloud.example.com."))
	m.AddResponseAnswer(test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.1"))
	m.AddResponseAnswer(test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.2"))
	m.AddResponseAnswer(test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.3"))
	m.AddResponseAnswer(test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.4"))
	m.AddResponseAnswer(test.MX("alpha.cloud.example.com.			300	IN	MX		1	mxa-alpha.cloud.example.com."))

	opt := new(dns.OPT)
	opt.Hdr.Name = "."
	opt.Hdr.Rrtype = dns.TypeOPT
	e := new(dns.EDNS0_LOCAL)
	e.Code = dns.EDNS0LOCALSTART
	e.Data = []byte("_rr_state=invalid")
	opt.Option = append(opt.Option, e)
	m.req.Req.Extra = append(m.req.Req.Extra, opt)

	clientState = NewStateless().Shuffle(m.req, m.res)
	if fmt.Sprintf("%v", getIPs(clientState)) != expected {
		t.Errorf("The Stateless shuffle is not working as expected. Expecting %v but got %v.", expected, getIPs(clientState))
	}

}


//func TestRoundRobinStatelessWithProvidedEmptyRequestState(t *testing.T) {
//	// rr_state=nil was provided at the beginning of transaction
//	// in stateless, the server doesn't know anything about client's record.
//	// the recorder order is sent back to the DNS server in each request
//	var clientState []dns.RR
//	for _, v := range arotations {
//		m := newMid()
//		m.AddResponseAnswer(test.CNAME("alpha.cloud.example.com.	300	IN	CNAME		beta.cloud.example.com."))
//		m.AddResponseAnswer(test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.1"))
//		m.AddResponseAnswer(test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.2"))
//		m.AddResponseAnswer(test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.3"))
//		m.AddResponseAnswer(test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.4"))
//		m.AddResponseAnswer(test.MX("alpha.cloud.example.com.			300	IN	MX		1	mxa-alpha.cloud.example.com."))
//		m.AddRequestOpt(clientState...)
//
//		clientState = NewStateless().Shuffle(m.req, m.res)
//		if fmt.Sprintf("%v", getIPs(clientState)) != v {
//			t.Errorf("The Stateless shuffle is not working as expected. Expecting %v but got %v.", v, clientState)
//		}
//	}
//}
