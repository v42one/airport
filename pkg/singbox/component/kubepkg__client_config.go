package kubepkg

import (
	"github.com/sagernet/sing-box/option"
	"github.com/v42one/airport/pkg/runtime"
	"github.com/v42one/airport/pkg/singbox"
)

func (k *SingBox) ClientConfig() *option.Options {
	return runtime.Build(
		runtime.With(&singbox.ClientDNS{
			Proxy: singbox.Addr{
				Scheme: "https",
				Host:   "1.1.1.1",
				Path:   "dns-query",
			},
			Direct: singbox.Addr{
				Scheme: "h3",
				Host:   "dns.alidns.com",
				Path:   "dns-query",
			},
			Resolver: singbox.Addr{
				Host: "223.5.5.5",
			},
			Block: singbox.Addr{
				Scheme: "rcode",
				Host:   "refused",
			},
		}),
		runtime.With(&singbox.RemoteRuleSet{
			Name: "geosite-cn",
		}),
		runtime.With(&singbox.RemoteRuleSet{
			Name: "geoip-cn",
		}),
		runtime.With(&singbox.RemoteRuleSet{
			Name: "geosite-tiktok",
		}),
		runtime.With(&singbox.RemoteRuleSet{
			Name: "geosite-feishu",
		}),
		runtime.With(&singbox.DNSRule{
			Server: "dns-proxy",
			RuleSets: []singbox.RuleSet{
				"geosite-tiktok",
			},
		}),
		runtime.With(&singbox.DNSRule{
			Server: "dns-direct",
			RuleSets: []singbox.RuleSet{
				"geosite-feishu",
			},
		}),
		runtime.With(&singbox.DNSRule{
			Server: "dns-proxy",
			RuleSets: []singbox.RuleSet{
				"!geosite-cn",
			},
		}),
		runtime.With(&singbox.DNSRule{
			Server: "dns-direct",
			RuleSets: []singbox.RuleSet{
				"geosite-cn",
			},
		}),
		runtime.With(&singbox.SniffRouteRule{}),
		runtime.With(&singbox.HijackDNSRouteRule{}),
		singbox.WithRouteOptions(func(route *option.RouteOptions) {
			route.Rules = append(route.Rules, *runtime.Build(
				singbox.WithDefaultRule(
					func(r *option.DefaultRule) {
						r.IPIsPrivate = true
					},
					singbox.WithDefaultRouteActionOptions(func(r *option.RouteActionOptions) {
						r.Outbound = "direct"
					}),
				),
			))
			route.Final = "direct"
			route.AutoDetectInterface = true
		}),
		runtime.With(&singbox.RouteRule{
			Outbound: "proxy",
			DomainSuffixes: []string{
				"huggingface.co",
				"hf.co",
				"pub.dev",
				"pkg.dev",
			},
		}),
		runtime.With(&singbox.RouteRule{
			Outbound: "proxy",
			RuleSets: []singbox.RuleSet{
				"geosite-tiktok",
			},
		}),
		runtime.With(&singbox.RouteRule{
			Outbound: "proxy",
			RuleSets: []singbox.RuleSet{
				"!geosite-cn",
				"!geoip-cn",
			},
		}),
		runtime.With(&singbox.RouteRule{
			Outbound: "direct",
			RuleSets: []singbox.RuleSet{
				"geosite-cn",
				"geoip-cn",
			},
		}),
		singbox.WithRouteOptions(func(route *option.RouteOptions) {
			route.Rules = append(route.Rules, *runtime.Build(
				singbox.WithDefaultRule(
					func(r *option.DefaultRule) {
						r.IPCIDR = []string{}
					},
					singbox.WithDefaultRouteActionOptions(func(r *option.RouteActionOptions) {
						r.Outbound = "direct"
					}),
				),
			))
		}),
		runtime.With(&singbox.InboundTun{
			Address: "172.19.0.1/30",
		}),
		runtime.With(&singbox.OutboundVmess{
			Name: "proxy",
			Clusters: []singbox.ProxyServer{
				{
					Name:   singbox.Tag(k.ServerName),
					Server: k.ServerIP,
					Secret: k.VMess.Secret,
					Ports: []uint16{
						k.VMess.ListenPort,
						k.VMess.ListenBackupPort,
					},
				},
			},
		}),
		singbox.WithLogOptions(func(l *option.LogOptions) {
			l.Level = "info"
			l.Timestamp = true
		}),
		singbox.WithExperimentalOptions(func(experimental *option.ExperimentalOptions) {
			experimental.CacheFile = &option.CacheFileOptions{
				Enabled: true,
			}
		}),
	)
}
