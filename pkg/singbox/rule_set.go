package singbox

import (
	"fmt"
	"strings"

	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common/json/badoption"

	"github.com/v42one/airport/pkg/runtime"
)

type RuleSet string

func (r RuleSet) Tag() Tag {
	return Tag(strings.TrimPrefix(string(r), "!"))
}

func (r RuleSet) Invert() bool {
	return len(r) > 0 && r[0] == '!'
}

type RemoteRuleSet struct {
	Name   Tag
	Detour Tag
}

func (d RemoteRuleSet) Tag() Tag {
	return d.Name
}

func (d RemoteRuleSet) ApplyTo(o *option.Options) {
	runtime.Apply(o, WithRouteOptions(func(r *option.RouteOptions) {
		remote := option.RemoteRuleSet{
			URL: fmt.Sprintf("https://gh-proxy.com/raw.githubusercontent.com/SagerNet/sing-%s/rule-set/%s.srs", strings.Split(string(d.Name), "-")[0], d.Name),
		}
		if d.Detour != "" {
			remote.HTTPClient = &option.HTTPClientOptions{
				DialerOptions: option.DialerOptions{
					Detour: string(d.Detour),
				},
			}
		}
		r.RuleSet = append(r.RuleSet, option.RuleSet{
			Tag:           badoption.Listable[string]{string(d.Tag())},
			Type:          "remote",
			Format:        "binary",
			RemoteOptions: remote,
		})
	}))
}
