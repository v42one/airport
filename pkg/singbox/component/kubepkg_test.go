package kubepkg

import (
	"cmp"
	"os"
	"testing"

	"github.com/sagernet/sing/common/json"
	"github.com/v42one/airport/pkg/singbox"
)

func TestSingBox(t *testing.T) {
	s := &SingBox{
		Version:    "1.11.3",
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
	}

	t.Run("client options", func(t *testing.T) {
		if err := WriteJSONFile("./build/client.json", s.ClientConfig()); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("server options", func(t *testing.T) {
		if err := WriteJSONFile("./build/server.json", s.ServerConfig()); err != nil {
			t.Fatal(err)
		}
	})
}

func WriteJSONFile(filename string, t any) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	e := json.NewEncoder(f)
	e.SetIndent("", "  ")
	return e.Encode(t)
}
