K0SCTL = go tool k0sctl

test:
	go test ./...

apply.%:
	cd ./build/$* && $(K0SCTL) apply --config cluster.yaml

reset.%:
	cd ./build/$* && $(K0SCTL) reset --config cluster.yaml

kubeconfig.%:
	cd ./build/$* && $(K0SCTL) kubeconfig --config cluster.yaml > ~/.kube_config/config--$*.yaml

