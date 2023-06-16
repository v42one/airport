export BUILDKIT_HOST =

dump:
	wagon -p . do cluster all --output ./build

MANIFESTS_ROOT=/var/lib/k0s/manifests/proxyprovider

apply.%:
	ssh root@proxy-$* "mkdir -p $(MANIFESTS_ROOT)"
	scp -r build/$*/proxyprovider.yaml root@proxy-$*:$(MANIFESTS_ROOT)/proxyprovider.yaml
	ssh root@proxy-$* "sudo k0s kubectl apply -f $(MANIFESTS_ROOT)/proxyprovider.yaml"

install:
	curl -sSLf https://get.k0s.sh | sudo sh
	sudo k0s config create > /etc/k0s/k0s.yaml
	sudo k0s install controller --single
	sudo k0s start

ping:
	ping fra-de-ping.vultr.com -c 100
#	ping hnd-jp-ping.vultr.com -c 5
#	ping syd-au-ping.vultr.com -c 5
#	ping par-fr-ping.vultr.com -c 5
