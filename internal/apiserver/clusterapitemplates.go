package apiserver

const SSHClusterTemplate = `
---
kind: Namespace
apiVersion: v1
metadata:
  name: {{ $.Name }}
---
apiVersion: "cluster.k8s.io/v1alpha1"
kind: Cluster
metadata:
  name: {{ $.Name }}
  namespace: {{ $.Name }}
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
{{ range $.ControlPlaneNodes }}
---
apiVersion: "cluster.k8s.io/v1alpha1"
kind: Machine
metadata:
  generateName: control-plane-
  namespace: {{ $.Name }}
  labels:
    controlPlaneVersion: {{ $.K8SVersion }}
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
    kubelet: {{ $.K8SVersion }}
    controlPlane: {{ $.K8SVersion }}
{{ end }}
{{ range $.WorkerNodes }}
---
apiVersion: "cluster.k8s.io/v1alpha1"
kind: Machine
metadata:
  generateName: worker-
  namespace: {{ $.Name }}
spec:
  providerConfig:
    value:
      apiVersion: "sshproviderconfig/v1alpha1"
      kind: "SSHMachineProviderConfig"
      roles:
        - Node
      sshConfig:
        username: {{ .Username }}
        host: {{ .Host }}
        port: {{ .Port }}
        secretName: cluster-private-key
  versions:
    kubelet: {{ $.K8SVersion }}
{{ end }}
---
apiVersion: v1
kind: Secret
type: Opaque
metadata:
  name: cluster-private-key
  namespace: {{ $.Name }}
data:
  private-key: {{ $.Base64PrivateKey }}
  pass-phrase: ""
`

const SSHMachineTemplate = `
{{ range $.ControlPlaneNodes }}
---
apiVersion: "cluster.k8s.io/v1alpha1"
kind: Machine
metadata:
  generateName: control-plane-
  namespace: {{ $.Name }}
  labels:
    controlPlaneVersion: {{ $.K8SVersion }}
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
    kubelet: {{ $.K8SVersion }}
    controlPlane: {{ $.K8SVersion }}
{{ end }}
{{ range $.WorkerNodes }}
---
apiVersion: "cluster.k8s.io/v1alpha1"
kind: Machine
metadata:
  generateName: worker-
  namespace: {{ $.Name }}
spec:
  providerConfig:
    value:
      apiVersion: "sshproviderconfig/v1alpha1"
      kind: "SSHMachineProviderConfig"
      roles:
        - Node
      sshConfig:
        username: {{ .Username }}
        host: {{ .Host }}
        port: {{ .Port }}
        secretName: cluster-private-key
  versions:
    kubelet: {{ $.K8SVersion }}
{{ end }}
`
