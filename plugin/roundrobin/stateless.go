package roundrobin

import (
	"encoding/json"
	"github.com/miekg/dns"
	"strings"
)

type state struct {
	Family uint16   `json:"type"`
	IPs    []string `json:"ip"`
	// IPs converted into map
	requestA map[string]bool
	// response.Answers[] converted into map which are A or AAAA. IP is a key
	responseA map[string]dns.RR
	// contains all response records which are not A nor AAAA
	responseNoA []dns.RR
}

func newState(request *dns.Msg, response *dns.Msg) (s *state) {
	const identifierPrefix ="_rr_state="
	if request == nil || response == nil {
		return
	}
	s = &state{}

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
			s.responseA[a.(*dns.A).String()] = a
		case dns.TypeAAAA:
			s.responseA[a.(*dns.AAAA).String()] = a
		default:
			s.responseNoA = append(s.responseNoA, a)
		}
	}
	for _, ip := range s.IPs {
		s.requestA[ip] = true
	}

	return
}

// normalize compare state records with response message records
// and cuts removed records or append new records to state
// normalize keeps records in the same order as they ar defined in tghe IPs field.
func (s *state) normalize()  *state {
	var newIPs []string

	// append if request IP exist in response
	for _, ip := range s.IPs {
		if _, found := s.responseA[ip]; found {
			newIPs = append(newIPs, ip)
		}
	}

	// append to the end of the IP list new records which doesn't exist in request but exist in response.
	for ip := range s.responseA {
		if !s.requestA[ip] {
			newIPs = append(newIPs, ip)
		}
	}
	s.IPs = newIPs
	return s
}

// rotate performs a cyclic rotation of the IPs records
func (s *state) rotate() *state {
	var newIPs []string
	l := len(s.IPs)
	for i := range s.IPs {
		newIPs = append(newIPs, s.IPs[(i+1) % l])
	}
	s.IPs = newIPs
	return s
}

func (s *state) getAnswers() []dns.RR{
	var shuffled []dns.RR
	for _, ip  := range s.IPs {
		shuffled = append(shuffled, s.responseA[ip])
	}
	return append(shuffled,s.responseNoA...)
}
