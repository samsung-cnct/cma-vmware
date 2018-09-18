package apiserver

const ClusterAPIProviderSSHTemplate = `
---
apiVersion: "cluster.k8s.io/v1alpha1"
kind: Cluster
metadata:
  name: {{ .Name }}
  namespace: {{ .Provider.Namespace }}
spec:
  clusterNetwork:
    services:
      cidrBlocks: ["10.96.0.0/12"]
    pods:
      cidrBlocks: ["10.244.0.0/16"]
    serviceDomain: "cluster.local"
  providerConfig:
    value:
      apiVersion: "sshproviderconfig/v1alpha1"
      kind: "SSHClusterProviderConfig"
{{ range .Provider.Machines }}
---
apiVersion: "cluster.k8s.io/v1alpha1"
kind: Machine
metadata:
  name: {{ .Name }}
spec:
  providerConfig:
    value:
      apiVersion: "sshproviderconfig/v1alpha1"
      kind: "SSHMachineProviderConfig"
      roles:
        - Master
        - Etcd
      sshConfig:
        username: {{ .Username }}
        host: {{ .Host }}
        port: {{ .Port }}
        secretName: cluster-private-key
  versions:
    kubelet: {{ $.Provider.K8SVersion }}
    controlPlane: {{ .ControlPlaneVersion }}
{{ end }}
---
apiVersion: v1
kind: Secret
type: Opaque
metadata:
  name: cluster-private-key
  namespace: default
data:
  private-key: {{ .Provider.PrivateKey }}
  pass-phrase: ""
`
