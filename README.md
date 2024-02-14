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
  - One of:
    - `cloud_provider_config`> `provider_settings` > `github_ids` one or more valid github id. The public keys associated to the id will be uplaoded to the servers.

## TODO

- [ ] add documentation to parts of the code (I need to work more with it)
- [ ] using systemctl, setup nginx, upload/download binary, upload service, restart service
- [ ] hetzner minimum viable config (2 boxes, 1 load balancer, 1 postgres, private network)
- [ ] postgres to master/slave
- [ ] o11y at infra level (grafana, tempo, loki, successor-to-prometheus) as a playbook
