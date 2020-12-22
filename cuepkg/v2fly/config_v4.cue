package v2fly

#V4: {
	#Config: {
		inbounds: [...#Inbound]
		outbounds: [...#Outbound]
		routing?:     #Routing
		observatory?: #Observatory
		log:          #Log
		api: {
			tag: "api"
			services: [...string]
		}
		stats: {}
		policy: {
			levels: {
				"0": {
					statsUserUplink:   true
					statsUserDownlink: true
				}
				"1": {
					statsUserUplink:   true
					statsUserDownlink: true
				}
			}
			"system": {
				statsInboundUplink:   true
				statsInboundDownlink: true
			}
		}
	}

	#Log: {
		loglevel: *"warning" | "debug" | "info" | "warning" | "error" | "none"
		access?:  string
		error?:   string
	}

	#Inbound: {
		protocol: string
		tag:      *"\(protocol)" | string
		port:     string | number
		listen:   *"0.0.0.0" | string
		settings: {...}
		allocate?: _
	}

	#Outbound: {
		sendThrough: *"0.0.0.0" | string
		protocol:    string
		tag:         *"\(protocol)" | string
		settings: {...}
		mux?: enabled: bool
	}

	#Routing: {
		domainStrategy: *"AsIs" | "AsIs" | "IPIfNonMatch" | "IPOnDemand"
		domainMatcher?: "linear" | "mph"
		rules?: [...#Rule]
		balancers?: [...#Balancer]
	}

	// https://www.v2fly.org/config/routing.html#ruleobject
	#Rule: {
		type:           "field"
		domainMatcher?: "linear" | "mph"
		domains?: [...string]
		ip?: [...string]
		port?:       string | int
		sourcePort?: string | int
		network?:    "tcp" | "udp" | "tcp,udp"
		protocol?: [...string]
		inboundTag?: [...string]
		balancerTag?: string
		outboundTag?: string
	}

	#Balancer: {
		tag: string
		selector: [...string]
		strategy: type: "random" | "leastPing"
	}

	#Observatory: {
		subjectSelector: [...string]
		probeURL:      string
		probeInterval: string
	}
}
