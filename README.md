# OpenPaas Core
Ties together Terraform, Ansible and Runtime execution of commands (local) to create a generic base for automation with these tools.

A generalisation of the original OpenPaas code, to be incorporated there later.


## TODO
- [ ] Terraform variables from environment
- [ ] Generate "diff" inventory: take old TF output and new TF output, and generate a diff inventory


# Bootstrap server
curl -sfL https://get.k3s.io | INSTALL_K3S_VERSION=v1.21.5+k3s1 sh -s - server --cluster-init --flannel-iface eth1 --node-ip 10.0.0.X --bind-address 10.0.0.X --advertise-address 10.0.0.X

sudo cat /var/lib/rancher/k3s/server/node-token


# other servers
curl -sfL https://get.k3s.io | INSTALL_K3S_VERSION=v1.21.5+k3s1 K3S_TOKEN=${NODE_TOKEN} sh -s - server --server https://10.0.0.2:6443 --flannel-iface eth1 --node-ip 10.0.0.Y --bind-address 10.0.0.Y --advertise-address 10.0.0.Y

# agents
curl -sfL https://get.k3s.io | INSTALL_K3S_VERSION=v1.21.5+k3s1 K3S_TOKEN=${NODE_TOKEN} K3S_URL=https://10.0.0.2:6443 sh -s - agent --flannel-iface eth1 --node-ip 10.0.0.Z


/etc/rancher/k3s/registries.yaml file.

Add your registry credentials:

mirrors:
  "myregistrydomain.com":
    endpoint:
      - "https://myregistrydomain.com"
configs:
  "myregistrydomain.com":
    auth:
      username: myusername
      password: mypassword


## Upgrade
Server nodes first:

kubectl drain <NODE_NAME> --ignore-daemonsets --delete-local-data
curl -sfL https://get.k3s.io | INSTALL_K3S_VERSION=<VERSION> sh -
kubectl uncordon <NODE_NAME>

## backup
k3s etcd-snapshot save
ls /var/lib/rancher/k3s/server/db/snapshots


https://chat.openai.com/c/0352ca9e-7841-4e36-88c7-f29c6f445e6f