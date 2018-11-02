#!/bin/bash

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

export CLUSTER_API=${CLUSTER_API:-cluster-manager-api.cnct.io}
export CLUSTER_API_PORT=${CLUSTER_API_PORT:-443}
export CLUSTER_NAME=${CLUSTER_NAME:-vmware-test-$(date +%s)}
export CLUSTER_API_NAMESPACE=${CLUSTER_API_NAMESPACE:-cma}

[[ -n $DEBUG ]] && set -o xtrace
set -o errexit
set -o nounset
set -o pipefail

get-kubeconfig(){
  kubectl get secrets --namespace $CLUSTER_API_NAMESPACE $CLUSTER_NAME -o yaml |
    grep kubernetes.kubeconfig |
    sed 's/  kubernetes.kubeconfig: //' |
    base64 --decode >kubeconfig.yaml
}

main() {
  ${__dir}/create-cluster.sh
  ${__dir}/get-cluster.sh

  get-kubeconfig

  ${__dir}/delete-cluster.sh
}

main
