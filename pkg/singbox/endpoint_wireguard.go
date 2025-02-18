package singbox

import (
	"net/netip"

	"github.com/sagernet/sing-box/option"
	"github.com/v42one/airport/pkg/runtime"
)

type WireguardEndpoint struct {
	PrivateKey    string
	PeerPublicKey string
}

func (e WireguardEndpoint) ApplyTo(o *option.Options) {
	w := &option.WireGuardEndpointOptions{}

	w.Name = "wg0"
	w.MTU = 1280
	w.Address = []netip.Prefix{
		netip.MustParsePrefix("172.16.0.2/32"),
		netip.MustParsePrefix("2606:4700:110:80e4:9ad5:251:1e0c:959e/128"),
	}
	w.PrivateKey = e.PrivateKey
	w.ListenPort = 10000

	w.Peers = append(w.Peers, option.WireGuardPeer{
		Address:   "engage.cloudflareclient.com",
		Port:      2408,
		PublicKey: e.PeerPublicKey,
		Reserved:  []uint8{26, 8, 60},
		AllowedIPs: []netip.Prefix{
			netip.MustParsePrefix("0.0.0.0/0"),
		},
		PersistentKeepaliveInterval: 30,
	})

	o.Endpoints = append(o.Endpoints, option.Endpoint{
		Type:    "wireguard",
		Tag:     "wireguard-out",
		Options: w,
	})

	o.Outbounds = append(o.Outbounds, option.Outbound{
		Type: "direct",
		Tag:  "direct",
		Options: runtime.Build(func(o *option.DirectOutboundOptions) {
		}),
	})
}
