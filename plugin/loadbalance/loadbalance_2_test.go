package loadbalance

import (
	"fmt"
	"github.com/coredns/coredns/plugin/test"
	"github.com/miekg/dns"
	"math"
	"math/rand"
	"sort"
	"strings"
	"testing"
)

func pickIndexes2(cdf []float32, l int) (p []int) {
	bucket := func(cdf []float32) int{
		r := rand.Float32()
		bucket := 0
		for r > cdf[bucket] {
			bucket++
		}
		return bucket
	}
	x := l
	from := 0
	to := 0
	for i := 0; i < l; i++ {
		from = to
		to = from + x
		ix := bucket(cdf[from:to])
		x--
		p = append(p,ix)
	}
	return p
}

func getCDF2(distribution []float32) (cdf []float32){
	l := len(distribution)
	sum := (l*(l+1))/2
	for i:= 0; i < sum; i++ {
		cdf = append(cdf,0.0)
	}
	x := 0
	// cdf looks like this: [[p1,p1+p2,p2+p3,p3+p4][p2,p2+p3,p3+p4][p3,p3+p4],[p4]]
	for i := 0; i < l; i++ {
		cdf[x] = distribution[i]
		var sum = distribution[i]
		x++
		for j := 1; j < l-i; j++ {
			cdf[x] = float32(math.Round(float64( (cdf[x-1] + distribution[j+i]) *100) )/100)
			sum += cdf[x]
			x++
		}

	}
	return
}

func weightRoundRobinShuffle(records []dns.RR) {
	switch l := len(records); l {
	case 0, 1:
		break
	default:
		for j := 0; j < l; j++ {
			p := j + (int(dns.Id()) % (l - j))
			if j == p {
				continue
			}
			records[j], records[p] = records[p], records[j]
		}
	}
}



func TestRRNoWeight(t *testing.T){
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