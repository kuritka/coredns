package roundrobin

import (
	"fmt"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"time"
)

const (
	// defaultTTLSeconds defines when state object is garbage collected
	defaultTTLSeconds = 600
	// garbageCollectionPeriodSeconds defines the period when garbage collection is triggered
	garbageCollectionPeriodSeconds = 1200
)

// <clientIP>_<clientSubnet>
type key string
// one client could hit many domains
type question string

type state struct {
	timestamp 	time.Time
	a []string
	aaaa []string
}

type stateful struct {
	state map[key]map[question]state
}

func newStateful() *stateful {
	return &stateful{
		state: make(map[key]map[question]state),
	}
}

func (s *stateful) handle(req *request.Request) error {
	if req == nil {
		return fmt.Errorf("nil request")
	}
	s.updateState(req)
	return nil
}

func (s *stateful) updateState(req *request.Request) {
	k := s.key(req)
	s.state[k] = s.newState(req.Req)

}

func (s *stateful) key(req *request.Request) key {
	subnet := s.readSubnet(req.Req)
	return key(fmt.Sprintf("%s_%s", req.IP(), subnet))
}

// readSubnet reads the option EDNS0_SUBNET which is usually filled by resolvers.
func (s *stateful) readSubnet(req *dns.Msg) string {
	for _, e := range req.Extra {
		opt := e.(*dns.OPT)
		if  opt == nil {
			continue
		}
		for _, o := range opt.Option {
			x := o.(*dns.EDNS0_SUBNET)
			if x == nil {
				continue
			}
			return o.(*dns.EDNS0_SUBNET).Address.String()
		}
	}
	return ""
}

func (s *stateful) newState(req *dns.Msg) (m map[question]state) {
	//todo:
	m = make(map[question]state)
	m["0"] = state{
		timestamp: time.Now(),
	}
	return
}
