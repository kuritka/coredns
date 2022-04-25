package roundrobin

import (
	"context"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/roundrobin/internal/strategy"
	"github.com/miekg/dns"
)

type RoundRobin struct {
	Next     plugin.Handler
	strategy strategy.Shuffler
}

const (
	strategyWeight    = "weight"
	strategyStateless = "stateless"
	strategyRandom    = "random"
	strategyStateful  = "stateful"
)

func New(next plugin.Handler, strategy strategy.Shuffler) *RoundRobin {
	return &RoundRobin{
		Next:     next,
		strategy: strategy,
	}
}

func (rr *RoundRobin) ServeDNS(ctx context.Context, w dns.ResponseWriter, msg *dns.Msg) (int, error) {
	wrr, err := NewMessageWriter(w, msg, rr.strategy)
	if err != nil {
		return dns.RcodeServerFailure, err
	}
	return plugin.NextOrFailure(rr.Name(), rr.Next, ctx, wrr, msg)
}

func (rr *RoundRobin) Name() string {
	return pluginName
}
