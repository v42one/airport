package singbox

import (
	"encoding/json"

	kubepkg "github.com/octohelm/kubepkgspec/cuepkg/kubepkg"
)

#Server: {
	#values: {
		clusters: [Name=string]: [...string]
		secret: string
		wireguard: {...}

		port: {
			config: 30100
			main:   30101
			backup: 30102
		}
	}

	kubepkg.#KubePkg & {
		metadata: {
			name: _ | *"sing-box"
		}

		spec: {
			version: _ | *"1.9.3"

			deploy: {
				kind: "Deployment"
				spec: replicas: _ | *1
			}

			services: "#": {
				ports: "http": #values.port.config
				ports: containers."sing-box".ports
				expose: {
					type: "NodePort"
				}
			}

			containers: "sing-box": {
				image: {
					name: *"ghcr.io/sagernet/sing-box" | string
					tag:  *"v\(spec.version)" | string
				}
				args: ["run", "-c", "\(volumes."config".mountPath)"]

				ports: {
					"tcp-main":   #values.port.main
					"tcp-backup": #values.port.backup
				}
			}

			containers: "config": {
				image: {
					name: _ | *"docker.io/library/nginx"
					tag:  _ | *"alpine"
				}
				ports: "http": 80
			}

			volumes: "provider": {
				type:      "ConfigMap"
				mountPath: "/usr/share/nginx/html/ui/\(#values.secret)"
				spec: data: {
					_config: #Config & {
						outbounds: [
							#SelectorOutbound & {
								"tag": "proxy"
								"outbounds": [
									"auto",
									"manual",
								]
							},
							#URLTestOutbound & {
								tag: "auto"
								outbounds: [
									for name, ips in #values.clusters {
										"\(name)"
									},
								]
							},
							#SelectorOutbound & {
								"tag": "manual"
								"outbounds": [
									for name, ips in #values.clusters {
										"\(name)"
									},
								]
							},
							for name, ips in #values.clusters for i, ip in ips for port in [#values.port.main, #values.port.backup] {
								#VmessOutbound & {
									tag:         "\(name)-\(i)-\(port)"
									server:      "\(ip)"
									server_port: port
									uuid:        #values.secret
									alter_id:    0
								}
							},

							for name, ips in #values.clusters {
								#URLTestOutbound & {
									tag: "\(name)"
									outbounds: [
										for i, ip in ips for port in [#values.port.main, #values.port.backup] {
											"\(name)-\(i)-\(port)"
										},
									]
								}
							},
							{
								"type": "direct"
								"tag":  "direct"
							},
							{
								"type": "block"
								"tag":  "block"
							},
							{
								"type": "dns"
								"tag":  "dns-out"
							},
						]
					}

					"sing-box.json": json.Marshal(_config)
				}
			}

			volumes: "config": {
				type:      "ConfigMap"
				mountPath: "/etc/sing-box/config.json"
				subPath:   "config.json"
				spec: data: {
					_config: {
						inbounds: [
							#VmessInbound & {
								tag:         "vmess-in"
								listen_port: #values.port.main
								users: [
									{
										"uuid":    "\(#values.secret)"
										"alterId": 0
									},
								]
								sniff:                      true
								sniff_override_destination: true
							},
							#VmessInbound & {
								tag:         "vmess-in-backup"
								listen_port: #values.port.backup
								users: [
									{
										"uuid":    "\(#values.secret)"
										"alterId": 0
									},
								]
								"sniff":                      true
								"sniff_override_destination": true
							},
						]
						route: {
							rules: [
								{
									"rule_set": [
										"geosite-openai",
									]
									"outbound": "warp-ipv4-out"
								},
							]
							rule_set: [
								{
									"tag":    "geosite-openai"
									"type":   "remote"
									"format": "binary"
									"url":    "https://raw.githubusercontent.com/SagerNet/sing-geosite/rule-set/geosite-openai.srs"
								},
							]
						}
						outbounds: [
							{
								type: "direct"
							},
							{
								"type":            "direct"
								"tag":             "warp-ipv4-out"
								"detour":          "wireguard-out"
								"domain_strategy": "ipv4_only"
							},
							{
								"type":            "direct"
								"tag":             "warp-ipv6-out"
								"detour":          "wireguard-out"
								"domain_strategy": "ipv6_only"
							},
							{
								#values.wireguard

								"type": "wireguard"
								"tag":  "wireguard-out"
							},
						]
					}

					"config.json": json.Marshal(_config)
				}
			}
		}
	}
}
