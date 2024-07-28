package main

import (
	"piper.octohelm.tech/kubepkg"
	"piper.octohelm.tech/wd"
	"piper.octohelm.tech/file"

	"github.com/v42one/airport/cuepkg/singbox"
	"github.com/v42one/airport/cuepkg/k0sctl"
)

// must be uuid
secret: string
clusters: [Name=string]: [...string] // ips
wireguard?: {...}

hosts: {
	local: wd.#Local & {
	}
}

for region, ips in clusters {
	actions: "\(region)": "singbox": {
		_manifests: kubepkg.#Manifests & {
			kubepkg: singbox.#Server & {
				#values: {
					"clusters":  clusters
					"secret":    secret
					"wireguard": wireguard
				}
			}
		}

		_write: file.#WriteAsYAML & {
			"data": _manifests.manifests
			"with": "asStream": true
			"outFile": {
				"wd":       hosts.local.dir
				"filename": "build/\(region)/manifests/sing-box/sing-box.yaml"
			}
		}
	}

	actions: "\(region)": "cluster": {
		_config: k0sctl.#Config & {
			#values: {
				name: "\(region)"
				hosts: "proxy-\(region)": {
					role: "single"
					ssh: {
						address: ips[0]
					}
					files: [
						{
							name:   "manifests"
							src:    "./manifests"
							dstDir: "/var/lib/k0s/manifests/"
						},
					]
				}
			}
			spec: k0s: {
				version: "1.30.2+k0s.0"
			}
		}

		_write: file.#WriteAsYAML & {
			data: _config
			"outFile": {
				"wd":       hosts.local.dir
				"filename": "build/\(region)/cluster.yaml"
			}
		}
	}
}
