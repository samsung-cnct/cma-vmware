# Introduction

`cma-vmware` is an implementation of the [`cluster-manager-api`](
https://github.com/samsung-cnct/cluster-manager-api). `cluster-manager-api`
is a facade in front of Kubernetes (k8s), Helm, and other tools used to build
and manage k8s clusters in both on-prem and cloud enviroments (e.g. AWS, Azure,
etc.) `cma-vmware` is an implementation for on-prem environments and is based
on [`cluster-api-provider-ssh`](
https://github.com/samsung-cnct/cluster-api-provider-ssh/).

# Quickstart

This guide assumes you have already deployed `cluster-api-provider-ssh`
according to its instructions.

First install `tiller` and `cert-manager` if they have not been already:

```bash
export KUBECONFIG=<path/to/kubeconfig>
kubectl create clusterrolebinding tiller-cluster-rule --clusterrole=cluster-admin --serviceaccount=kube-system:tiller
helm init --service-account tiller
helm install --tiller-namespace=kube-system --name cert-manager --namespace kube-system stable/cert-manager
```

Then install `cma-vmware`:

```
helm install --tiller-namespace=kube-system --name cma-vmware --namespace cma-vmware
```

`cma-vmware` contains a built-in GUI which can be useful in determining the 
structure of the API and submitting requests. See the helm values file and
documentation for your k8s environment to determine how to access it.

# Development

## To download, build, and run `cma-vmware` locally

On OSX:

```bash
go get github.com/samsung-cnct/cma-vmware
cd $GOPATH/src/github.com/samsung-cnct/cma-vmware
make -f build/Makefile cmavmw-bin-darwin
./cma-vmware
```

On Linux, replace the `make` command above with this one:

```
make -f build/Makefile cmavmw-bin-linux-amd64
```

## To submit requests for testing purposes

Either use the GUI:

```
open http://127.0.0.1:9020/swagger-ui/#/
``` 

Use `curl`:

```
curl 127.0.0.1:9020/api/v1/cluster -X POST -d @dewdrops/cma.json
curl 127.0.0.1:9020/api/v1/cluster?name=dewdrops -X DELETE
```

Or use another http client.

## To change the API and regenerate golang bindings

```
vi api/api.proto
make -f build/Makefile generators
```
