# Tests

These tests demonstrate base end to end functionality of the VMware
helper of the Cluster Manager API.

# How to run

1.  populate the environment with required VARs (defaults are provided)
2.  execute `full-test.sh`

## Sequence of the tests

1.  create a client cluster via a parent CMA-VMware helper:
    `create-cluster.sh`
2.  get the kubeconfig for the client cluster from the parent cluster
    K8S API `get-kubeconfig()` in `full-test.sh`
3.  create a simple system in the client cluster (using nginx-ingress)
4.  verify the simple system functions
5.  tear down the client cluster: `delete-cluster.sh`
