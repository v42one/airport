package singbox

import (
	"fmt"
	"slices"
	"time"

	"github.com/octohelm/exp/xiter"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common/json/badoption"
	"github.com/v42one/airport/pkg/runtime"
)

type ProxyServer struct {
	Name   Tag
	Server string
	Ports  []uint16
	Secret string
}

type OutboundVmess struct {
	Name     Tag
	Clusters []ProxyServer
}

func (x OutboundVmess) ApplyTo(opt *option.Options) {
	baseOutbounds := xiter.Map(xiter.Of(x.Clusters...), func(c ProxyServer) string {
		return string(c.Name)
	})

	auto := option.Outbound{
		Tag:  "auto",
		Type: "urltest",
		Options: runtime.Build(func(o *option.URLTestOutboundOptions) {
			o.Outbounds = slices.Collect(baseOutbounds)

			o.URL = "https://www.gstatic.com/generate_204"
			o.Interval = badoption.Duration(1 * time.Minute)
			o.Tolerance = 50
		}),
	}

	manual := option.Outbound{
		Tag:  "manual",
		Type: "selector",
		Options: runtime.Build(func(o *option.URLTestOutboundOptions) {
			o.Outbounds = slices.Collect(baseOutbounds)
		}),
	}

	proxy := option.Outbound{
		Tag:  string(x.Name),
		Type: "selector",
		Options: runtime.Build(func(o *option.SelectorOutboundOptions) {
			o.Outbounds = []string{
				"auto",
				"manual",
			}
		}),
	}

	opt.Outbounds = append(opt.Outbounds,
		option.Outbound{
			Tag:     "direct",
			Type:    "direct",
			Options: &option.DirectOutboundOptions{},
		},
		proxy,
		auto,
		manual,
	)

	for _, c := range x.Clusters {
		lb := runtime.Build(func(o *option.URLTestOutboundOptions) {
			o.URL = "https://www.gstatic.com/generate_204"
			o.Interval = badoption.Duration(1 * time.Minute)
			o.Tolerance = 50
		})

		opt.Outbounds = append(opt.Outbounds, option.Outbound{
			Tag:     string(c.Name),
			Type:    "urltest",
			Options: lb,
		})

		for _, port := range c.Ports {
			tag := fmt.Sprintf("%s-%d", c.Name, port)

			lb.Outbounds = append(lb.Outbounds, tag)

			opt.Outbounds = append(opt.Outbounds, option.Outbound{
				Tag:  tag,
				Type: "vmess",
				Options: runtime.Build(func(o *option.VMessOutboundOptions) {
					o.UUID = c.Secret
					o.Security = "auto"
					o.Server = c.Server
					o.ServerPort = port
				}),
			})
		}
	}

}
