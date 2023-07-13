export BUILDKIT_HOST =

gen:
	wagon -p . do all --output ./build

apply.%:
	cd ./build/$* && k0sctl apply --config cluster.yaml

kubeconfig.%:
	cd ./build/$* && k0sctl kubeconfig --config cluster.yaml > ~/.kube_config/config--proxy-$*.yaml

reset.%:
	cd ./build/$* && k0sctl reset --config cluster.yaml
