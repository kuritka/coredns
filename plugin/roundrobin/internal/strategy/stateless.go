package strategy

import (
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

type Stateless struct {
}

func NewStateless() *Stateless {
	return &Stateless{}
}

func (r *Stateless) Shuffle(req request.Request, msg *dns.Msg) ([]dns.RR, error) {
	return newStateless(req.Req, msg).updateState().rotate().getAnswers()
}
