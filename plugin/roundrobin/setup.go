package roundrobin

import (
	"fmt"
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
)

const (
	pluginName = "roundrobin"
)

var log = clog.NewWithPlugin(pluginName)

func init() { plugin.Register(pluginName, setup) }

func setup(c *caddy.Controller) error {
	strategy, err := parse(c)
	if err != nil {
		return plugin.Error(pluginName, err)
	}
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return New(next, strategy)
	})
	return nil
}

func parse(c *caddy.Controller) (shuffler,  error) {
	for c.Next() {
		args := c.RemainingArgs()
		if len(args) == 0 {
			return NewConsistent(), nil
		}
		switch args[0] {
		case strategyConsistent:
			return NewConsistent(), nil
		case strategyWeight:
			return nil, fmt.Errorf("not implemented %s", args[0])
		case strategyRandom:
			return NewRandom(), nil
		}
	}
	return nil, fmt.Errorf("unknown roundrobin type")
}
