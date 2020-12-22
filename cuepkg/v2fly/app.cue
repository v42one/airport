package v2fly

import (
	"strconv"
	"encoding/json"

	"github.com/innoai-tech/runtime/cuepkg/kube"
)

#V2Fly: kube.#App & {
	#values: {
		server: {
			uuid: string
		}
		upstreams: [Name=string]: [Server=string]: [Port=string]: {
			uuid: string
		}
		expose: [Type=string]: int
	}

	app: {
		name:    "v2fly"
		version: *"5.1.0" | string
	}

	services: "\(app.name)": {
		selector: "app": "\(app.name)"

		ports: "proxy-http": containers."v2fly".ports."proxy-http"
	}

	for t, port in #values.expose {
		services: "\(app.name)-\(t)": {
			selector: "app": "\(app.name)"
			expose: type:    "NodePort"
			ports: "\(t)":   port
		}
	}

	containers: "v2fly": {
		image: {
			name: *"docker.io/v2fly/v2fly-core" | string
			tag:  *"v\(app.version)" | string
		}

		args: ["run", "-c", "\(volumes."config".mountPath)"]

		ports: "proxy-http":  80
		ports: "proxy-vmess": 10086
	}

	_config: #V4.#Config & {
//		log: loglevel: "info"

		api: services: [
			"StatsService",
			"ObservatoryService",
		]

		inbounds: [
			{
				tag:      "api"
				protocol: "dokodemo-door"
				listen:   "127.0.0.1"
				settings: address: "\(listen)"
				port: 8080
			},
			{
				protocol: "http"
				tag:      "inbound-http"
				port:     containers."v2fly".ports."proxy-http"
			},
			{
				tag:      "inbound-vmess"
				protocol: "vmess"
				port:     containers."v2fly".ports."proxy-vmess"
				settings: {
					clients: [
						{
							id:      "\(#values.server.uuid)"
							alterId: 0
						},
					]
				}
			},
		]
		outbounds: [
			for name, sc in #values.upstreams
			for address, portC in sc
			for port, c in portC {
				{
					protocol: "vmess"
					tag:      "naive-\(name)-\(address)-\(port)"
					mux: enabled: true
					settings: {
						vnext: [
							{
								"address": "\(address)"
								"port":    strconv.ParseInt(port, 10, 64)
								"users": [
									{
										"id":      c.uuid
										"alterId": 0
									},
								]
							},
						]
					}
				}
			},
		]
		routing: {
			domainStrategy: "AsIs"
			domainMatcher:  "mph"
			rules: [
				{
					inboundTag: [
						"api",
					]
					outboundTag: "api"
				},
				{
					network:     "tcp,udp"
					balancerTag: "balancer"
				},
			]
			balancers: [
				for name, _ in #values.upstreams {
					{
						tag: "balancer-\(name)"
						selector: [
							"naive-\(name)",
						]
						strategy: type: "random"
					}
				},
				{
					tag: "balancer"
					selector: [
						for name, sc in #values.upstreams {
							"balancer-\(name)"
						},
					]
					strategy: type: "leastPing"
				},
			]
		}
		observatory: {
			subjectSelector: [
				"naive",
			]
			probeURL:      "https://www.google.com/generate_204"
			probeInterval: "10s"
		}
	}

	volumes: "config": {
		mountPath: "/etc/v2fly/config.json"
		subPath:   "config.json"
		source: {
			type: "configMap"
			name: "\(app.name)-config"
			spec: {
				data: {
					"config.json": json.Marshal(_config)
				}
			}
		}
	}
}
