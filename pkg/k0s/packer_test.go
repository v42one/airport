package k0s

import (
	"cmp"
	"context"
	"os"
	"testing"

	kubepkgv1alpha1 "github.com/octohelm/kubepkgspec/pkg/apis/kubepkg/v1alpha1"
	"github.com/octohelm/unifs/pkg/filesystem/local"

	"github.com/v42one/airport/pkg/runtime"
	"github.com/v42one/airport/pkg/singbox"
	singboxcomponent "github.com/v42one/airport/pkg/singbox/component"
)

func TestPacker(t *testing.T) {
	root := local.NewFS("../../build")

	p := &Cluster{
		Name:         "proxy-sg",
		K0sVersion:   "1.34.1+k0s.1",
		RemoteServer: cmp.Or(os.Getenv("VMESS_REMOTE_SERVER"), "127.0.0.1"),

		Components: []*kubepkgv1alpha1.KubePkg{
			runtime.Build(
				runtime.With(&singboxcomponent.SingBox{
					Version:    "1.12.12",
					ServerName: "sg",
					ServerIP:   cmp.Or(os.Getenv("VMESS_REMOTE_SERVER"), "127.0.0.1"),
					VMess: &singbox.InboundVMess{
						Secret:           cmp.Or(os.Getenv("VMESS_REMOTE_SECRET"), "---"),
						ListenPort:       30101,
						ListenBackupPort: 30102,
					},
					Wireguard: &singbox.WireguardEndpoint{
						PrivateKey:    cmp.Or(os.Getenv("WIREGUARD_PRIVATE_KEY"), "---"),
						PeerPublicKey: cmp.Or(os.Getenv("WIREGUARD_PEER_PUBLIC_KEY"), "---"),
					},
				}),
			),
		},
	}

	err := p.ExportTo(context.Background(), root)
	if err != nil {
		t.Fatal(err)
	}
}
