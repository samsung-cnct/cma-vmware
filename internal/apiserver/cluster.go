package apiserver

import (
	"fmt"
	"golang.org/x/net/context"

	"github.com/samsung-cnct/cluster-api-provider-ssh/cloud/ssh/providerconfig/v1alpha1"
	clusterapi "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"

	pb "github.com/samsung-cnct/cma-vmware/pkg/generated/api"
)

func (s *Server) CreateCluster(ctx context.Context, in *pb.CreateClusterMsg) (*pb.CreateClusterReply, error) {
	manifests, err := GetManifests(in)
	if err != nil {
		// TODO: Make this consistent with how the CMA does logging...
		fmt.Printf("ERROR: CreateCluster, GetManifests, name %v, err %v", in.Name, err)
		return &pb.CreateClusterReply{
			Ok: false,
			Cluster: &pb.ClusterItem{
				Id:     "stub",
				Name:   in.Name,
				Status: "Failed",
			},
		}, nil
	}

	err = ApplyManifests(manifests)
	if err != nil {
		// TODO: Make this consistent with how the CMA does logging...
		fmt.Printf("ERROR: CreateCluster, name %v, err %v", in.Name, err)
		return &pb.CreateClusterReply{
			Ok: false,
			Cluster: &pb.ClusterItem{
				Id:     "stub",
				Name:   in.Name,
				Status: "PossiblePartialFailure",
			},
		}, nil
	}

	return &pb.CreateClusterReply{
		Ok: true,
		Cluster: &pb.ClusterItem{
			Id:     "stub",
			Name:   in.Name,
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
	_ = clusterapi.Cluster{}
	_ = v1alpha1.SSHClusterProviderConfig{}

	return &pb.DeleteClusterReply{Ok: true, Status: "stub"}, nil
}

func (s *Server) GetClusterList(ctx context.Context, in *pb.GetClusterListMsg) (reply *pb.GetClusterListReply, err error) {
	reply = &pb.GetClusterListReply{
		Ok: true,
	}
	return reply, nil
}

func (s *Server) AdjustClusterNodes(ctx context.Context, in *pb.AdjustClusterMsg) (*pb.AdjustClusterReply, error) {
	return &pb.AdjustClusterReply{}, fmt.Errorf("adjust cluster nodes not implemented yet")
}

func (s *Server) GetUpgradeClusterInformation(ctx context.Context, in *pb.GetUpgradeClusterInformationMsg) (*pb.GetUpgradeClusterInformationReply, error) {
	return &pb.GetUpgradeClusterInformationReply{}, fmt.Errorf("get cluster upgrade information not implemented yet")
}

func (s *Server) UpgradeCluster(ctx context.Context, in *pb.UpgradeClusterMsg) (*pb.UpgradeClusterReply, error) {
	return &pb.UpgradeClusterReply{}, fmt.Errorf("cluster upgrade not implemented yet")
}
