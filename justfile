k0sctl := "go tool k0sctl"

test:
    go test ./...

apply cluster:
    cd ./build/{{ cluster }} && {{ k0sctl }} apply --config cluster.yaml

reset cluster:
    cd ./build/{{ cluster }} && {{ k0sctl }} reset --config cluster.yaml

kubeconfig cluster:
    cd ./build/{{ cluster }} && {{ k0sctl }} kubeconfig --config cluster.yaml > ~/.kube_config/config--$*.yaml
