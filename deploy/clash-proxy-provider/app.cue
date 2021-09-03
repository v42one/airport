package shadowsocks

import (
	"github.com/octohelm/cuem/release"
)

release.#Release & {
	#name:      "clash-proxy-provider"
	#namespace: "clash-proxy"
	#context:   "hw-sg"

	spec: serviceAccounts: "\(#name)": {
		#role: "ClusterRole"
		#rules: [
			{
				verbs: ["*"]
				apiGroups: ["*"]
				resources: ["*"]
			},
			{
				verbs: ["*"]
				nonResourceURLs: ["*"]
			},
		]
	}

	spec: services: "\(#name)-expose": {
		spec: selector: app: "\(#name)"
		spec: type: "NodePort"
		spec: ports: [
			{
				name:       "tcp-\(#values.expose.port)"
				port:       #values.expose.port
				targetPort: 80
				protocol:   "TCP"
				nodePort:   #values.expose.port
			},
		]
	}

	spec: deployments: "\(#name)": {
		spec: replicas: 1
		spec: template: spec: serviceAccountName: "\(#name)"

		#containers: "clash-proxy-provider": {
			image:           "\(#values.image.hub)/\(#values.image.name):\(#values.image.tag)"
			imagePullPolicy: "\(#values.image.pullPolicy)"

			#envVars: {
				WATCH_NAMESPACE: "\(#namespace)"
			}

			#ports: {
				tcp: 80
			}
		}
	}

}
