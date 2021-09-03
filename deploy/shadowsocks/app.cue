package shadowsocks

import (
	"strings"
	"strconv"
	"uuid"

	"github.com/octohelm/cuem/release"
)

release.#Release & {
	#name:      "shadowsocks"
	#namespace: "clash-proxy"
	#context:   "hw-sg"

	spec: services: "shadowsocks-expose-tcp": {
		spec: selector: app: "\(#name)"
		spec: type: "NodePort"
		spec: ports: [
			for np in #values.ports {
				{
					name:       "np-tcp-\(strconv.FormatInt(np, 10))"
					port:       np
					targetPort: 8388
					protocol:   "TCP"
					nodePort:   np
				}
			},
		]
	}

	spec: deployments: "\(#name)": {
		_password : "\(uuid.SHA1("6ba7b810-9dad-11d1-80b4-00c04fd430c8", ''))"

		metadata: labels: "clash-proxy-type": "ss"

		spec: template: metadata: annotations: {
			"clash-proxy-cipher":   "\(#values.cipher)"
			"clash-proxy-password": "\(_password)"
			"clash-proxy-ports":    "\(strings.Join([ for p in #values.ports {strconv.FormatInt(p, 10)}], ","))"
		}

		spec: replicas: 3

		#containers: "shadowsocks": {
			image:           "\(#values.image.hub)/\(#values.image.name):\(#values.image.tag)"
			imagePullPolicy: "\(#values.image.pullPolicy)"

			#envVars: {
				ARGS:      "-v"
				DNS_ADDRS: "1.1.1.1"
				METHOD:    "\(#values.cipher)"
				PASSWORD:  "\(_password)"
			}

			#ports: {
				tcp: 8388
			}
		}
	}
}
