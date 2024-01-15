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
- [ ] Terraform variables from environment
- [ ] Generate "diff" inventory: take old TF output and new TF output, and generate a diff inventory
- [ ] Provide `config.yaml` or file to process as parameter, with a default value of `./config.yaml`
- [ ] Add help flag (--help)