package singbox

import (
	"net/netip"

	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common/json/badoption"
)

type InboundVMess struct {
	Secret           string
	ListenPort       uint16
	ListenBackupPort uint16
}

func (d InboundVMess) ApplyTo(o *option.Options) {
	vmessIn := option.VMessInboundOptions{
		Listen:     new(badoption.Addr(netip.MustParseAddr("::"))),
		ListenPort: d.ListenPort,
		Users: []option.VMessUser{{
			UUID: d.Secret,
		}}}

	o.Inbounds = append(o.Inbounds, option.Inbound{
		Tag:     "vmess-in",
		Type:    "vmess",
		Options: vmessIn,
	})

	vmessInBackup := option.VMessInboundOptions{
		Listen:     new(badoption.Addr(netip.MustParseAddr("::"))),
		ListenPort: d.ListenBackupPort,
		Users: []option.VMessUser{{
			UUID: d.Secret,
		}}}

	o.Inbounds = append(o.Inbounds, option.Inbound{
		Tag: "vmess-in-backup",

		Type:    "vmess",
		Options: vmessInBackup,
	})
}
