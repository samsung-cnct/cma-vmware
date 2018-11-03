#!/bin/bash

CLUSTER_API=${CLUSTER_API:-cluster-manager-api.cnct.io}
CLUSTER_API_PORT=${CLUSTER_API_PORT:-443}
CLUSTER_NAME=${CLUSTER_NAME:-vmware-test-$(date +%s)}

[[ -n $DEBUG ]] && set -o xtrace
set -o errexit
set -o nounset
set -o pipefail

main() {
  curl -X GET \
    "https://${CLUSTER_API}:${CLUSTER_API_PORT}/api/v1/cluster?provider=vmware&name=${CLUSTER_NAME}" \
    -H 'Cache-Control: no-cache' \
    -H 'Content-Type: application/json' \
    -iks
}

main
