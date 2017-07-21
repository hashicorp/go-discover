help:
	@echo "build   - terraform get plan apply"
	@echo "test    - test go-discover"
	@echo "destroy - terraform destroy"
	@echo "clean   - remove temp files and terraform state"

creds:
	@if [ -z "$$ARM_TENANT_ID" -o -z "$$ARM_SUBSCRIPTION_ID" -o -z "$$ARM_CLIENT_ID" -o -z "$$ARM_CLIENT_SECRET" ] ; then \
		echo "ARM_{TENANT_ID,SUBSCRIPTION_ID,CLIENT_ID,CLIENT_SECRET} not found" ; \
		exit 1 ; \
	fi

build: creds compile sshkeygen
	@echo "*** Creating Terraform environment"
	@echo "*** This takes ~ 5min"
	@echo "*** See tf.log for progress"
	@echo
	terraform get   >> tf.log 2>&1
	terraform plan  >> tf.log 2>&1
	terraform apply >> tf.log 2>&1

compile:
	@echo "*** Compiling github.com/hashicorp/go-discover/cmd/discover"
	GOOS=linux GOARCH=amd64 go build -i -ldflags '-s -w' github.com/hashicorp/go-discover/cmd/discover

sshkeygen:
	@echo "*** Generating ssh key tf_rsa for Terraform"
	test -r tf_rsa || ssh-keygen -t rsa -b 2048 -f tf_rsa -N '' -C 'terraform'

test: creds
	@echo "*** Running go-discover test on Azure"
	@./test.bash

destroy: creds
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

.PHONY: creds compile sshkeygen build test destroy clean
