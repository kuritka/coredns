package roundrobin

import (
	"fmt"
	"testing"

	"github.com/coredns/coredns/plugin/test"
	"github.com/miekg/dns"
)

// In stateless, the client sends back previous state to the CoreDNS per DNS query,
// the server doesn't know anything about client's data.
// This test examines the response, different situations when we send different
// types of data to the server (uninitialized data or valid or invalid) and tests that we get the expected response.
func TestRoundRobinStatelessInitialize(t *testing.T) {
	expected := "[10.240.0.2 10.240.0.3 10.240.0.4 10.240.0.1]"
	cname := test.CNAME("alpha.cloud.example.com.	300	IN	CNAME		beta.cloud.example.com.")
	mx := test.MX("alpha.cloud.example.com.			300	IN	MX		1	mxa-alpha.cloud.example.com.")
	testValues := []struct {
		value    string
		expected string
	}{
		{"_rr_state=invalid", expected},
		{`_rr_state=`, expected},
		{``, expected},
		{`_rr_state=nil`, expected},
		{`_rr_state={}`, expected},
		{`_rr_state={"ip":[]}`, expected},
		{`_rr_state={"ip":["10.240.0.1"]}`, expected},
		{`_rr_state={"ip":[10.240.0.2 10.240.0.3 10.240.0.4 10.240.0.1]}`, expected},
		{`_rr_state={"ip":[10.240.0.10 10.240.0.20 10.240.0.40 10.240.0.111]}`, expected},
		{`_rr_state={"ip":["10.0.0.1","10.2.2.1","10.1.1.2","10.1.1.3","10.2.2.2","10.2.2.3","10.0.0.2","10.0.0.3","10.1.1.1"]}`, expected},
		{`blah=`, expected},
	}
	for _, raw := range testValues {
		m := newMid()

		m.AddResponseAnswer(cname)
		m.AddResponseAnswer(test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.1"))
		m.AddResponseAnswer(test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.2"))
		m.AddResponseAnswer(test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.3"))
		m.AddResponseAnswer(test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.4"))
		m.AddResponseAnswer(mx)
		m.AddRequestOptRaw(raw.value)

		var clientState = NewStateless().Shuffle(m.req, m.res)

		if len(clientState) != len(m.res.Answer) {
			t.Errorf("The stateless retrieved different number of records. Expected %v got %v", len(m.res.Answer), len(clientState))
		}

		if fmt.Sprintf("%v", getIPs(clientState)) != raw.expected {
			t.Errorf("The stateless shuffle is not working as expected. For %s Expecting %v but got %v.", raw.value, expected, getIPs(clientState))
		}

		// end of the list contains additional records in the order of they are defined in the response
		if clientState[4].String() != cname.String() {
			t.Errorf("Expecting %s result but got %s", cname, clientState[4].String())
		}
		if clientState[5].String() != mx.String() {
			t.Errorf("Expecting %s result but got %s",mx, clientState[5].String())
		}
	}
}

func TestRoundRobinStatelessShuffleA(t *testing.T) {
	var arotations =[]string{"[10.240.0.2 10.240.0.3 10.240.0.4 10.240.0.1]","[10.240.0.3 10.240.0.4 10.240.0.1 10.240.0.2]",
		"[10.240.0.4 10.240.0.1 10.240.0.2 10.240.0.3]", "[10.240.0.1 10.240.0.2 10.240.0.3 10.240.0.4]"}
	var rr  []dns.RR
	for i := 0; i < 10; i++ {
		m := newMid()
		m.AddResponseAnswer(test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.1"))
		m.AddResponseAnswer(test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.2"))
		m.AddResponseAnswer(test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.3"))
		m.AddResponseAnswer(test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.4"))
		// state, from the previous loop that arrived in the DNS query
		m.AddRequestOpt(rr...)
		var clientState = NewStateless().Shuffle(m.req, m.res)

		// save the new state for the next query
		rr = filterAandAAAA(clientState)

		if fmt.Sprintf("%v",getIPs(clientState)) != fmt.Sprintf("%v",arotations[i %len(arotations)]) {
			t.Errorf("%v: The stateless rotation is not working. Expecting %v but got %v.",i, arotations[i %len(arotations)], getIPs(clientState))
		}
	}
}

func TestRoundRobinStatelessShuffleAAAA(t *testing.T) {
	var aaaarotations =[]string{"[4001:a1:1014::8a 4001:a1:1014::8b 4001:a1:1014::89]","[4001:a1:1014::8b 4001:a1:1014::89 4001:a1:1014::8a]",
		"[4001:a1:1014::89 4001:a1:1014::8a 4001:a1:1014::8b]"}
	var rr  []dns.RR
	var cname = test.CNAME("ipv6.cloud.example.com.	300	IN	CNAME		beta.cloud.example.com.")

	for i := 0; i < 10; i++ {
		m := newMid()
		m.AddResponseAnswer(test.AAAA("ipv6.cloud.example.com.		300	IN	AAAA			4001:a1:1014::89"))
		m.AddResponseAnswer(cname)
		m.AddResponseAnswer(test.AAAA("ipv6.cloud.example.com.		300	IN	AAAA			4001:a1:1014::8a"))
		m.AddResponseAnswer(test.AAAA("ipv6.cloud.example.com.		300	IN	AAAA			4001:a1:1014::8b"))
		// state, from the previous loop that arrived in the DNS query
		m.AddRequestOpt(rr...)
		var clientState = NewStateless().Shuffle(m.req, m.res)

		// save the new state for the next query
		rr = filterAandAAAA(clientState)

		if fmt.Sprintf("%v",getIPs(clientState)) != fmt.Sprintf("%v",aaaarotations[i %len(aaaarotations)]) {
			t.Errorf("%v: The stateless rotation is not working. Expecting %v but got %v.",i, aaaarotations[i %len(aaaarotations)], getIPs(clientState))
		}

		if len(clientState) != len(m.res.Answer) {
			t.Errorf("The stateless retrieved different number of records. Expected %v got %v", len(m.res.Answer), len(clientState))
		}

		if clientState[3].String() != cname.String() {
			t.Errorf("Expecting %s result but got %s", cname, clientState[3].String())
		}
	}
}

func TestRoundRobinStatelessShuffleEmpty(t *testing.T) {
	m := newMid()
	m.AddRequestOptRaw(`_rr_state={"ip":[10.240.0.2 10.240.0.3 10.240.0.4 10.240.0.1]}`)
	if len(NewStateless().Shuffle(m.req, m.res)) != 0 {
		t.Errorf("The stateless retrieved different number of records. Expected %v got %v", len(m.res.Answer), 0)
	}

	m = newMid()
	if len(NewStateless().Shuffle(m.req, m.res)) != 0 {
		t.Errorf("The stateless retrieved different number of records. Expected %v got %v", len(m.res.Answer), 0)
	}
}

func TestRoundRobinStatelessShuffleOne(t *testing.T) {
	a := test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.1")
	m := newMid()
	m.AddRequestOptRaw(`_rr_state={"ip":[10.240.0.2 10.240.0.3 10.240.0.4 10.240.0.1]}`)
	m.AddResponseAnswer(a)
	var answers = NewStateless().Shuffle(m.req, m.res)
	if len(answers) != 1 {
		t.Errorf("The stateless retrieved different number of records. Expected %v got %v", len(m.res.Answer), 0)
	}
	if answers[0].String() != a.String() {
		t.Errorf("The stateless shuffle doesnt work.  Expected %s got %s",answers[0],a)
	}

	m = newMid()
	m.AddResponseAnswer(a)
	answers = NewStateless().Shuffle(m.req, m.res)
	if len(answers) != 1 {
		t.Errorf("The stateless retrieved different number of records. Expected %v got %v", len(m.res.Answer), 0)
	}
	if answers[0].String() != a.String() {
		t.Errorf("The stateless shuffle doesnt work.  Expected %s got %s",answers[0],a)
	}
}