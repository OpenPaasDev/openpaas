# OpenPaas Core
Ties together Terraform, Ansible and Runtime execution of commands (local) to create a generic base for automation with these tools.

A generalisation of the original OpenPaas code, to be incorporated there later.

## Requirements

- [make](https://formulae.brew.sh/formula/make) for build commands.`brew install make`
- [golangci-lint](https://golangci-lint.run/usage/install/) for linting. `brew install golangci-lint`
- [pkgsite](https://pkg.go.dev/golang.org/x/pkgsite/cmd/pkgsite) for local documentation. `go install golang.org/x/pkgsite/cmd/pkgsite@latest`
- [Terraform](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli) for Infrastructure as Code (IaC). Recommended to use [tfenv](https://github.com/tfutils/tfenv) for installation.
- [Ansible](https://www.ansible.com) for automation of server configuration. `brew install ansible`

### Configuration

- Your `Hetzner` token must be available as the environment variable `HETZNER_TOKEN`
- You must provide your config in `config.yaml`
  - in `cloud_provider_config`> `allowed_ips` you must have the IP of the machine doing the deployment, otherwise the firewall will block the connection
  - in `cloud_provider_config`> `provider_settings` > `ssh_keys` you must have the fingerprint of at least 1 ssh key already uploaded to Heztner
  - in `cloud_provider_config`> `ssh_key` you must have the path to a private ssh key matching one of the public keys provided as authorised ssh keys

## TODO


- [ ] Terraform variables from environment


- [ ] merge with bootstrap ansible branch
- [ ] add documentation to parts of the code

- [ ] copy ansible playbook from artaml and use it as template to run nginx as binary, directly
- [ ] using systemctl, setup nginx, upload/download binary, upload service, restart service
- [ ] enable remote backend for coordination, based on some flag

- [ ] hetzner minimum viable config (2 boxes, 1 load balancer, 1 postgres, private network)
- [ ] Generate "diff" inventory: take old TF output and new TF output, and generate a diff inventory

- [ ] Issue? the firewall only allows predefined ip, how to access if ip cycles? may need a default strong root password?

