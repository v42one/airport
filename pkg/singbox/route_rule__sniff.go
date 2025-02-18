package singbox

import (
	"github.com/sagernet/sing-box/option"
	"github.com/v42one/airport/pkg/runtime"
)

type SniffRouteRule struct {
}

func (SniffRouteRule) ApplyTo(o *option.Options) {
	runtime.Apply(o, WithRouteOptions(func(route *option.RouteOptions) {
		route.Rules = append(route.Rules, *runtime.Build(
			WithDefaultRule(func(r *option.DefaultRule) {
				r.Action = "sniff"
			}),
		))
	}))
}
