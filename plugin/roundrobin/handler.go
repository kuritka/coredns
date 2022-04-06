package roundrobin

import (
	"context"
	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
)

type roundRobin struct {
	Next     plugin.Handler
	strategy shuffler
}

const (
	strategyWeight     = "weight"
	strategyConsistent = "consistent"
	strategyRandom     = "random"
)

func (rr roundRobin) ServeDNS(ctx context.Context, w dns.ResponseWriter, msg *dns.Msg) (int, error) {
	wrr, err := NewMessageWriter(w, msg, rr.strategy)
	if err != nil {
		return dns.RcodeServerFailure, err
	}
	return plugin.NextOrFailure(rr.Name(), rr.Next, ctx, wrr, msg)
}

func (rr roundRobin) Name() string {
	return pluginName
}
