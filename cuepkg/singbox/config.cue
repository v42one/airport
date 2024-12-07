package singbox

#Config: {
	log:   #Log
	dns:   _ | *#DefaultDNS
	route: _ | *#DefaultRoute
	inbounds: [...] | *[#DefaultClientInbound]
	outbounds: [...]
	experimental: {
		cache_file: {
			enabled: true // required to save rule-set cache
		}
	}
}

#Log: {
	level:     string | *"info"
	timestamp: bool | *true
	output?:   string
	disabled?: bool
}

#Listen: {
	listen:      string | *"::"
	listen_port: int

	tcp_fast_open?:  bool
	tcp_multi_path?: bool
	udp_fragment?:   bool
	udp_timeout?:    int

	detour?:                     string
	sniff?:                      bool
	sniff_override_destination?: bool
	sniff_timeout?:              string

	domain_strategy?:              "ipv4_only" | "ipv6_only" | "prefer_ipv4" | "prefer_ipv6"
	udp_disable_domain_unmapping?: bool
}

#VmessInbound: {
	type: "vmess"
	tag:  string

	users: [
		...{
			uuid:    string
			alterId: int
		},
	]

	#Listen
}

#VmessOutbound: {
	type:        "vmess"
	tag:         string
	server:      string
	server_port: int
	uuid:        string
	security:    "auto"
	alter_id:    int
}

#URLTestOutbound: {
	"type": "urltest"
	"tag":  string
	"outbounds": [...string]
	"url":                         string | *"https://www.gstatic.com/generate_204"
	"interval":                    string | *"1m"
	"tolerance":                   int | *50
	"interrupt_exist_connections": bool | *false
}

#SelectorOutbound: {
	"type": "selector"
	"tag":  string
	"outbounds": [...string]
}
