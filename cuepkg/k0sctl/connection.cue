package k0sctl

#Connection: {
	ssh?: #SSH
}

#SSH: {
	address:  string
	user:     string | *"root"
	port:     int | *22
	keyPath:  string | *"~/.ssh/id_rsa"
	hostKey?: string
	bastion?: #SSH
}
