# go-discover tests for Kubernetes

## Prerequisites
- [docker](https://www.docker.com/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [minikube](https://github.com/kubernetes/minikube)

## Running the tests

- `make build`
    - Should create a Minikube VM
    - Should create a Docker image for the Discover binary
    - Should create a deployment for NGINX
    - Should create a deployment for the Discover binary
- `make test`
    - Should test that the Discover binary discovers three
    NGINX instances, and validate that each of the three addresses
    are valid IP addresses.
- `make destroy`
    - Destroys all the Minikube resources
