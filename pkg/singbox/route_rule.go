package singbox

import (
	"github.com/sagernet/sing-box/option"
	"github.com/v42one/airport/pkg/runtime"
)

type RouteRule struct {
	Outbound Tag

	RuleSets       []RuleSet
	DomainSuffixes []string
}

func (x RouteRule) ApplyTo(o *option.Options) {
	if x.Outbound == "" {
		return
	}

	runtime.Apply(o, WithRouteOptions(func(r *option.RouteOptions) {
		if and := x.RuleSets; len(and) > 0 {
			r.Rules = append(r.Rules, *runtime.Build(
				WithLogicalRule(
					func(r *option.LogicalRule) {
						r.Mode = "and"

						for _, rs := range and {
							r.Rules = append(r.Rules, *runtime.Build(
								WithDefaultRule(func(d *option.DefaultRule) {
									d.RuleSet = []string{string(rs.Tag())}
									d.Invert = rs.Invert()
								}),
							))
						}
					},
					WithLogicalRouteActionOptions(func(action *option.RouteActionOptions) {
						action.Outbound = string(x.Outbound)
					}),
				),
			))

			return
		}

		if domainSuffix := x.DomainSuffixes; len(domainSuffix) > 0 {
			r.Rules = append(r.Rules, *runtime.Build(
				WithDefaultRule(
					func(r *option.DefaultRule) {
						r.DomainSuffix = domainSuffix
					},
					WithDefaultRouteActionOptions(func(r *option.RouteActionOptions) {
						r.Outbound = string(x.Outbound)
					}),
				),
			))
			return
		}
	}))
}
