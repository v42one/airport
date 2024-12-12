package srsutil

import (
	"bytes"
	"fmt"
	"maps"
	"net"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/sagernet/sing-box/common/srs"
)

func Test(t *testing.T) {
	raw, err := os.ReadFile("../../build/rule-site/geosite-github.srs")
	if err != nil {
		t.Fatal(err)
	}

	rs, err := srs.Read(bytes.NewBuffer(raw), true)
	if err != nil {
		t.Fatal(err)
	}

	hosts := map[string]string{}

	resolveHost := func(host string) {
		ips, err := net.LookupIP(host)
		if err != nil {
			return
		}
		for _, ip := range ips {
			if !strings.Contains(ip.String(), ":") {
				hosts[ip.String()] = host
			}
		}
	}

	for _, r := range rs.Options.Rules {
		for _, domain := range r.DefaultOptions.Domain {
			resolveHost(domain)
		}

		for _, domain := range r.DefaultOptions.DomainSuffix {
			resolveHost(domain)
		}
	}

	for _, ip := range slices.Sorted(maps.Keys(hosts)) {
		fmt.Println(ip, hosts[ip])
	}
}
