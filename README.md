![facebook_cover_photo_1.png](static%2Fbranding%2Ffacebook_cover_photo_1.png)

# Container Risk Inspection and Protection Tool 

![GitHub](https://img.shields.io/github/license/launchboxio/cript)
![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/launchboxio/cript)
![GitHub release (with filter)](https://img.shields.io/github/v/release/launchboxio/cript)

## Development

```shell
# Start minikube 
$ minikube start --profile cript 
# Install redis for storage
$ helm install redis oci://registry-1.docker.io/bitnamicharts/redis --version 17.13.2 -f deploy/dev/redis.yaml
# Build the image for scanning purposes
$ eval $(minikube docker-env --profile cript)
# Start the operator 
$ make install run
```
