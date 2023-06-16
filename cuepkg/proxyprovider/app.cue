package proxyprovider

import (
	"list"
	"encoding/json"
	"encoding/yaml"

	"github.com/v42one/clash-proxy/cuepkg/clash"
	"github.com/v42one/clash-proxy/cuepkg/v2fly"

	kubepkg "github.com/octohelm/kubepkg/cuepkg/kubepkg"
)

#ProxyProvider: {
	#values: {
		clusters: [Name=string]: [...string]
		config: clash.#Config
		secret: string

		port: {
			config: 30000

			main:   30001
			backup: 30002

			from: 30010
			to:   30020
		}

		for name, ips in clusters
		for i, ip in ips {
			config: proxies: "\(name)-\(i)": {
				vmess: "\(ip)": {
					"\(port.main)": {
						uuid:    secret
						alterId: 0
						cipher:  "auto"
					}
					"\(port.backup)": {
						uuid:    secret
						alterId: 0
						cipher:  "auto"
					}
				}
			}
		}
	}

	kubepkg.#KubePkg & {
		metadata: {
			name: _ | *"proxyprovider"
		}

		spec: {
			version: _ | *"5.4.1"

			deploy: {
				kind: "Deployment"
				spec: replicas: _ | *1
			}

			services: "#": {
				ports: "http": #values.port.config
				ports: containers.v2fly.ports
				expose: {
					type: "NodePort"
				}
			}

			containers: "v2fly": {
				image: {
					name: *"docker.io/v2fly/v2fly-core" | string
					tag:  *"v\(spec.version)" | string
				}
				args: ["run", "-c", "\(volumes."config".mountPath)"]

				ports: {
					"tcp-main":   #values.port.main
					"tcp-backup": #values.port.backup

					for i, v in list.Repeat([int], (#values.port.to+1)-#values.port.from) {
						"tcp-\(#values.port.from+i)": #values.port.from + i
					}
				}
			}

			containers: "config": {
				image: {
					name: _ | *"docker.io/library/nginx"
					tag:  _ | *"alpine"
				}
				ports: "http": 80
			}

			_config: v2fly.#V4.#Config & {
				inbounds: [
					{
						tag:      "vmess"
						protocol: "vmess"
						port:     #values.port.main
						settings: {
							clients: [
								{
									id:      "\(#values.secret)"
									alterId: 0
								},
							]
						}
					},
					{
						tag:      "vmess-backup"
						port:     #values.port.backup
						protocol: "vmess"
						settings: {
							clients: [
								{
									id:      "\(#values.secret)"
									alterId: 0
								},
							]
						}
					},
					{
						tag:      "vmess-dynamic"
						port:     #values.port.from
						protocol: "vmess"
						settings: {
							clients: [
								{
									id:      "\(#values.secret)"
									alterId: 0
								},
							]
							detour: {
								to: "vmess-dynamic-ports"
							}
						}
					},
					{
						tag:      "vmess-dynamic-ports"
						protocol: "vmess"
						port:     "\(#values.port.from+1)-\(#values.port.to)"
						settings: {
							default: {
								alterId: 0
							}
						}
						allocate: {
							strategy:    "random"
							concurrency: 2
							refresh:     5
						}
					},
				]
				outbounds: [
					{
						protocol: "freedom"
					},
				]
			}

			volumes: "provider": {
				type:      "ConfigMap"
				mountPath: "/usr/share/nginx/html/ui/\(#values.secret)"
				spec: data: {
					"clash-proxy.yaml": yaml.Marshal(#values.config.output)
				}
			}

			volumes: "config": {
				type:      "ConfigMap"
				mountPath: "/etc/v2fly/config.json"
				subPath:   "config.json"
				spec: data: {
					"config.json": json.Indent(json.Marshal(_config), "", "  ")
				}
			}
		}
	}
}
