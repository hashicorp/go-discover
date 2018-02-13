help:
	@echo "build   - minikube get plan apply"
	@echo "test    - test go-discover"
	@echo "destroy - minikube destroy"
	@echo "clean   - remove temp files"

compile:
	@echo "*** Compiling github.com/hashicorp/go-discover/cmd/discover"
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -i -ldflags '-s -w' github.com/hashicorp/go-discover/cmd/discover

build: compile
	@echo "*** Creating Minikube environment"
	@echo "*** This takes ~ 5min"
	@echo
	@minikube start --memory 8192 --cpus 4 --vm-driver=virtualbox --extra-config=kubelet.CAdvisorPort=4194
	@sleep 30
	@minikube addons enable dashboard
	@minikube addons enable default-storageclass
	@minikube addons enable heapster
	@minikube addons enable kube-dns
	@sleep 60
	@minikube dashboard
	@echo "*** Creating Docker image within Minikube environment"
	@echo
	@eval $(minikube docker-env)
	@docker build -t hashicorp/discover:latest .
	@echo "*** Deploying NGINX so that its discoverable"
	@echo
	@kubectl create -f nginx.yaml
	@sleep 10
	@kubectl create -f discover.yaml

test:
	@echo "*** Running go-discover test"
	@./test.bash

destroy:
	@echo "*** Destroying Minikube environment"
	@echo
	minikube stop
	minikube delete

clean:
	rm -f discover

.PHONY: compile build test destroy clean
