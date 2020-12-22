package clash

import (
	"strconv"
)

#HealthCheck: {
	enable:   bool
	interval: int
	url:      string | *"https://cp.cloudflare.com/"
}

#UrlTest: {
	type:     string | *"url-test"
	url:      string | *"https://cp.cloudflare.com/"
	interval: 300
}

#Config: {
	proxies: [Name=string]: [Type=string]: [IP=string]: [Port=string]: {
		port:   strconv.ParseInt(Port, 10, 64)
		type:   Type
		server: IP
		name:   "\(Name)-\(type)-\(server)-\(port)"
		...
	}
	proxyProviders: [Name=string]: [Type=string]: string
	ruleProviders: [T=string]: [N=string]:        string
	rules: *["MATCH,PROXY"] | [...string]

	proxyGroups: {
		PROXY: {
			"type": "select"
			"proxies": [
				"AUTO", "MANUAL",
			]
		}

		MANUAL: {
			"type": "select"
			"proxies": [
				for n, c in proxyProviders {
					n
				},
				for n, c in proxies {
					n
				},
			]
		}

		FALLBACK: {
			type: "select"
			proxies: [
				"DIRECT",
				"PROXY",
			]
		}

		for n, typeUrl in proxyProviders {
			"\(n)": {
				#UrlTest

				proxies: [
					for t, url in typeUrl {
						"\(n)-\(t)"
					},
				]
			}

			for t, url in typeUrl {
				"\(n)-\(t)": {
					#UrlTest & {
						type: "load-balance"
					}
					use: [
						"\(n)-\(t)",
					]
				}
			}
		}

		for n, typeIpPorts in proxies {
			"\(n)": {
				#UrlTest & {
					type: "load-balance"
				}

				proxies: [
					for type, ipC in typeIpPorts
					for ip, portC in ipC
					for port, c in portC {
						c.name
					},
				]
			}
		}

		AUTO: {
			#UrlTest & {
				type: "url-test"
			}
			proxies: MANUAL.proxies
		}
	}

	output: {
		"proxy-providers": {
			for n, typeURL in proxyProviders
			for t, u in typeURL {
				"\(n)-\(t)": {
					"type": "http"
					"url":  u
					"path": "./\(n)-\(t).yaml"

					"health-check": #HealthCheck & {
						enable:   true
						interval: 300
					}
				}
			}
		}

		"proxies": [
			for n, typeAndConf in proxies
			for t, ipAndConf in typeAndConf
			for ip, portAndConf in ipAndConf
			for port, conf in portAndConf {
				conf
			},
		]

		"proxy-groups": [
			for n, c in proxyGroups {
				{
					c
					name: "\(n)"
				}
			},
		]

		"rule-providers": {
			for b, privders in ruleProviders
			for n, u in privders {
				"\(n)": {
					type:     "http"
					behavior: "domain"
					url:      u
					path:     "./ruleset/\(n).yaml"
					interval: 86400
				}
			}
		}

		"rules": rules

		"port":       7890
		"mixed-port": 7890

		"ipv6":                true
		"allow-lan":           true
		"mode":                "Rule"
		"log-level":           "silent"
		"external-controller": "0.0.0.0:9090"
		"secret":              ""
	}
}
