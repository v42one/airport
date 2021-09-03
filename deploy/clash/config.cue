package clash

_proxyGroups: [Name=_]: {
	name: Name
}

_proxyGroups: "PROXY": {
	type: "select"
	proxies: ["AUTO", "MANUAL"]
}

_proxyGroups: "MANUAL": {
	type: "select"
	proxies: [
		for name, url in #values.proxyPrivders {
			"\(name)"
		},
	]
}

_proxyGroups: "AUTO": {
	type:     "url-test"
	url:      "http://www.gstatic.com/generate_204"
	interval: 300
	proxies: [
		for name, url in #values.proxyPrivders {
			"\(name)"
		},
	]
}

_proxyGroups: {
	for name, url in #values.proxyPrivders {
		"\(name)": {
			type:     "load-balance"
			url:      "http://www.gstatic.com/generate_204"
			interval: 300
			use: [
				for name, url in #values.proxyPrivders {
					"\(name)"
				},
			]
		}
	}
}

_config: {
	"port":                7890
	"socks-port":          7891
	"redir-port":          7892
	"allow-lan":           true
	"mode":                "Rule"
	"log-level":           "info"
	"external-controller": "0.0.0.0:9090"
	"rules": [
		"MATCH,PROXY",
	]
	"proxy-groups": [
		for group in _proxyGroups {
			group
		},
	]
	"proxy-providers": {
		for n, _ in #values.proxyPrivders {
			"\(n)": {
				url:  #values.proxyPrivders[n]
				path: "\(n).yaml"
				type: "http"
				"health-check": {
					enable:   true
					interval: 300
					url:      "http://www.gstatic.com/generate_204"
				}
			}
		}
	}
}
