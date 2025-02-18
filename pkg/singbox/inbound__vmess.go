package singbox

import (
	"github.com/octohelm/x/ptr"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common/json/badoption"
	"net/netip"
)

type InboundVMess struct {
	Secret           string
	ListenPort       uint16
	ListenBackupPort uint16
}

func (d InboundVMess) ApplyTo(o *option.Options) {
	vmessIn := option.VMessInboundOptions{}

	vmessIn.Listen = ptr.Ptr(badoption.Addr(netip.MustParseAddr("::")))
	vmessIn.ListenPort = d.ListenPort

	vmessIn.Users = []option.VMessUser{{
		UUID: d.Secret,
	}}

	o.Inbounds = append(o.Inbounds, option.Inbound{
		Tag: "vmess-in",

		Type:    "vmess",
		Options: vmessIn,
	})

	vmessInBackup := option.VMessInboundOptions{}
	vmessInBackup.Listen = ptr.Ptr(badoption.Addr(netip.MustParseAddr("::")))
	vmessInBackup.ListenPort = d.ListenBackupPort
	vmessInBackup.Users = []option.VMessUser{{
		UUID: d.Secret,
	}}

	o.Inbounds = append(o.Inbounds, option.Inbound{
		Tag: "vmess-in-backup",

		Type:    "vmess",
		Options: vmessInBackup,
	})
}
