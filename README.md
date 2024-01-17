# OpenPaas Core
Ties together Terraform, Ansible and Runtime execution of commands (local) to create a generic base for automation with these tools.

A generalisation of the original OpenPaas code, to be incorporated there later.

## Requirements

- [make](https://formulae.brew.sh/formula/make) for build commands.`brew install make`
- [golangci-lint](https://golangci-lint.run/usage/install/) for linting. `brew install golangci-lint`
- [pkgsite](https://pkg.go.dev/golang.org/x/pkgsite/cmd/pkgsite) for local documentation. `go install golang.org/x/pkgsite/cmd/pkgsite@latest`
- [Terraform](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli) for Infrastructure as Code (IaC). Recommended to use [tfenv](https://github.com/tfutils/tfenv) for installation.
- [Ansible](https://www.ansible.com) for automation of server configuration. `brew install ansible`

## TODO

- [ ] modify config to use existing pre-uploaded keys from hetzner for ssh (see test)
- [ ] Terraform variables from environment
- [ ] Generate "diff" inventory: take old TF output and new TF output, and generate a diff inventory
- [ ] Provide `config.yaml` or file to process as parameter, with a default value of `./config.yaml`
- [ ] Complete Command configurationmfor command and main subcommands, including any new flags (take it from the package documentation)
- [ ] Init terraform version, parameterise?
