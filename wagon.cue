package main

import (
	"encoding/yaml"
	"wagon.octohelm.tech/core"

	"github.com/v42one/clash-proxy/cuepkg/proxyprovider"
	"github.com/v42one/clash-proxy/cuepkg/clash"
	"github.com/v42one/clash-proxy/cuepkg/k0sctl"
	"github.com/octohelm/kubepkg/cuepkg/kubepkgcli"
)

// must be uuid
secret: string

clusters: [Name=string]: [...string] // ips

config: clash.#Config & {
	rules:         clash.#DefaultRules
	ruleProviders: clash.#DefaultRuleProviders
}

for region, ips in clusters {
	actions: "proxyprovider": "\(region)": {
		_manifests: kubepkgcli.#Manifests & {
			image: tag: "v0.5.3-0.20230710081347-f7c2f4a798c6"
			kubepkg: proxyprovider.#ProxyProvider & {
				#values: {
					"config":   config
					"clusters": clusters
					"secret":   secret
				}
			}
			flags: {
				"--namespace": "default"
				"--output":    "/\(region)/manifests/proxyprovider/proxyprovider.yaml"
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
				version: "1.27.3+k0s.0"
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
			actions.proxyprovider["\(region)"].output
		},
		for region, ips in clusters {
			actions.cluster["\(region)"].output
		},
	]
}
