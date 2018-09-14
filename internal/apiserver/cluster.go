package apiserver

import (
	"golang.org/x/net/context"

	"github.com/samsung-cnct/cluster-api-provider-ssh/cloud/ssh/providerconfig/v1alpha1"
	clusterapi "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"

	pb "github.com/samsung-cnct/cma-vmware/pkg/generated/api"
)

func (s *Server) CreateCluster(ctx context.Context, in *pb.CreateClusterMsg) (*pb.CreateClusterReply, error) {
	_ = clusterapi.Cluster{}
	_ = v1alpha1.SSHClusterProviderConfig{}

	return &pb.CreateClusterReply{
		Ok: true,
		Cluster: &pb.ClusterItem{
			Id:     "stub",
			Name:   "stub",
			Status: "Creating",
		},
	}, nil
}

func (s *Server) GetCluster(ctx context.Context, in *pb.GetClusterMsg) (*pb.GetClusterReply, error) {
	return &pb.GetClusterReply{
		Ok: true,
		Cluster: &pb.ClusterDetailItem{
			Id:         "stub",
			Name:       "stub",
			Status:     "stub",
			Kubeconfig: "xyz",
		},
	}, nil
}

func (s *Server) DeleteCluster(ctx context.Context, in *pb.DeleteClusterMsg) (*pb.DeleteClusterReply, error) {
	return &pb.DeleteClusterReply{Ok: true, Status: "stub"}, nil
}

func (s *Server) GetClusterList(ctx context.Context, in *pb.GetClusterListMsg) (reply *pb.GetClusterListReply, err error) {
	reply = &pb.GetClusterListReply{
		Ok: true,
	}
	return reply, nil
}
