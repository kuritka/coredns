package roundrobin

import (
	"encoding/json"
	"fmt"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"strings"
)

type Consistent struct {
}

type State struct {
	Family uint16   `json:"type"`
	IPs    []string `json:"ip"`
}

func (s State) IsProvided() bool {
	return !(s.Family == 0 && len(s.IPs) == 0)
}

func NewConsistent() *Consistent {
	return &Consistent{}
}

func (r *Consistent) Shuffle(req request.Request, msg *dns.Msg) []dns.RR{
	clientState := r.extractClientState(req.Req)
	fmt.Println(clientState)
	return msg.Answer
}

// extracts state from client in case it was sent
func (r *Consistent) extractClientState(msg *dns.Msg) (state State) {
	state = State{}
	if msg == nil {
		return
	}
	for _, e := range msg.Extra {
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
			if strings.HasPrefix(data, "_rr_state="){
				str := strings.Replace(data, "_rr_state=","",-1)
				json.Unmarshal([]byte(str),&state)
			}
		}
	}
	return
}
