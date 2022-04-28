package strategy

import (
	"fmt"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"time"
)

const (
	// defaultTTLSeconds defines the period after which the resource is removed
	defaultTTLSeconds = 600
	// garbageCollectionPeriodSeconds defines the period when garbage collection is triggered
	garbageCollectionPeriodSeconds = 5
	missingSubnet                  = "missing-subnet"
	emptySubnet                    = "empty-subnet"
)

// <clientIP>_<clientSubnet>
type key string

// one client could hit many domains
type question string

type state struct {
	timestamp time.Time
	ip        []string
}

type stateful struct {
	state map[key]map[question]state
}

func newStateful() *stateful {
	this := new(stateful)
	this.state = make(map[key]map[question]state)
	gc := newGarbageCollector(&this.state)
	go func() {
		for range time.Tick(time.Second * garbageCollectionPeriodSeconds) {
			gc.collect()
		}
	}()
	return this
}

func (s *stateful) update(req *request.Request, res *dns.Msg) (rr []dns.RR, err error) {
	if req == nil {
		err = fmt.Errorf("nil request")
		return
	}
	if res == nil {
		err = fmt.Errorf("nil response")
		return
	}
	if len(req.Req.Question) == 0 {
		err = fmt.Errorf("empty request question")
		return
	}
	return s.updateState(req, res)
}

func (s *stateful) updateState(req *request.Request, res *dns.Msg) (answer []dns.RR, err error) {
	q := question(req.Req.Question[0].Name)
	k := s.key(req)
	responseA, responseIPs, responseNoA := parseAnswerSection(res.Answer)
	s.refresh(k, q, responseA, responseIPs)
	for _, ip := range s.state[k][q].ip {
		answer = append(answer, responseA[ip])
	}
	return append(answer, responseNoA...), nil
}

func (s *stateful) key(req *request.Request) key {
	subnet := s.readSubnet(req.Req)
	return key(fmt.Sprintf("%s", subnet))
}

// readSubnet reads the option EDNS0_SUBNET which is usually filled by resolvers.
func (s *stateful) readSubnet(req *dns.Msg) string {
	for _, e := range req.Extra {
		opt := e.(*dns.OPT)
		if opt == nil {
			continue
		}
		for _, o := range opt.Option {
			x := o.(*dns.EDNS0_SUBNET)
			if x == nil {
				continue
			}
			if o.(*dns.EDNS0_SUBNET).Address.String() == "<nil>" {
				return emptySubnet
			}
			return o.(*dns.EDNS0_SUBNET).Address.String()
		}
	}
	return missingSubnet
}

func (s *stateful) refresh(k key, q question, responseA map[string]dns.RR, responseIPs []string) {
	var st state
	if _, found := s.state[k]; !found {
		s.state[k] = make(map[question]state)
	}
	if _, found := s.state[k][q]; !found {
		s.state[k][q] = state{
			ip:        []string{},
			timestamp: time.Now(),
		}
	}
	st = s.state[k][q]
	st.updateState(responseA, responseIPs)
	st.ip = rotate(st.ip)
	s.state[k][q] = st
}

func (s *state) updateState(responseA map[string]dns.RR, responseIPs []string) {
	var newIPs []string
	currentA := ipsToSet(s.ip)

	// append only such IP which exist in response
	for _, ip := range s.ip {
		if _, found := responseA[ip]; found {
			newIPs = append(newIPs, ip)
		}
	}

	// to the end of the IP list append new records which doesn't exist in request but exist in response.
	for _, ip := range responseIPs {
		if !currentA[ip] {
			newIPs = append(newIPs, ip)
		}
	}
	s.ip = newIPs
}
