package singbox

import (
	"github.com/sagernet/sing-box/option"
	"github.com/v42one/airport/pkg/runtime"
)

type HijackDNSRouteRule struct {
}

func (HijackDNSRouteRule) ApplyTo(o *option.Options) {
	runtime.Apply(o, WithRouteOptions(func(route *option.RouteOptions) {
		route.Rules = append(route.Rules, *runtime.Build(
			WithDefaultRule(func(r *option.DefaultRule) {
				r.Protocol = []string{"dns"}
				r.Action = "hijack-dns"
			}),
		))
	}))
}
