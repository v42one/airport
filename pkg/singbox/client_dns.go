package singbox

import (
	"cmp"

	"github.com/sagernet/sing-box/option"
	dns "github.com/sagernet/sing-dns"

	"github.com/v42one/airport/pkg/runtime"
)

type ClientDNS struct {
	Proxy    Addr
	Direct   Addr
	Resolver Addr
	Block    Addr
}

func (d *ClientDNS) ApplyTo(o *option.Options) {
	runtime.Apply(o,
		WithDNSOptions(func(dnsOptions *option.DNSOptions) {
			local := option.DNSServerOptions{}
			local.Tag = "local"
			local.Type = "local"
			local.Options = make(map[string]interface{})

			resolver := option.DNSServerOptions{}
			resolver.Tag = "dns-resolver"
			resolver.Type = cmp.Or(d.Resolver.Scheme, "udp")
			resolver.Options = runtime.Build(func(o *option.RemoteDNSServerOptions) {
				o.Server = d.Resolver.Host
			})

			direct := option.DNSServerOptions{}
			direct.Tag = "dns-direct"
			direct.Type = cmp.Or(d.Direct.Scheme, "udp")
			direct.Options = runtime.Build(func(o *option.RemoteDNSServerOptions) {
				o.Server = d.Direct.Host

				o.DomainResolver = runtime.Build(func(dr *option.DomainResolveOptions) {
					dr.Server = resolver.Tag
				})
			})

			proxy := option.DNSServerOptions{}
			proxy.Tag = "dns-proxy"
			proxy.Type = cmp.Or(d.Proxy.Scheme, "udp")
			proxy.Options = runtime.Build(func(o *option.RemoteDNSServerOptions) {
				o.Server = d.Proxy.Host
				o.Detour = "proxy"

				o.DomainResolver = runtime.Build(func(dr *option.DomainResolveOptions) {
					dr.Server = resolver.Tag
				})
			})

			dnsOptions.Servers = append(dnsOptions.Servers,
				proxy,
				direct,
				resolver,
				local,
			)

			dnsOptions.Strategy = option.DomainStrategy(dns.DomainStrategyUseIPv4)
		}),
		WithRouteOptions(func(r *option.RouteOptions) {
			r.DefaultDomainResolver = runtime.Build(func(dr *option.DomainResolveOptions) {
				dr.Server = "local"
			})
		}),
	)
}
