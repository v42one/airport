package shadowsocks

#values: {
	cipher: *"xchacha20-ietf-poly1305" | string
	ports:  *[30001, 30002] | [int]

	image: {
		hub:        *"docker.io/shadowsocks" | string
		name:       *"shadowsocks-libev" | string
		tag:        *"latest" | string
		pullPolicy: *"IfNotPresent" | string
	}
}
