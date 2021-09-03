package clash

import (
	"encoding/yaml"

	"github.com/octohelm/cuem/release"
)

_configYAML: yaml.Marshal(_config)

release.#Release & {
	#name:      "clash"
	#namespace: "registry"

	spec: configMaps: "\(#name)-config": data: "config.yaml": _configYAML

	spec: deployments: "\(#name)": {
		#volumes: "config": {
			mountPath: "/etc/clash"
			volume: configMap: name: "\(#name)-config"
		}

		#containers: "clash": {
			image:           "\(#values.image.repository):\(#values.image.tag)"
			imagePullPolicy: "\(#values.image.pullPolicy)"
			args: ["-f", "\(#volumes.config.mountPath)/config.yaml"]

			#ports: {
				"http-9090":     9090
				"http-proxy":    7890
				"sockets-proxy": 7891
			}
		}
	}
}
