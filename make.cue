package main

import (
	"encoding/yaml"

	"wagon.octohelm.tech/core"

	"github.com/v42one/clash-proxy/cuepkg/proxyprovider"
	"github.com/v42one/clash-proxy/cuepkg/clash"
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
			_app: (proxyprovider.#ProxyProvider & {
				#values: {
					"config":   actions.config
					"clusters": clusters
					"secret":   secret
				}
			})

			_manifests: yaml.MarshalStream([
					for k, ms in _app.kube if k != "namespace" for n, m in ms {
					m
				},
			])

			_write: core.#WriteFile & {
				contents: _manifests
				path:     "\(region)/proxyprovider.yaml"
			}

			output: _write.output
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
