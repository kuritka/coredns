package loadbalance

import (
	"fmt"
	"github.com/coredns/coredns/plugin/test"
	"github.com/miekg/dns"
	"sort"
	"strings"
	"testing"
)

func TestRRWeight(t *testing.T){
	inputs := []dns.RR{
		test.A("one.example.com.		300	IN	A			10.0.0.1"),
		test.A("two.example.com.		300	IN	A			10.0.0.2"),
		test.A("three.example.com.		300	IN	A			10.0.0.3"),
		test.A("four.example.com.		300	IN	A			10.0.0.4"),
	}

	variationsQuantity := map[string]int{}

	for i := 0; i < 1000000; i++ {
		// function roundRobinShuffle is shuffler in coreDNS loadbalance plugin
		roundRobinShuffle(inputs)
		// stringify test.A records to string format like: "10.0.0.1,10.0.0.2,10.0.0.3,10.0.0.4"
		key := getRRkeyFromA(inputs)
		// adds key to map and increase counter
		variationsQuantity[key]++
	}

	print(variationsQuantity)
}

// stringify test.A records to format like: "10.0.0.1,10.0.0.2,10.0.0.3,10.0.0.4"
func getRRkeyFromA(rr []dns.RR) (key string){
	var ips []string
	for _,v :=  range rr {
		ips = append(ips, v.(*dns.A).A.String())
	}
	return strings.Join(ips,",")
}

func print(m map[string]int){
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for i,v := range keys {
		fmt.Printf("%v. [%s]: %v\n",i+1,v,m[v])
	}
}