package kubepkg

import (
	"github.com/sagernet/sing-box/option"

	"github.com/v42one/airport/pkg/runtime"
	"github.com/v42one/airport/pkg/singbox"
)

func (k *SingBox) ServerConfig() *option.Options {
	return runtime.Build(
		runtime.With(&singbox.SniffRouteRule{}),
		runtime.With(&singbox.RouteRule{
			Outbound: "wireguard-out",
			RuleSets: []singbox.RuleSet{
				"geosite-openai",
			},
		}),
		singbox.WithRouteOptions(func(route *option.RouteOptions) {
			route.Final = "direct"
			route.AutoDetectInterface = true
		}),
		runtime.With(k.Wireguard),
		runtime.With(k.VMess),
		runtime.With(&singbox.RemoteRuleSet{
			Name: "geosite-openai",
		}),
		singbox.WithLogOptions(func(l *option.LogOptions) {
			l.Level = "info"
			l.Timestamp = true
		}),
	)
}
