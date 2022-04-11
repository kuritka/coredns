package roundrobin

import (

	"fmt"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"

)

type Consistent struct {
}



func NewConsistent() *Consistent {
	return &Consistent{}
}

func (r *Consistent) Shuffle(req request.Request, msg *dns.Msg) []dns.RR{
	clientState := newState(req.Req, msg).normalize().rotate()
	fmt.Println(clientState)
	return msg.Answer
}

