# Pod-Deployer

Simple app to deploy a pod to kubernetes cluster from getting user inputs through UI. It uses kubeconfig from the machine it is accessed.

![index page](img/index.png)

## Motivation for the project

The motivation for creating this project is to get familiarize with [client-go](https://github.com/kubernetes/client-go) and to learn go web programming.

## Usage

#### Pre-requisite

- Running kubernetes cluster and kube config available in the machine

### Build Locally

- Since pre-built binary is not included you need to build yourself hence you need to have go `1.16+` installed in your system

Build the application

```go
git clone https://github.com/kube-go/pod-deployer.git
cd pod-deployer
go build
```

#### Execute

```go
./pod-deployer
```

### Docker way

```go
DATE=$(date +"%Y%m%d")
docker build . -t pod-deployer:$DATE
```

#### Run as docker container

*Note:* Mounted kube context directory as volume in container

```go
docker run -p 8080:8080 -v $HOME/.kube:/home/deployer/.kube pod-deployer:$DATE
```


## TODO

- [x] Provide exact reason for failure in UI

- [ ] Option to view pod status, description, logs in UI

- [ ] Include binary

## References

- [client-go](https://github.com/kubernetes/client-go)

- [gowebexamples](https://gowebexamples.com/)
