package roundrobin

import (
	"fmt"
	"testing"

	"github.com/coredns/coredns/plugin/test"
	"github.com/miekg/dns"
)

func TestRoundRobinStatefulShuffleA(t *testing.T) {
	var arotations = []string{"[10.240.0.2 10.240.0.3 10.240.0.4 10.240.0.1]", "[10.240.0.3 10.240.0.4 10.240.0.1 10.240.0.2]",
		"[10.240.0.4 10.240.0.1 10.240.0.2 10.240.0.3]", "[10.240.0.1 10.240.0.2 10.240.0.3 10.240.0.4]"}
	cname := test.CNAME("alpha.cloud.example.com.	300	IN	CNAME		beta.cloud.example.com.")
	mx := test.MX("alpha.cloud.example.com.			300	IN	MX		1	mxa-alpha.cloud.example.com.")
	var s = NewStateful()
	for i := 0; i < 10; i++ {
		m := newMid()
		m.SetQuestion("alpha.cloud.example.com.", dns.TypeA)
		m.AddResponseAnswer(cname)
		m.AddResponseAnswer(test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.1"))
		m.AddResponseAnswer(test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.2"))
		m.AddResponseAnswer(test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.3"))
		m.AddResponseAnswer(mx)
		m.AddResponseAnswer(test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.4"))

		clientState := s.Shuffle(m.req, m.res)

		if fmt.Sprintf("%v", getIPs(clientState)) != fmt.Sprintf("%v", arotations[i%len(arotations)]) {
			t.Errorf("%v: The stateful rotation is not working. Expecting %v but got %v.", i, arotations[i%len(arotations)], getIPs(clientState))
		}

		if len(clientState) != len(m.res.Answer) {
			t.Errorf("The stateful retrieved different number of records. Expected %v got %v", len(m.res.Answer), len(clientState))
		}

		if clientState[len(clientState)-1] != mx {
			t.Errorf("Expecting %s result but got %s", mx, clientState[len(clientState)-1].String())
		}

		if clientState[len(clientState)-2] != cname {
			t.Errorf("Expecting %s result but got %s", cname, clientState[len(clientState)-2].String())
		}
	}
}

func TestRoundRobinStatefulShuffleAAAA(t *testing.T) {
	var aaaarotations = []string{"[4001:a1:1014::8a 4001:a1:1014::8b 4001:a1:1014::89]", "[4001:a1:1014::8b 4001:a1:1014::89 4001:a1:1014::8a]",
		"[4001:a1:1014::89 4001:a1:1014::8a 4001:a1:1014::8b]"}
	var s = NewStateful()
	for i := 0; i < 10; i++ {
		cname := test.CNAME("alpha.cloud.example.com.	300	IN	CNAME		beta.cloud.example.com.")
		m := newMid()
		m.SetQuestion("alpha.cloud.example.com.", dns.TypeAAAA)
		m.AddResponseAnswer(test.AAAA("ipv6.cloud.example.com.		300	IN	AAAA			4001:a1:1014::89"))
		m.AddResponseAnswer(cname)
		m.AddResponseAnswer(test.AAAA("ipv6.cloud.example.com.		300	IN	AAAA			4001:a1:1014::8a"))
		m.AddResponseAnswer(test.AAAA("ipv6.cloud.example.com.		300	IN	AAAA			4001:a1:1014::8b"))

		clientState := s.Shuffle(m.req, m.res)

		if fmt.Sprintf("%v", getIPs(clientState)) != fmt.Sprintf("%v", aaaarotations[i%len(aaaarotations)]) {
			t.Errorf("%v: The stateful rotation is not working. Expecting %v but got %v.", i, aaaarotations[i%len(aaaarotations)], getIPs(clientState))
		}

		if len(clientState) != len(m.res.Answer) {
			t.Errorf("The stateful retrieved different number of records. Expected %v got %v", len(m.res.Answer), len(clientState))
		}

		if clientState[len(clientState)-1].String() != cname.String() {
			t.Errorf("Expecting %s result but got %s", cname, clientState[len(clientState)-1].String())
		}
	}
}

func TestRoundRobinStatefulShuffleEmpty(t *testing.T) {
	m := newMid()
	m.SetQuestion("alpha.cloud.example.com.", dns.TypeA)
	if len(NewStateful().Shuffle(m.req, m.res)) != 0 {
		t.Errorf("The stateful retrieved different number of records. Expected %v got %v", len(m.res.Answer), 0)
	}
}

func TestRoundRobinStatefulShuffleOne(t *testing.T) {
	a := test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.1")
	m := newMid()
	m.SetQuestion("alpha.cloud.example.com.", dns.TypeA)
	m.AddResponseAnswer(a)
	answers := NewStateful().Shuffle(m.req, m.res)
	if len(answers) != 1 {
		t.Errorf("The stateful retrieved different number of records. Expected %v got %v", len(m.res.Answer), 0)
	}
	if answers[0].String() != a.String() {
		t.Errorf("The stateful shuffle doesnt work. Expected %s got %s", answers[0], a)
	}
}

func TestRoundRobinStatefulNoQuestion(t *testing.T) {
	m := newMid()
	if len(NewStateful().Shuffle(m.req, m.res)) != 0 {
		t.Errorf("The stateful retrieved different number of records. Expected %v got %v", len(m.res.Answer), 0)
	}
}

func TestRoundRobinStatefulState(t *testing.T) {
	s := NewStateful()
	tests := []struct {
		question    string
		from        string
		expectedKey string
		setSubnet   bool
		rr          []dns.RR
	}{
		{"alpha.cloud.example.com.", "200.10.0.0", "200.10.0.0", true,
			[]dns.RR{
				test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.1"),
				test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.2")}},
		{"alpha.cloud.example.com.", "101.203.0.0", "101.203.0.0", true,
			[]dns.RR{
				test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.1"),
				test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.2")}},
		{"beta.cloud.example.com.", "101.203.0.0", "101.203.0.0", true,
			[]dns.RR{
				test.A("beta.cloud.example.com.		300	IN	A			20.100.0.1")}},
		{"beta.cloud.example.com.", "102.203.0.0", "102.203.0.0", true,
			[]dns.RR{
				test.A("beta.cloud.example.com.		300	IN	A			20.100.0.1")}},
		{"ipv6-subnet.cloud.example.com.", "4001:a1:1014::8a", "4001:a1:1014::8a", true,
			[]dns.RR{
				test.A("ipv6.cloud.example.com.		300	IN	A			10.240.0.1"),
				test.A("ipv6.cloud.example.com.		300	IN	A			10.240.0.2"),
				test.A("ipv6.cloud.example.com.		300	IN	A			10.240.0.3")}},
		{"no-records.cloud.example.com.", "200.10.0.0", "200.10.0.0", true, []dns.RR{}},
		{"empty-subnet-and-empty-response.cloud.com.", "", emptySubnet, true, []dns.RR{}},
		{"empty-subnet-with-a-records.cloud.com.", "", emptySubnet, true, []dns.RR{
			test.A("empty-subnet-with-a-records.cloud.com.		300	IN	A			10.240.0.1"),
			test.A("empty-subnet-with-a-records.cloud.com.		300	IN	A			10.240.0.2"),
		}},
		{"missing-subnet-and-empty-response.cloud.com.", "", missingSubnet, false, []dns.RR{}},
		{"missing-subnet-with-a-records.cloud.com.", "", missingSubnet, false, []dns.RR{
			test.A("missing-subnet-with-a-records.cloud.com.		300	IN	A			10.240.0.1"),
			test.A("missing-subnet-with-a-records.cloud.com.		300	IN	A			10.240.0.2"),
		}},
		{"only-cname.cloud.example.com.", "200.10.0.0", "200.10.0.0", true, []dns.RR{
			test.CNAME("alpha.cloud.example.com.	300	IN	CNAME		beta.cloud.example.com."),
		}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Request %s from %s", test.question, test.from), func(t *testing.T) {

			// arrange
			m := newMid()
			m.SetQuestion(test.question, dns.TypeA)
			if test.setSubnet {
				m.SetSubnet(test.from)
			}
			for _, a := range test.rr {
				m.AddResponseAnswer(a)
			}

			//act
			_ = s.Shuffle(m.req, m.res)

			// assert
			ipMap := ipsToSet(getIPs(test.rr))
			if len(s.state.state[key(test.expectedKey)][question(test.question)].ip) != len(getIPs(test.rr)) {
				t.Errorf("the number of records in the test (%v) and the state (%v) do not match.",
					len(test.rr), len(s.state.state[key(test.from)][question(test.question)].ip))
			}
			for _, ip := range s.state.state[key(test.expectedKey)][question(test.question)].ip {
				if !ipMap[ip] {
					t.Errorf("Can't find %s for state[%s][%s] ", ip, test.from, test.question)
				}
			}
		})
	}
}
