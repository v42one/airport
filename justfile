k0sctl := "go tool k0sctl"

fmt:
    go tool fmt .

test:
    go test --count=1 ./...

apply cluster:
    cd ./build/{{ cluster }} && {{ k0sctl }} apply --config cluster.yaml

reset cluster:
    cd ./build/{{ cluster }} && {{ k0sctl }} reset --config cluster.yaml

kubeconfig cluster:
    cd ./build/{{ cluster }} && {{ k0sctl }} kubeconfig --config cluster.yaml > ~/.kube_config/config--$*.yaml

