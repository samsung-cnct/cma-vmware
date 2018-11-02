#!/bin/bash

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

export CLUSTER_API=${CLUSTER_API:-cluster-manager-api.cnct.io}
export CLUSTER_API_PORT=${CLUSTER_API_PORT:-443}
export CLUSTER_NAME=${CLUSTER_NAME:-vmware-test-$(date +%s)}

[[ -n $DEBUG ]] && set -o xtrace
set -o errexit
set -o nounset
set -o pipefail

main() {
  ${__dir}/create-cluster.sh
  ${__dir}/get-cluster.sh
  ${__dir}/delete-cluster.sh
}

main
