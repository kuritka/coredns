package roundrobin

import (
	"fmt"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

type Subnet struct {
	 state *stateful
}

func NewSubnet() *Subnet {
	return &Subnet{
		state: newStateful(),
	}
}

func (s *Subnet) Shuffle(req request.Request, msg *dns.Msg) []dns.RR {
	s.state.process(req.Req,msg)
	return nil
}


type stateful struct {
}

func newStateful() (state *stateful){
	return &stateful{}
}

func (s *stateful) process(request *dns.Msg, response *dns.Msg) {
	if request == nil || response == nil {
		return
	}
	// extracting subnet from client request
	for _, e := range request.Extra {
		opt := e.(*dns.OPT)
		if  opt == nil {
			continue
		}
		for _, o := range opt.Option {
			e := o.(*dns.EDNS0_SUBNET)
			if e == nil {
				continue
			}
			ip := o.(*dns.EDNS0_SUBNET).Address.String()
			fmt.Println(ip)
		}
	}
}


