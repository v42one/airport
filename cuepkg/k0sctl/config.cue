package k0sctl

#Config: {
	#values: {
		name: string

		hosts: [Name=string]: #Host & {
			hostname:     Name
			role:         _ | *"controller+worker"
			installFlags: _ | *[
					"--disable-components helm",
			]
		}

		ipv6?: bool
	}

	#Cluster & {
		metadata: {
			name: "\(#values.name)"
		}

		spec: hosts: [
			for h in #values.hosts {
				h
			},
		]

		spec: k0s: {
			config: {
				if #values.ipv6 & true != _|_ {
					spec: network: {
						provider: "calico"
						calico: mode: "bird"
						dualStack: {
							enabled:         true
							IPv6podCIDR:     "fd00::/108"
							IPv6serviceCIDR: "fd01::/108"
						}
					}
				}

				spec: extensions: {
					storage: type: "openebs_local_storage"
				}
			}
		}
	}
}
