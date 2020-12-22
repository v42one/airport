# Clash Proxy

Auto generate proxy.yaml as provider-provider

# setup k3s

```
export PUBLIC_IP=<public_ip>

curl -sfL https://get.k3s.io | sh -s - --node-external-ip=${PUBLIC_IP}

cat /etc/rancher/k3s/k3s.yaml | sed "s/127.0.0.1/${PUBLIC_IP}/g" > clipboard
```