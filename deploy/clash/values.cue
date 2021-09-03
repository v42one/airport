package clash

#values: {
	image: {
		repository: *"docker.io/dreamacro/clash" | string
		tag:        *"v1.6.5" | string
		pullPolicy: *"IfNotPresent" | string
	}

	proxyPrivders: *{
			"hw-sg": "http://0.0.0.0:30000/proxy.yaml"
	} | {[string]: string}
}
