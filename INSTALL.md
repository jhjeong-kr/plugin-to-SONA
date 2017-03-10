
# Installing plugin-to-SONA

## Prerequisite: Dep

[dep](https://github.com/golang/dep) is an up-and-coming dependency management tool,
which has the goal of being accepted as part of the standard go toolchain. It
is currently pre-alpha. However, it comes the closest to working easily out of
the box.

```sh
$ go get github.com/golang/dep
$ go install github.com/golang/dep/cmd/dep

$ dep init
$ dep ensure k8s.io/client-go@^2.0.0
```

This will set up a `vendor` directory in your current directory.

## Build
```sh
$ make build
```

## Run

### In-clustering
```sh
# Make sure the docker image is accessable by kubernetes.
$ sudo docker build -t plugin2sona .

$ kubectl create -f plugin2sona.yaml

```

### Out-of-clustering
```sh
$ sudo ./plugin2sona -kubeconfig=/etc/kubernetes/admin.conf
```
