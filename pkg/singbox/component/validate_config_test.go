package kubepkg

import (
	"os"
	"testing"

	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common/json"
)

func TestValidateGeneratedConfigs(t *testing.T) {
	for _, f := range []string{"./build/client.json", "./build/server.json"} {
		data, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}

		// 提取 rule_set 部分验证（无需 registry context）
		var partial struct {
			Route struct {
				RuleSet []option.RuleSet `json:"rule_set"`
			} `json:"route"`
		}
		if err := json.Unmarshal(data, &partial); err != nil {
			t.Fatalf("%s: %v", f, err)
		}

		for _, rs := range partial.Route.RuleSet {
			if rs.Type != "remote" {
				continue
			}
			if rs.RemoteOptions.HTTPClient == nil {
				t.Fatalf("%s: rule-set %v missing http_client", f, rs.Tag)
			}
			if rs.RemoteOptions.HTTPClient.Detour == "" {
				t.Fatalf("%s: rule-set %v http_client missing detour", f, rs.Tag)
			}
		}
	}
}
