package apiserver

import (
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
				Status: pb.ClusterStatus_ERROR,
			},
		}, nil
	}

	return &pb.CreateClusterReply{
		Ok: true,
		Cluster: &pb.ClusterItem{
			Id:     "stub",
			Name:   in.Name,
			Status: pb.ClusterStatus_PROVISIONING,
		},
	}, nil
}

func (s *Server) GetCluster(ctx context.Context, in *pb.GetClusterMsg) (*pb.GetClusterReply, error) {
	kubeconfigBytes, err := GetKubeConfig(in.Name)
	if err != nil {
		fmt.Printf("INFO: GetCluster unable to getKubeconfig secret for cluster %s\n", in.Name)
	}
	clusterExists, err := ClusterExists(in.Name)
	if !clusterExists {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	clusterStatus, err := GetSSHClusterStatus(in, kubeconfigBytes)
	if err != nil {
		fmt.Printf("ERROR: GetCluster, %v, err %v\n", in.Name, err)
		return &pb.GetClusterReply{
			Ok: true,
			Cluster: &pb.ClusterDetailItem{
				// TODO: Standardize on the nil UUID for bare metal?
				// 00000000-0000-0000-0000-000000 000000?
				Id:         "stub",
				Name:       in.Name,
				Status:     pb.ClusterStatus_ERROR,
				Kubeconfig: "",
			},
		}, nil
	}

	return &pb.GetClusterReply{
		Ok: true,
		Cluster: &pb.ClusterDetailItem{
			Id:         "stub",
			Name:       in.Name,
			Status:     clusterStatus,
			Kubeconfig: string(kubeconfigBytes),
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
			Status: pb.ClusterStatus_RUNNING,
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
	// the machine-setup ConfigMap defined here: https://goo.gl/Wi81Z9
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
	kubeconfigBytes, err := GetKubeConfig(in.Name)
	if err != nil {
		fmt.Printf("INFO: UpgradeCluster unable to getKubeconfig secret for cluster %s\n", in.Name)
		return nil, status.Error(codes.NotFound, err.Error())
	}
	err = UpgradeSSHCluster(in.Name, in.Version, kubeconfigBytes)
	if err != nil {
		fmt.Printf("ERROR: UpgradeCluster, %v, err %v\n", in.Name, err)
		return &pb.UpgradeClusterReply{
			Ok: true,
		}, err
	}
	return &pb.UpgradeClusterReply{
		Ok: true,
	}, nil
}
