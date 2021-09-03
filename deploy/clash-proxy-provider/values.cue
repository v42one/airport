package shadowsocks

#values: {
	expose: port: *30000 | int

	image: {
		hub:        *"docker.io/morlay" | string
		name:       *"clash-proxy-provider" | string
		tag:        *"0.2.0" | string
		pullPolicy: *"IfNotPresent" | string
	}
}
