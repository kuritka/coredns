package roundrobin

import (
	"fmt"
	"testing"

	"github.com/coredns/coredns/plugin/test"
)

var (
	avariations =[]string{"[10.240.0.1 10.240.0.2 10.240.0.3]","[10.240.0.1 10.240.0.3 10.240.0.2]",
	"[10.240.0.3 10.240.0.2 10.240.0.1]","[10.240.0.3 10.240.0.1 10.240.0.2]",
	"[10.240.0.2 10.240.0.1 10.240.0.3]","[10.240.0.2 10.240.0.3 10.240.0.1]"}

	aaaavariations = []string{"[4001:a1:1014::89 4001:a1:1014::8a 4001:a1:1014::8b]","[4001:a1:1014::89 4001:a1:1014::8b 4001:a1:1014::8a]",
	"[4001:a1:1014::8b 4001:a1:1014::8a 4001:a1:1014::89]","[4001:a1:1014::8b 4001:a1:1014::89 4001:a1:1014::8a]",
	"[4001:a1:1014::8a 4001:a1:1014::89 4001:a1:1014::8b]","[4001:a1:1014::8a 4001:a1:1014::8b 4001:a1:1014::89]"}
)

func TestRoundRobinRandomAMixed(t *testing.T){
	m := newMid()
	m.AddResponseAnswer(test.CNAME("alpha.cloud.example.com.	300	IN	CNAME		beta.cloud.example.com."))
	m.AddResponseAnswer(test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.1"))
	m.AddResponseAnswer(test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.2"))
	m.AddResponseAnswer(test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.3"))
	m.AddResponseAnswer(test.MX("alpha.cloud.example.com.			300	IN	MX		1	mxa-alpha.cloud.example.com."))

	for _, v := range avariations {
		var x int
		for x = 0; x < 100; x++ {
			result := NewRandom().Shuffle(m.req, m.res)
			if fmt.Sprintf("%v", getIPs(result)) == v {
				break
			}
		}
		if x == 100 {
			t.Errorf("The Random shuffle is not working as expected. %v didn't occure during %v attempts", v, x)
		}
	}
}

func TestRoundRobinRandomAOnly(t *testing.T){
	m := newMid()
	m.AddResponseAnswer(test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.1"))
	m.AddResponseAnswer(test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.2"))
	m.AddResponseAnswer(test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.3"))

	for _, v := range avariations {
		var x int
		for x = 0; x < 100; x++ {
			result := NewRandom().Shuffle(m.req, m.res)
			if fmt.Sprintf("%v", getIPs(result)) == v {
				break
			}
		}
		if x == 100 {
			t.Errorf("The Random shuffle is not working as expected. %v didn't occure during %v attempts", v, x)
		}
	}
}

func TestRoundRobinRandomAAAAOnly(t *testing.T){
	m := newMid()
	m.AddResponseAnswer(test.AAAA("ipv6.cloud.example.com.		300	IN	AAAA			4001:a1:1014::89"))
	m.AddResponseAnswer(test.AAAA("ipv6.cloud.example.com.		300	IN	AAAA			4001:a1:1014::8a"))
	m.AddResponseAnswer(test.AAAA("ipv6.cloud.example.com.		300	IN	AAAA			4001:a1:1014::8b"))

	for _, v := range aaaavariations {
		var x int
		for x = 0; x < 100; x++ {
			result := NewRandom().Shuffle(m.req, m.res)
			if fmt.Sprintf("%v", getIPs(result)) == v {
				break
			}
		}
		if x == 100 {
			t.Errorf("The Random shuffle is not working as expected. %v didn't occure during %v attempts", v, x)
		}
	}
}

func TestRoundRobinRandomEmptyAnswer(t *testing.T){
	m := newMid()
	result := NewRandom().Shuffle(m.req, m.res)
	if len(result) != 0 {
		t.Errorf("Expecting empty result but got %v", result)
	}
}

func TestRoundRobinRandomARecordsGoesFirst(t *testing.T){
	m := newMid()
	m.AddResponseAnswer(test.CNAME("alpha.cloud.example.com.	300	IN	CNAME		beta.cloud.example.com."))
	m.AddResponseAnswer(test.A("alpha.cloud.example.com.		300	IN	A			10.240.0.1"))
	result := NewRandom().Shuffle(m.req, m.res)
	if result[0].String() != m.res.Answer[1].String() {
		t.Errorf("Expecting %s result but got %s",result[0].String(), m.res.Answer[1].String())
	}
	if result[1].String() != m.res.Answer[0].String() {
		t.Errorf("Expecting %s result but got %s",result[1].String(), m.res.Answer[0].String())
	}
}

func TestRoundRobinRandomStableOrderForNonAandAAA(t *testing.T){
	m := newMid()
	m.AddResponseAnswer(test.CNAME("alpha.cloud.example.com.	300	IN	CNAME		beta.cloud.example.com."))
	m.AddResponseAnswer(test.MX("alpha.cloud.example.com.			300	IN	MX		1	mxa-alpha.cloud.example.com."))
	m.AddResponseAnswer(test.MX("alpha.cloud.example.com.			300	IN	MX		1	mx-beta.cloud.example.com."))
	result := NewRandom().Shuffle(m.req, m.res)
	if len(result) != 3 {
		t.Errorf("Expecting %v result but got %v", len(m.res.Answer), len(result))
	}
	for i, v := range result {
		if v.String() != m.res.Answer[i].String() {
			t.Errorf("Expecting %s result but got %s", v.String(), m.res.Answer[i].String())
		}
	}
}
