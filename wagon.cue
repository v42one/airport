package main

import (
	"encoding/yaml"
	"wagon.octohelm.tech/core"

	"github.com/v42one/airport/cuepkg/singbox"
	"github.com/v42one/airport/cuepkg/k0sctl"
	"github.com/octohelm/kubepkg/cuepkg/kubepkgcli"
)

// must be uuid
secret: string
clusters: [Name=string]: [...string] // ips
wireguard?: {...}

for region, ips in clusters {
	actions: "singbox": "\(region)": {
		_manifests: kubepkgcli.#Manifests & {
			image: tag: "v0.5.3-0.20230710081347-f7c2f4a798c6"
			kubepkg: singbox.#Server & {
				#values: {
					"clusters":  clusters
					"secret":    secret
					"wireguard": wireguard
				}
			}
			flags: {
				"--namespace": "default"
				"--output":    "/\(region)/manifests/sing-box/sing-box.yaml"
			}
		}

		_copy: core.#Copy & {
			contents: _manifests.output.rootfs
			include: [
				"\(region)",
			]
		}

		output: _copy.output
	}

	actions: "cluster": "\(region)": {
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
				version: "1.28.9+k0s.0"
			}
		}

		core.#WriteFile & {
			path:     "\(region)/cluster.yaml"
			contents: yaml.Marshal(_config)
		}
	}
}

actions: all: core.#Merge & {
	inputs: [
		for region, ips in clusters {
			actions.singbox["\(region)"].output
		},
		for region, ips in clusters {
			actions.cluster["\(region)"].output
		},
	]
}
