package roundrobin

import (
	"fmt"
	"testing"

	"github.com/coredns/coredns/plugin/test"
	"github.com/miekg/dns"
)

func TestRoundRobinStatefulShuffleA(t *testing.T) {
	var arotations =[]string{"[10.240.0.2 10.240.0.3 10.240.0.4 10.240.0.1]","[10.240.0.3 10.240.0.4 10.240.0.1 10.240.0.2]",
		"[10.240.0.4 10.240.0.1 10.240.0.2 10.240.0.3]", "[10.240.0.1 10.240.0.2 10.240.0.3 10.240.0.4]"}
	var s = NewStateful()
	for i := 0; i < 10; i++ {
		cname := test.CNAME("alpha.cloud.example.com.	300	IN	CNAME		beta.cloud.example.com.")
		m := newMid()
		m.SetQuestion("alpha.cloud.example.com.", dns.TypeA)
		m.AddResponseAnswer(cname)
		m.AddResponseAnswer(test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.1"))
		m.AddResponseAnswer(test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.2"))
		m.AddResponseAnswer(test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.3"))
		m.AddResponseAnswer(test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.4"))

		clientState := s.Shuffle(m.req, m.res)

		if fmt.Sprintf("%v",getIPs(clientState)) != fmt.Sprintf("%v",arotations[i %len(arotations)]) {
			t.Errorf("%v: The stateful rotation is not working. Expecting %v but got %v.",i, arotations[i %len(arotations)], getIPs(clientState))
		}

		if len(clientState) != len(m.res.Answer) {
			t.Errorf("The stateful retrieved different number of records. Expected %v got %v", len(m.res.Answer), len(clientState))
		}
	}
}

func TestRoundRobinStatefulShuffleAAAA(t *testing.T) {
	var aaaarotations =[]string{"[4001:a1:1014::8a 4001:a1:1014::8b 4001:a1:1014::89]","[4001:a1:1014::8b 4001:a1:1014::89 4001:a1:1014::8a]",
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

		if fmt.Sprintf("%v",getIPs(clientState)) != fmt.Sprintf("%v",aaaarotations[i %len(aaaarotations)]) {
			t.Errorf("%v: The stateful rotation is not working. Expecting %v but got %v.",i, aaaarotations[i %len(aaaarotations)], getIPs(clientState))
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
		t.Errorf("The stateful shuffle doesnt work. Expected %s got %s",answers[0],a)
	}
}

func TestRoundRobinStatefulNoQuestion(t *testing.T){
	m := newMid()
	if len(NewStateful().Shuffle(m.req, m.res)) != 0 {
		t.Errorf("The stateful retrieved different number of records. Expected %v got %v", len(m.res.Answer), 0)
	}
}