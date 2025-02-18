package singbox

import (
	"net/netip"

	"github.com/sagernet/sing-box/option"
)

type InboundTun struct {
	Address string
}

func (d InboundTun) ApplyTo(o *option.Options) {
	tunInbound := option.TunInboundOptions{}
	tunInbound.Address = append(tunInbound.Address, netip.MustParsePrefix(d.Address))
	tunInbound.Stack = "system"
	tunInbound.AutoRoute = true
	tunInbound.StrictRoute = true

	o.Inbounds = append(o.Inbounds, option.Inbound{
		Type:    "tun",
		Tag:     "in",
		Options: tunInbound,
	})
}
