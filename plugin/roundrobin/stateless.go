package roundrobin

import (
	"encoding/json"
	"github.com/miekg/dns"
	"strings"
)

type stateless struct {
	IPs    []string `json:"ip"`
	// IPs converted into map
	requestA map[string]bool
	// response.Answers[] converted into map which are A or AAAA. IP is a key
	responseA map[string]dns.RR
	// contains all response records which are not A nor AAAA
	responseNoA []dns.RR
}

func newStateless(request *dns.Msg, response *dns.Msg) (s *stateless) {
	const identifierPrefix ="_rr_state="
	if request == nil || response == nil {
		return
	}
	s = &stateless{
		requestA: map[string]bool{},
		responseA: map[string]dns.RR{},
		responseNoA: []dns.RR{},
	}

	// extracting records from client request
	for _, e := range request.Extra {
		opt := e.(*dns.OPT)
		if  opt == nil {
			continue
		}
		for _, o := range opt.Option {
			e := o.(*dns.EDNS0_LOCAL)
			if e == nil {
				continue
			}
			data := string(o.(*dns.EDNS0_LOCAL).Data)
			if strings.HasPrefix(data, identifierPrefix){
				str := strings.Replace(data, identifierPrefix,"",-1)
				_ = json.Unmarshal([]byte(str),&s)
			}
		}
	}

	for _, a := range  response.Answer {
		switch a.Header().Rrtype {
		case dns.TypeA:
			s.responseA[a.(*dns.A).A.String()] = a
		case dns.TypeAAAA:
			s.responseA[a.(*dns.AAAA).AAAA.String()] = a
		default:
			s.responseNoA = append(s.responseNoA, a)
		}
	}
	for _, ip := range s.IPs {
		s.requestA[ip] = true
	}

	return
}

// updateState compare stateless records with response message records
// and cuts removed records or append new records to stateless
// updateState keeps records in the same order as they ar defined in tghe IPs field.
func (s *stateless) updateState()  *stateless {
	var newIPs []string

	// append only such IP which exist in response
	for _, ip := range s.IPs {
		if _, found := s.responseA[ip]; found {
			newIPs = append(newIPs, ip)
		}
	}

	// to the end of the IP list append new records which doesn't exist in request but exist in response.
	for ip := range s.responseA {
		if !s.requestA[ip] {
			newIPs = append(newIPs, ip)
		}
	}

	s.IPs = newIPs
	return s
}

// rotate performs a cyclic rotation of the IPs records
func (s *stateless) rotate() *stateless {
	var newIPs []string
	l := len(s.IPs)
	for i := range s.IPs {
		newIPs = append(newIPs, s.IPs[(i+1) % l])
	}
	s.IPs = newIPs
	return s
}

// getAnswers recreates Answer slice from original answers
func (s *stateless) getAnswers() []dns.RR{
	var shuffled []dns.RR
	for _, ip  := range s.IPs {
		shuffled = append(shuffled, s.responseA[ip])
	}
	return append(shuffled,s.responseNoA...)
}
