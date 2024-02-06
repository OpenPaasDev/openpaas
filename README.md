# OpenPaas Core
Ties together Terraform, Ansible and Runtime execution of commands (local) to create a generic base for automation with these tools.


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

- [ ] add documentation to parts of the code (I need to work more with it)
- [ ] enable remote backend for TF for coordination, based on some flag


- [ ] copy ansible playbook from artaml and use it as template to run nginx as binary, directly
- [ ] using systemctl, setup nginx, upload/download binary, upload service, restart service

- [ ] hetzner minimum viable config (2 boxes, 1 load balancer, 1 postgres, private network)

- [ ] set up postgres and backups? single ndoe with continues backup + daily snapshot?
- [ ] postgres to master/slave

- [ ] o11y at infra level (grafana, tempo, loki, successor-to-prometheus) as a playbook
