package singbox

import (
	"github.com/sagernet/sing-box/option"
	"github.com/v42one/airport/pkg/runtime"
)

type DNSRule struct {
	Server         Tag
	DomainSuffixes []string
	RuleSets       []RuleSet
}

func (x DNSRule) ApplyTo(o *option.Options) {
	if x.Server == "" {
		return
	}

	runtime.Apply(o, WithDNSOptions(func(target *option.DNSOptions) {
		if and := x.RuleSets; len(and) > 0 {
			target.Rules = append(target.Rules, *runtime.Build(
				WithLogicalDNSRule(
					func(r *option.LogicalDNSRule) {
						r.Mode = "and"

						for _, rs := range and {
							r.Rules = append(r.Rules, *runtime.Build(
								WithDefaultDNSRule(func(r *option.DefaultDNSRule) {
									r.RuleSet = []string{string(rs.Tag())}
									r.Invert = rs.Invert()
								})),
							)
						}
					},
					WithLogicalDNSRouteActionOptions(func(target *option.DNSRouteActionOptions) {
						target.Server = string(x.Server)
					}),
				),
			))

			return
		}

		if domainSuffix := x.DomainSuffixes; len(domainSuffix) > 0 {
			target.Rules = append(target.Rules, *runtime.Build(
				WithDefaultDNSRule(
					func(r *option.DefaultDNSRule) {
						r.DomainSuffix = domainSuffix
					},
					WithDefaultDNSRouteActionOptions(func(target *option.DNSRouteActionOptions) {
						target.Server = string(x.Server)
					}),
				),
			),
			)
			return
		}
	}))
}
