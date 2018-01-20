SHELL=/usr/bin/env bash

help:
	@echo "build   - terraform get plan apply"
	@echo "test    - test go-discover"
	@echo "destroy - terraform destroy"
	@echo "clean   - remove temp files and terraform state"

prerequisites:
	@if ! command -v az >/dev/null; then \
		echo "!!! azure-cli is required."; \
		echo "!!! Follow these instructions to install azure-cli: https://docs.microsoft.com/en-us/cli/azure/install-azure-cli"; \
		exit 1; \
	fi
	@if ! command -v jq >/dev/null; then \
		echo "!!! jq is required."; \
		echo "!!! Follow these instructions to install jq: https://stedolan.github.io/jq/download/"; \
		exit 1; \
	fi
	@if ! command -v envsubst >/dev/null; then \
		echo "!!! envsubst is required."; \
		echo "!!! Follow these instructions to install envsubst: https://www.gnu.org/software/gettext/"; \
		exit 1; \
	fi

build: prerequisites compile sshkeygen
	@echo "*** Creating service principal and role"
	@echo
	@./setup_vmss_account.bash
	@echo "*** Creating Terraform environment"
	@echo "*** This takes ~ 5min"
	@echo "*** See tf.log for progress"
	@echo
	terraform init  >> tf.log 2>&1
	terraform get   >> tf.log 2>&1
	terraform plan  >> tf.log 2>&1
	terraform apply >> tf.log 2>&1

compile:
	@echo "*** Compiling github.com/hashicorp/go-discover/cmd/discover"
	GOOS=linux GOARCH=amd64 go build -i -ldflags '-s -w' github.com/hashicorp/go-discover/cmd/discover

sshkeygen:
	@echo "*** Generating ssh key tf_rsa for Terraform"
	test -r tf_rsa || ssh-keygen -t rsa -b 2048 -f tf_rsa -N '' -C 'terraform'

test:
	@echo "*** Running go-discover test on Azure VMSS"
	@./test.bash

destroy:
	@echo "*** Destroying service principal and role"
	@echo
	@./remove_vmss_account.bash
	@echo "*** Destroying Terraform environment"
	@echo "*** This takes ~ 2min. See tf.log for progress"
	@echo "*** See tf.log for progress"
	@echo
	terraform destroy -force >> tf.log 2>&1

clean:
	rm -f discover
	rm -rf .terraform
	rm -f terraform.tfstate{,.backup}
	rm -f tf_rsa{,.pub}
	rm -f tf.log
	rm -rf cli-outputs

.PHONY: compile sshkeygen build test destroy clean prerequisites
