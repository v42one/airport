package singbox

import (
	"strings"
)

#DefaultClientInbound: {
	type: "tun"
	address: [
		"172.19.0.1/30",
	]
	auto_route:                 true
	strict_route:               true
	stack:                      "system"
	sniff:                      true
	sniff_override_destination: false
}

#DefaultRoute: {
	rules: [
		{
			protocol: "dns"
			outbound: "dns-out"
		},
		{
			"domain_suffix": [
				"huggingface.co",
				"hf.co",
				"pub.dev",
				"pkg.dev",
			]
			outbound: "proxy"
		},
		{
			type: "logical"
			mode: "and"
			rules: [
				{
					rule_set: "geosite-cn"
					invert:   true
				},
				{
					rule_set: "geoip-cn"
					invert:   true
				},
			]
			outbound: "proxy"
		},
		{
			type: "logical"
			mode: "and"
			rules: [
				{
					rule_set: "geosite-cn"
				},
				{
					rule_set: "geoip-cn"
				},
			]
			outbound: "direct"
		},
		{
			"ip_is_private": true
			"outbound":      "direct"
		},
		{
			ip_cidr: []
			outbound: "direct"
		},
	]

	rule_set: [
		for tag in [
			"geosite-cn",
			"geoip-cn",
		] {
			{
				"tag":    tag
				"type":   "remote"
				"format": "binary"
				"url": [
					if strings.HasPrefix(tag, "geosite-") {
						"https://raw.githubusercontent.com/SagerNet/sing-geosite/rule-set/\(tag).srs"
					},
					"https://raw.githubusercontent.com/SagerNet/sing-geoip/rule-set/\(tag).srs",
				][0]
				"download_detour": "proxy"
			}
		},
	]

	final:                 "direct"
	auto_detect_interface: true
}

#DefaultDNS: {
	"servers": [
		{
			"tag":              "dns_proxy"
			"address":          "https://1.1.1.1/dns-query"
			"address_resolver": "dns_resolver"
			"strategy":         "ipv4_only"
			"detour":           "proxy"
		},
		{
			"tag":              "dns_direct"
			"address":          "h3://dns.alidns.com/dns-query"
			"address_resolver": "dns_resolver"
			"strategy":         "ipv4_only"
			"detour":           "direct"
		},
		{
			"tag":      "dns_resolver"
			"address":  "223.5.5.5"
			"strategy": "ipv4_only"
			"detour":   "direct"
		},
		{
			"tag":     "dns_block"
			"address": "rcode://refused"
		},
	]
	"rules": [
		{
			"domain_suffix": [
				// https://www.feishu.cn/hc/zh-CN/articles/360044683233-%E9%85%8D%E7%BD%AE%E4%BC%81%E4%B8%9A%E5%86%85%E7%BD%91%E9%98%B2%E7%81%AB%E5%A2%99%E5%9F%9F%E5%90%8D%E5%92%8C%E7%99%BD%E5%90%8D%E5%8D%95
				"feishu.net",
				"feishu.cn",
				"larkoffice.com",
				"feishucdn.com",
				"zjurl.cn",
				"snssdk.com",
				"pstatp.com",
				"byteimg.com",
				"bytedance.net",
				"bytedance.com",
				"byted-static.com",
				"bytegoofy.com",
				"weboffice.feishu-3rd-party-services.com",
				"bytehwm.com",
				"ttwebview.com",
				"bytegecko.com",
				"bytescm.com",
				"kundou.cn",
				"bytetos.com",
				"zijieapi.com",
				"byteeffecttos.com",
				"bytednsdoc.com",
				"bytedanceapi.com",
				"volcvideo.com",
				"feishuimg.com",
				"feishuapp.cn",
				"getfeishu.cn",
				"feishupkg.com",
			]
			"server": "dns_direct"
		},
		{
			"domain_suffix": [
				"pkg.dev",
			]
			"server": "dns_proxy"
		},
		{
			"rule_set": "geosite-cn"
			"invert":   true
			"server":   "dns_proxy"
		},
		{
			"rule_set": "geosite-cn"
			"server":   "dns_direct"
		},
		{
			"outbound": [
				"any",
			]
			"server": "dns_resolver"
		},
	]
}
