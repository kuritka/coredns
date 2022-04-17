package roundrobin

import (
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

type MessageWriter struct {
	dns.ResponseWriter
	strategy shuffler
	state request.Request
}



func NewMessageWriter(w dns.ResponseWriter, msg *dns.Msg, strategy shuffler) (*MessageWriter, error) {
	return &MessageWriter{
		state: request.Request{W: w, Req: msg},
		ResponseWriter: w,
		strategy: strategy,
	}, nil
}

// WriteMsg implements the dns.ResponseWriter interface.
func (r *MessageWriter) WriteMsg(msg *dns.Msg) error {
	if msg.Rcode != dns.RcodeSuccess {
		return r.ResponseWriter.WriteMsg(msg)
	}
	q := msg.Question[0]

	if q.Qtype == dns.TypeAXFR || q.Qtype == dns.TypeIXFR {
		return r.ResponseWriter.WriteMsg(msg)
	}

	msg.Answer = r.strategy.Shuffle(r.state, msg)

	return r.ResponseWriter.WriteMsg(msg)
}

// Write implements the dns.ResponseWriter interface.
func (r *MessageWriter) Write(buf []byte) (int, error) {
	// Should we pack and unpack here to fiddle with the packet... Not likely.
	log.Warning("RoundRobin called with Write: not shuffling records")
	n, err := r.ResponseWriter.Write(buf)
	return n, err
}
