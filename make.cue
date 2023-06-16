package main

import (
	"wagon.octohelm.tech/core"

	"github.com/v42one/clash-proxy/cuepkg/proxyprovider"
	"github.com/v42one/clash-proxy/cuepkg/clash"
	"github.com/octohelm/kubepkg/cuepkg/kubepkgcli"
)

// must be uuid
secret: string

clusters: [Name=string]: [...string] // ips

actions: config: clash.#Config & {
	rules:         clash.#DefaultRules
	ruleProviders: clash.#DefaultRuleProviders
}

for region, ips in clusters {
	actions: cluster: {
		"\(region)": {
			_manifests: kubepkgcli.#Manifests & {
				image: tag: "v0.5.1-20230616092629-29f3314eb178"
				kubepkg: proxyprovider.#ProxyProvider & {
					#values: {
						"config":   actions.config
						"clusters": clusters
						"secret":   secret
					}
				}
				flags: {
					"--output": "/\(region)/proxyprovider.yaml"
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
	}
}

actions: cluster: all: core.#Merge & {
	inputs: [
		for region, ips in clusters {
			actions.cluster["\(region)"].output
		},
	]
}
