package apiserver

import (
	"fmt"
	"golang.org/x/net/context"

	pb "github.com/samsung-cnct/cma-vmware/pkg/generated/api"
)

func (s *Server) CreateCluster(ctx context.Context, in *pb.CreateClusterMsg) (*pb.CreateClusterReply, error) {
	cluster := TranslateAPI(in)
	err := ApplyManifests(cluster)
	if err != nil {
		// TODO: Make this consistent with how the CMA does logging...
		fmt.Printf("ERROR: CreateCluster, name %v, err %v\n", in.Name, err)
		return &pb.CreateClusterReply{
			Ok: false,
			Cluster: &pb.ClusterItem{
				Id:     "stub",
				Name:   in.Name,
				Status: "CreateFailed",
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
	err := DeleteManifests(in.Name)
	if err != nil {
		return &pb.DeleteClusterReply{Ok: false, Status: "DeleteFailed"}, nil
	}

	return &pb.DeleteClusterReply{Ok: true, Status: "Deleted"}, nil
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
