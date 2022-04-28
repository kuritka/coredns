package strategy

import (
	"fmt"
	"testing"
	"time"
)

func TestStatefulGCRemoveItem(t *testing.T) {
	tests := []struct {
		key       string
		question  string
		timestamp time.Time
		ips       []string
	}{
		{"10.20.30.40", "test.example.com.", time.Now().Add(time.Hour*-5), []string{"10.10.10.10"}},
		{"10.20.30.40", "alpha.example.com.", time.Now().Add(time.Hour*-5), []string{"10.10.10.10", "20.20.20.20"}},
		{"10.20.30.40", "beta.example.com.", time.Now().Add(time.Hour*-5), []string{}},
		{"11.11.11.11", "gc.test.com.", time.Now(), []string{"10.10.10.10"}},
	}

	s := make(map[key]map[question]state)
	for _, test := range tests {
		if _, ok := s[key(test.key)]; !ok {
			s[key(test.key)] = make(map[question]state)
		}
		s[key(test.key)][question(test.question)] = state{test.timestamp, test.ips}
	}

	newGarbageCollector(&s).collect()
	s["xxx"] = make(map[question]state)
	s["xxx"]["blah"] = state{ip: []string{"1"}}
	fmt.Println(s)

}
