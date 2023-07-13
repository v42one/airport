package k0sctl

#Cluster: {
	apiVersion: "k0sctl.k0sproject.io/v1beta1"
	kind:       "Cluster"
	metadata: {
		name: string
	}
	spec: #Spec
}

#Spec: {
	hosts: [...#Host]
	k0s?: #K0s
}

#Host: {
	role:              "single" | "controller" | "worker" | "controller+worker"
	reset?:            bool
	privateInterface?: string
	privateAddress?:   string
	dataDir?:          string
	environment?: [Name=string]: string

	installFlags?: [...string]
	files?: [...#UploadFile]
	os?:       string
	hostname?: string
	noTaints?: bool
	hooks?:    #Hooks

	uploadBinary?:  bool
	k0sBinaryPath?: string

	#Connection
}

// Hooks define a list of hooks such as hooks["apply"]["before"] = ["ls -al", "rm foo.txt"]
#Hooks: [Action=string]: [Hook=string]: [...string]

#UploadFile: {
	name?:    string
	src:      string
	dstDir?:  string
	dst?:     string
	perm:     int | *0o755
	dirPerm?: int
	user?:    string
	group?:   string
}

#K0s: {
	version?:       string
	dynamicConfig?: bool
	config: {...}
}
