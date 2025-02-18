package singbox

import (
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
	runtime.Apply(o, WithDNSOptions(func(dnsOptions *option.DNSOptions) {
		block := option.DNSServerOptions{}
		block.Tag = "dns-block"
		block.Address = d.Block.String()

		resolver := option.DNSServerOptions{}
		resolver.Tag = "dns-resolver"
		resolver.Strategy = option.DomainStrategy(dns.DomainStrategyUseIPv4)
		resolver.Address = d.Resolver.String()
		resolver.Detour = "direct"

		direct := option.DNSServerOptions{}
		direct.Tag = "dns-direct"
		direct.Strategy = option.DomainStrategy(dns.DomainStrategyUseIPv4)
		direct.Address = d.Direct.String()
		direct.Detour = "direct"
		direct.AddressResolver = "dns-resolver"

		proxy := option.DNSServerOptions{}
		proxy.Tag = "dns-proxy"
		proxy.Strategy = option.DomainStrategy(dns.DomainStrategyUseIPv4)
		proxy.Address = d.Proxy.String()
		proxy.Detour = "proxy"
		proxy.AddressResolver = "dns-resolver"

		dnsOptions.Servers = append(dnsOptions.Servers,
			proxy,
			direct,
			resolver,
			block,
		)
	}))
}
