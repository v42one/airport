package clash

#DefaultRuleProviders: {
	_rules: [Group=string]: [Name=string]: ""
	_rules: {
		domain: {
			reject:       _
			icloud:       _
			apple:        _
			google:       _
			proxy:        _
			direct:       _
			gfw:          _
			greatfire:    _
			"tld-not-cn": _
		}
		ipcidr: {
			cncidr:  _
			lancidr: _
		}
	}

	for g, values in _rules
	for n, v in values {
		"\(g)": "\(n)": "https://gh-proxy.com/https://raw.githubusercontent.com/Loyalsoldier/clash-rules/release/\(n).txt"
	}
}

#DefaultRules: [
	"DOMAIN-SUFFIX,gvt0.com,PROXY",
	"DOMAIN-SUFFIX,gvt1.com,PROXY",
	"DOMAIN-SUFFIX,gvt3.com,PROXY",
	"DOMAIN-SUFFIX,xn--ngstr-lra8j.com,PROXY",
	"DOMAIN-KEYWORD,google,PROXY",
	"PROCESS-NAME,v2ray,DIRECT",
	"PROCESS-NAME,Surge%203,DIRECT",
	"PROCESS-NAME,ss-local,DIRECT",
	"PROCESS-NAME,privoxy,DIRECT",
	"PROCESS-NAME,trojan,DIRECT",
	"PROCESS-NAME,trojan-go,DIRECT",
	"PROCESS-NAME,naive,DIRECT",
	"PROCESS-NAME,Thunder,DIRECT",
	"PROCESS-NAME,DownloadService,DIRECT",
	"PROCESS-NAME,qBittorrent,DIRECT",
	"PROCESS-NAME,Transmission,DIRECT",
	"PROCESS-NAME,fdm,DIRECT",
	"PROCESS-NAME,aria2c,DIRECT",
	"PROCESS-NAME,Folx,DIRECT",
	"PROCESS-NAME,NetTransport,DIRECT",
	"PROCESS-NAME,uTorrent,DIRECT",
	"PROCESS-NAME,WebTorrent,DIRECT",
	"DOMAIN-SUFFIX,innoai.tech,DIRECT",
	"RULE-SET,icloud,DIRECT",
	"RULE-SET,apple,DIRECT",
	"RULE-SET,google,DIRECT",
	"RULE-SET,proxy,PROXY",
	"RULE-SET,direct,DIRECT",
	"GEOIP,,DIRECT",
	"GEOIP,CN,DIRECT",
	"MATCH,FALLBACK",
]
