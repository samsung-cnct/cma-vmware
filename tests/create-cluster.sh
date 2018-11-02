#!/bin/bash

CLUSTER_API=${CLUSTER_API:-cluster-manager-api.cnct.io}
CLUSTER_API_PORT=${CLUSTER_API_PORT:-443}
CLUSTER_NAME=${CLUSTER_NAME:-vmware-test-$(date +%s)}
MASTER_IP=${MASTER_IP:-182.195.81.132}
WORKER_1_IP=${WORKER_IP:-182.195.81.137}
K8S_VERSION=${K8S_VERSION:-1.10.6}
NODE_USER=${NODE_USER:-root}

[[ -n $DEBUG ]] && set -o xtrace
set -o errexit
set -o nounset
set -o pipefail

DATA=$(
  cat <<JSON
{
  "name": "${CLUSTER_NAME}",
  "provider": {
    "name": "vmware",
    "k8s_version": "${K8S_VERSION}",
    "vmware": {
      "control_plane_nodes": [
        {
          "username": "${NODE_USER}",
          "host": "${MASTER_IP}",
          "port": 22,
          "password": "${NODE_PASSWORD}",
          "labels": [
            {
              "name": "master",
              "value": "1"
            }
          ]
        }
      ],
      "worker_nodes": [
        {
          "username": "${NODE_USER}",
          "host": "${WORKER_1_IP}",
          "port": 22,
          "password": "${NODE_PASSWORD}",
          "labels": [
            {
              "name": "worker",
              "value": "1"
            }
          ]
        }
      ],
      "api_endpoint": "172.10.10.1:443"
    },
    "high_availability": false,
    "network_fabric": "canal"
  },
  "callback": {
    "url": "unused",
    "request_id": "1234"
  }
}
JSON
)

main() {
  curl -X POST \
    "https://${CLUSTER_API}:${CLUSTER_API_PORT}/api/v1/cluster" \
    -H 'Cache-Control: no-cache' \
    -H 'Content-Type: application/json' \
    -d "${DATA}" \
    -iks
}

main
