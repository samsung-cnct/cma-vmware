package apiserver

import (
	"fmt"
	"golang.org/x/net/context"

	pb "github.com/samsung-cnct/cma-vmware/pkg/generated/api"
)

func (s *Server) CreateCluster(ctx context.Context, in *pb.CreateClusterMsg) (*pb.CreateClusterReply, error) {
	err := CreateSSHCluster(in)
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
	kubeconfig, err := GetKubeConfig(in.Name)
	if err != nil {
		return &pb.GetClusterReply{
			Ok: true,
			Cluster: &pb.ClusterDetailItem{
				Id:         "stub",
				Name:       in.Name,
				Status:     "GetFailed",
				Kubeconfig: "",
			},
		}, nil
	}

	return &pb.GetClusterReply{
		Ok: true,
		Cluster: &pb.ClusterDetailItem{
			Id:         "stub",
			Name:       in.Name,
			Status:     "Yes",
			Kubeconfig: kubeconfig,
		},
	}, nil
}

func (s *Server) DeleteCluster(ctx context.Context, in *pb.DeleteClusterMsg) (*pb.DeleteClusterReply, error) {
	err := DeleteSSHCluster(in.Name)
	if err != nil {
		return &pb.DeleteClusterReply{Ok: false, Status: "DeleteFailed"}, nil
	}

	return &pb.DeleteClusterReply{Ok: true, Status: "Deleted"}, nil
}

func (s *Server) GetClusterList(ctx context.Context, in *pb.GetClusterListMsg) (reply *pb.GetClusterListReply, err error) {
	clusterNames, err := ListSSHClusters()
	if err != nil {
		return &pb.GetClusterListReply{
			Ok: false,
		}, err
	}

	var clusters []*pb.ClusterItem
	for _, name := range clusterNames {
		clusters = append(clusters, &pb.ClusterItem{
			Id:     "stub",
			Name:   name,
			Status: "Yes",
		})
	}

	return &pb.GetClusterListReply{
		Ok:       true,
		Clusters: clusters,
	}, nil
}

func (s *Server) AdjustClusterNodes(ctx context.Context, in *pb.AdjustClusterMsg) (*pb.AdjustClusterReply, error) {
	err := AdjustSSHCluster(in)
	if err != nil {
		return &pb.AdjustClusterReply{Ok: false}, nil
	}

	return &pb.AdjustClusterReply{Ok: true}, nil
}

func (s *Server) GetUpgradeClusterInformation(ctx context.Context, in *pb.GetUpgradeClusterInformationMsg) (*pb.GetUpgradeClusterInformationReply, error) {
	// TODO: Do not hard code this list. Before adding versions, update
	// the machine-setup ConfigMap in defined here: https://goo.gl/Wi81Z9
	return &pb.GetUpgradeClusterInformationReply{
		Versions: []string{
			"1.10.4",
			"1.10.6",
			"1.11.2",
			"1.11.3",
		},
	}, nil
}

func (s *Server) UpgradeCluster(ctx context.Context, in *pb.UpgradeClusterMsg) (*pb.UpgradeClusterReply, error) {
	err := UpgradeSSHCluster(in.Name, in.Version)
	if err != nil {
		return &pb.UpgradeClusterReply{
			Ok: true,
		}, err
	}
	return &pb.UpgradeClusterReply{
		Ok: true,
	}, nil
}
