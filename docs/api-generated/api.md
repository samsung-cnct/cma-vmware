# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [api.proto](#api.proto)
    - [ClusterDetailItem](#cmavmware.ClusterDetailItem)
    - [ClusterItem](#cmavmware.ClusterItem)
    - [CreateClusterMsg](#cmavmware.CreateClusterMsg)
    - [CreateClusterProviderSpec](#cmavmware.CreateClusterProviderSpec)
    - [CreateClusterReply](#cmavmware.CreateClusterReply)
    - [CreateClusterVMWareSpec](#cmavmware.CreateClusterVMWareSpec)
    - [DeleteClusterMsg](#cmavmware.DeleteClusterMsg)
    - [DeleteClusterReply](#cmavmware.DeleteClusterReply)
    - [GetClusterListMsg](#cmavmware.GetClusterListMsg)
    - [GetClusterListReply](#cmavmware.GetClusterListReply)
    - [GetClusterMsg](#cmavmware.GetClusterMsg)
    - [GetClusterReply](#cmavmware.GetClusterReply)
    - [GetVersionMsg](#cmavmware.GetVersionMsg)
    - [GetVersionReply](#cmavmware.GetVersionReply)
    - [GetVersionReply.VersionInformation](#cmavmware.GetVersionReply.VersionInformation)
    - [MachineSpec](#cmavmware.MachineSpec)
  
  
  
    - [Cluster](#cmavmware.Cluster)
  

- [api.proto](#api.proto)
    - [ClusterDetailItem](#cmavmware.ClusterDetailItem)
    - [ClusterItem](#cmavmware.ClusterItem)
    - [CreateClusterMsg](#cmavmware.CreateClusterMsg)
    - [CreateClusterProviderSpec](#cmavmware.CreateClusterProviderSpec)
    - [CreateClusterReply](#cmavmware.CreateClusterReply)
    - [CreateClusterVMWareSpec](#cmavmware.CreateClusterVMWareSpec)
    - [DeleteClusterMsg](#cmavmware.DeleteClusterMsg)
    - [DeleteClusterReply](#cmavmware.DeleteClusterReply)
    - [GetClusterListMsg](#cmavmware.GetClusterListMsg)
    - [GetClusterListReply](#cmavmware.GetClusterListReply)
    - [GetClusterMsg](#cmavmware.GetClusterMsg)
    - [GetClusterReply](#cmavmware.GetClusterReply)
    - [GetVersionMsg](#cmavmware.GetVersionMsg)
    - [GetVersionReply](#cmavmware.GetVersionReply)
    - [GetVersionReply.VersionInformation](#cmavmware.GetVersionReply.VersionInformation)
    - [MachineSpec](#cmavmware.MachineSpec)
  
  
  
    - [Cluster](#cmavmware.Cluster)
  

- [Scalar Value Types](#scalar-value-types)



<a name="api.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api.proto



<a name="cmavmware.ClusterDetailItem"></a>

### ClusterDetailItem



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | ID of the cluster |
| name | [string](#string) |  | Name of the cluster |
| status | [string](#string) |  | What is the status of the cluster |
| kubeconfig | [string](#string) |  | What is the kubeconfig to connect to the cluster |






<a name="cmavmware.ClusterItem"></a>

### ClusterItem



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | ID of the cluster |
| name | [string](#string) |  | Name of the cluster |
| status | [string](#string) |  | What is the status of the cluster |






<a name="cmavmware.CreateClusterMsg"></a>

### CreateClusterMsg



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of the cluster to be provisioned |
| provider | [CreateClusterProviderSpec](#cmavmware.CreateClusterProviderSpec) |  | The provider specification |






<a name="cmavmware.CreateClusterProviderSpec"></a>

### CreateClusterProviderSpec



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | What is the provider - like vmware |
| k8s_version | [string](#string) |  | The version of Kubernetes for worker nodes. Control plane versions are determined by the MachineSpec. |
| vmware | [CreateClusterVMWareSpec](#cmavmware.CreateClusterVMWareSpec) |  | The VMWare specification |
| high_availability | [bool](#bool) |  | Whether or not the cluster is HA |
| network_fabric | [string](#string) |  | The fabric to be used |






<a name="cmavmware.CreateClusterReply"></a>

### CreateClusterReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ok | [bool](#bool) |  | Whether or not the cluster was provisioned by this request |
| cluster | [ClusterItem](#cmavmware.ClusterItem) |  | The details of the cluster request response |






<a name="cmavmware.CreateClusterVMWareSpec"></a>

### CreateClusterVMWareSpec



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| namespace | [string](#string) |  | This namespace along with the clustername with CreateClusterProviderSpec uniquely identify a managed cluster |
| machines | [MachineSpec](#cmavmware.MachineSpec) | repeated | Machines which comprise the cluster |






<a name="cmavmware.DeleteClusterMsg"></a>

### DeleteClusterMsg



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | What is the cluster&#39;s name to destroy |






<a name="cmavmware.DeleteClusterReply"></a>

### DeleteClusterReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ok | [bool](#bool) |  | Could the cluster be destroyed |
| status | [string](#string) |  | Status of the request |






<a name="cmavmware.GetClusterListMsg"></a>

### GetClusterListMsg







<a name="cmavmware.GetClusterListReply"></a>

### GetClusterListReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ok | [bool](#bool) |  | Is the cluster in the system |
| clusters | [ClusterItem](#cmavmware.ClusterItem) | repeated | List of clusters |






<a name="cmavmware.GetClusterMsg"></a>

### GetClusterMsg



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of the cluster to be looked up |






<a name="cmavmware.GetClusterReply"></a>

### GetClusterReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ok | [bool](#bool) |  | Is the cluster in the system |
| cluster | [ClusterDetailItem](#cmavmware.ClusterDetailItem) |  |  |






<a name="cmavmware.GetVersionMsg"></a>

### GetVersionMsg
Get version of API Server






<a name="cmavmware.GetVersionReply"></a>

### GetVersionReply
Reply for version request


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ok | [bool](#bool) |  | If operation was OK |
| version_information | [GetVersionReply.VersionInformation](#cmavmware.GetVersionReply.VersionInformation) |  | Version Information |






<a name="cmavmware.GetVersionReply.VersionInformation"></a>

### GetVersionReply.VersionInformation



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| git_version | [string](#string) |  | The tag on the git repository |
| git_commit | [string](#string) |  | The hash of the git commit |
| git_tree_state | [string](#string) |  | Whether or not the tree was clean when built |
| build_date | [string](#string) |  | Date of build |
| go_version | [string](#string) |  | Version of go used to compile |
| compiler | [string](#string) |  | Compiler used |
| platform | [string](#string) |  | Platform it was compiled for / running on |






<a name="cmavmware.MachineSpec"></a>

### MachineSpec
The credentials to use for creating the cluster


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| username | [string](#string) |  | The username for SSH access |
| host | [string](#string) |  | The host for SSH access |
| port | [string](#string) |  | The port for SSH access |
| private_key | [string](#string) |  | The private key for SSH. This must be preconfigured on the VMWare instances |
| control_plane_version | [string](#string) |  | The k8s version for the control plane. This node is only a master if this field is defined. |





 

 

 


<a name="cmavmware.Cluster"></a>

### Cluster


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateCluster | [CreateClusterMsg](#cmavmware.CreateClusterMsg) | [CreateClusterReply](#cmavmware.CreateClusterReply) | Will provision a cluster |
| GetCluster | [GetClusterMsg](#cmavmware.GetClusterMsg) | [GetClusterReply](#cmavmware.GetClusterReply) | Will retrieve the status of a cluster and its kubeconfig for connectivity |
| DeleteCluster | [DeleteClusterMsg](#cmavmware.DeleteClusterMsg) | [DeleteClusterReply](#cmavmware.DeleteClusterReply) | Will delete a cluster |
| GetClusterList | [GetClusterListMsg](#cmavmware.GetClusterListMsg) | [GetClusterListReply](#cmavmware.GetClusterListReply) | Will retrieve a list of clusters |
| GetVersionInformation | [GetVersionMsg](#cmavmware.GetVersionMsg) | [GetVersionReply](#cmavmware.GetVersionReply) | Will return version information about api server |

 



<a name="api.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api.proto



<a name="cmavmware.ClusterDetailItem"></a>

### ClusterDetailItem



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | ID of the cluster |
| name | [string](#string) |  | Name of the cluster |
| status | [string](#string) |  | What is the status of the cluster |
| kubeconfig | [string](#string) |  | What is the kubeconfig to connect to the cluster |






<a name="cmavmware.ClusterItem"></a>

### ClusterItem



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | ID of the cluster |
| name | [string](#string) |  | Name of the cluster |
| status | [string](#string) |  | What is the status of the cluster |






<a name="cmavmware.CreateClusterMsg"></a>

### CreateClusterMsg



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of the cluster to be provisioned |
| provider | [CreateClusterProviderSpec](#cmavmware.CreateClusterProviderSpec) |  | The provider specification |






<a name="cmavmware.CreateClusterProviderSpec"></a>

### CreateClusterProviderSpec



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | What is the provider - like vmware |
| k8s_version | [string](#string) |  | The version of Kubernetes for worker nodes. Control plane versions are determined by the MachineSpec. |
| vmware | [CreateClusterVMWareSpec](#cmavmware.CreateClusterVMWareSpec) |  | The VMWare specification |
| high_availability | [bool](#bool) |  | Whether or not the cluster is HA |
| network_fabric | [string](#string) |  | The fabric to be used |






<a name="cmavmware.CreateClusterReply"></a>

### CreateClusterReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ok | [bool](#bool) |  | Whether or not the cluster was provisioned by this request |
| cluster | [ClusterItem](#cmavmware.ClusterItem) |  | The details of the cluster request response |






<a name="cmavmware.CreateClusterVMWareSpec"></a>

### CreateClusterVMWareSpec



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| namespace | [string](#string) |  | This namespace along with the clustername with CreateClusterProviderSpec uniquely identify a managed cluster |
| machines | [MachineSpec](#cmavmware.MachineSpec) | repeated | Machines which comprise the cluster |






<a name="cmavmware.DeleteClusterMsg"></a>

### DeleteClusterMsg



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | What is the cluster&#39;s name to destroy |






<a name="cmavmware.DeleteClusterReply"></a>

### DeleteClusterReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ok | [bool](#bool) |  | Could the cluster be destroyed |
| status | [string](#string) |  | Status of the request |






<a name="cmavmware.GetClusterListMsg"></a>

### GetClusterListMsg







<a name="cmavmware.GetClusterListReply"></a>

### GetClusterListReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ok | [bool](#bool) |  | Is the cluster in the system |
| clusters | [ClusterItem](#cmavmware.ClusterItem) | repeated | List of clusters |






<a name="cmavmware.GetClusterMsg"></a>

### GetClusterMsg



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of the cluster to be looked up |






<a name="cmavmware.GetClusterReply"></a>

### GetClusterReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ok | [bool](#bool) |  | Is the cluster in the system |
| cluster | [ClusterDetailItem](#cmavmware.ClusterDetailItem) |  |  |






<a name="cmavmware.GetVersionMsg"></a>

### GetVersionMsg
Get version of API Server






<a name="cmavmware.GetVersionReply"></a>

### GetVersionReply
Reply for version request


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ok | [bool](#bool) |  | If operation was OK |
| version_information | [GetVersionReply.VersionInformation](#cmavmware.GetVersionReply.VersionInformation) |  | Version Information |






<a name="cmavmware.GetVersionReply.VersionInformation"></a>

### GetVersionReply.VersionInformation



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| git_version | [string](#string) |  | The tag on the git repository |
| git_commit | [string](#string) |  | The hash of the git commit |
| git_tree_state | [string](#string) |  | Whether or not the tree was clean when built |
| build_date | [string](#string) |  | Date of build |
| go_version | [string](#string) |  | Version of go used to compile |
| compiler | [string](#string) |  | Compiler used |
| platform | [string](#string) |  | Platform it was compiled for / running on |






<a name="cmavmware.MachineSpec"></a>

### MachineSpec
The credentials to use for creating the cluster


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| username | [string](#string) |  | The username for SSH access |
| host | [string](#string) |  | The host for SSH access |
| port | [string](#string) |  | The port for SSH access |
| private_key | [string](#string) |  | The private key for SSH. This must be preconfigured on the VMWare instances |
| control_plane_version | [string](#string) |  | The k8s version for the control plane. This node is only a master if this field is defined. |





 

 

 


<a name="cmavmware.Cluster"></a>

### Cluster


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateCluster | [CreateClusterMsg](#cmavmware.CreateClusterMsg) | [CreateClusterReply](#cmavmware.CreateClusterReply) | Will provision a cluster |
| GetCluster | [GetClusterMsg](#cmavmware.GetClusterMsg) | [GetClusterReply](#cmavmware.GetClusterReply) | Will retrieve the status of a cluster and its kubeconfig for connectivity |
| DeleteCluster | [DeleteClusterMsg](#cmavmware.DeleteClusterMsg) | [DeleteClusterReply](#cmavmware.DeleteClusterReply) | Will delete a cluster |
| GetClusterList | [GetClusterListMsg](#cmavmware.GetClusterListMsg) | [GetClusterListReply](#cmavmware.GetClusterListReply) | Will retrieve a list of clusters |
| GetVersionInformation | [GetVersionMsg](#cmavmware.GetVersionMsg) | [GetVersionReply](#cmavmware.GetVersionReply) | Will return version information about api server |

 



## Scalar Value Types

| .proto Type | Notes | C++ Type | Java Type | Python Type |
| ----------- | ----- | -------- | --------- | ----------- |
| <a name="double" /> double |  | double | double | float |
| <a name="float" /> float |  | float | float | float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long |
| <a name="bool" /> bool |  | bool | boolean | boolean |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str |

