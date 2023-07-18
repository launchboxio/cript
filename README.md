# cript 

## Development 

```shell
# Start minikube 
$ minikube start --profile cript 
# Install redis for storage
$ helm install redis oci://registry-1.docker.io/bitnamicharts/redis --version 17.13.2 -f deploy/dev/redis.yaml
# Generate the configmap 

```