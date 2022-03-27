package api

import (
	"context"
	"fmt"
	"github.com/multycloud/multy/api/proto"
	"github.com/multycloud/multy/api/proto/common"
	"github.com/multycloud/multy/api/proto/resources"
	"github.com/multycloud/multy/api/services/database"
	"github.com/multycloud/multy/api/services/kubernetes_cluster"
	"github.com/multycloud/multy/api/services/kubernetes_node_pool"
	"github.com/multycloud/multy/api/services/lambda"
	"github.com/multycloud/multy/api/services/network_interface"
	"github.com/multycloud/multy/api/services/network_security_group"
	"github.com/multycloud/multy/api/services/object_storage"
	"github.com/multycloud/multy/api/services/object_storage_object"
	"github.com/multycloud/multy/api/services/public_ip"
	"github.com/multycloud/multy/api/services/route_table"
	"github.com/multycloud/multy/api/services/route_table_association"
	"github.com/multycloud/multy/api/services/subnet"
	"github.com/multycloud/multy/api/services/vault"
	"github.com/multycloud/multy/api/services/virtual_machine"
	"github.com/multycloud/multy/api/services/virtual_network"
	"github.com/multycloud/multy/db"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
)

type Server struct {
	proto.UnimplementedMultyResourceServiceServer
	virtual_network.VnService
	subnet.SubnetService
	network_interface.NetworkInterfaceService
	route_table.RouteTableService
	route_table_association.RouteTableAssociationService
	network_security_group.NetworkSecurityGroupService
	database.DatabaseService
	object_storage.ObjectStorageService
	object_storage_object.ObjectStorageObjectService
	public_ip.PublicIpService
	kubernetes_cluster.KubernetesClusterService
	kubernetes_node_pool.KubernetesNodePoolService
	lambda.LambdaService
	vault.VaultService
	vault.VaultAccessPolicyService
	vault.VaultSecretService
	virtual_machine.VirtualMachineService
}

func RunServer(ctx context.Context, port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	endpoint := "api.multy.dev"
	certFile := fmt.Sprintf("/etc/letsencrypt/live/%s/fullchain.pem", endpoint)
	keyFile := fmt.Sprintf("/etc/letsencrypt/live/%s/privkey.pem", endpoint)
	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	var s *grpc.Server
	if err != nil {
		log.Printf("unable to read certificate (%s), running in insecure mode", err.Error())
		s = grpc.NewServer()
	} else {
		s = grpc.NewServer(grpc.Creds(creds))
	}

	go func() {
		<-ctx.Done()
		s.GracefulStop()
		_ = lis.Close()
	}()
	d, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("failed to load db: %v", err)
	}
	server := Server{
		proto.UnimplementedMultyResourceServiceServer{},
		virtual_network.NewVnService(d),
		subnet.NewSubnetService(d),
		network_interface.NewNetworkInterfaceService(d),
		route_table.NewRouteTableService(d),
		route_table_association.NewRouteTableAssociationService(d),
		network_security_group.NewNetworkSecurityGroupService(d),
		database.NewDatabaseService(d),
		object_storage.NewObjectStorageService(d),
		object_storage_object.NewObjectStorageObjectService(d),
		public_ip.NewPublicIpService(d),
		kubernetes_cluster.NewKubernetesClusterService(d),
		kubernetes_node_pool.NewKubernetesNodePoolService(d),
		lambda.NewLambdaService(d),
		vault.NewVaultService(d),
		vault.NewVaultAccessPolicyService(d),
		vault.NewVaultSecretService(d),
		virtual_machine.NewVirtualMachineService(d),
	}
	proto.RegisterMultyResourceServiceServer(s, &server)
	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *Server) CreateVirtualNetwork(ctx context.Context, in *resources.CreateVirtualNetworkRequest) (*resources.VirtualNetworkResource, error) {
	return s.VnService.Service.Create(ctx, in)
}
func (s *Server) ReadVirtualNetwork(ctx context.Context, in *resources.ReadVirtualNetworkRequest) (*resources.VirtualNetworkResource, error) {
	return s.VnService.Service.Read(ctx, in)
}
func (s *Server) UpdateVirtualNetwork(ctx context.Context, in *resources.UpdateVirtualNetworkRequest) (*resources.VirtualNetworkResource, error) {
	return s.VnService.Service.Update(ctx, in)
}
func (s *Server) DeleteVirtualNetwork(ctx context.Context, in *resources.DeleteVirtualNetworkRequest) (*common.Empty, error) {
	return s.VnService.Service.Delete(ctx, in)
}

func (s *Server) CreateSubnet(ctx context.Context, in *resources.CreateSubnetRequest) (*resources.SubnetResource, error) {
	return s.SubnetService.Service.Create(ctx, in)
}
func (s *Server) ReadSubnet(ctx context.Context, in *resources.ReadSubnetRequest) (*resources.SubnetResource, error) {
	return s.SubnetService.Service.Read(ctx, in)
}
func (s *Server) UpdateSubnet(ctx context.Context, in *resources.UpdateSubnetRequest) (*resources.SubnetResource, error) {
	return s.SubnetService.Service.Update(ctx, in)
}
func (s *Server) DeleteSubnet(ctx context.Context, in *resources.DeleteSubnetRequest) (*common.Empty, error) {
	return s.SubnetService.Service.Delete(ctx, in)
}

func (s *Server) CreateNetworkInterface(ctx context.Context, in *resources.CreateNetworkInterfaceRequest) (*resources.NetworkInterfaceResource, error) {
	return s.NetworkInterfaceService.Service.Create(ctx, in)
}
func (s *Server) ReadNetworkInterface(ctx context.Context, in *resources.ReadNetworkInterfaceRequest) (*resources.NetworkInterfaceResource, error) {
	return s.NetworkInterfaceService.Service.Read(ctx, in)
}
func (s *Server) UpdateNetworkInterface(ctx context.Context, in *resources.UpdateNetworkInterfaceRequest) (*resources.NetworkInterfaceResource, error) {
	return s.NetworkInterfaceService.Service.Update(ctx, in)
}
func (s *Server) DeleteNetworkInterface(ctx context.Context, in *resources.DeleteNetworkInterfaceRequest) (*common.Empty, error) {
	return s.NetworkInterfaceService.Service.Delete(ctx, in)
}

func (s *Server) CreateRouteTable(ctx context.Context, in *resources.CreateRouteTableRequest) (*resources.RouteTableResource, error) {
	return s.RouteTableService.Service.Create(ctx, in)
}
func (s *Server) ReadRouteTable(ctx context.Context, in *resources.ReadRouteTableRequest) (*resources.RouteTableResource, error) {
	return s.RouteTableService.Service.Read(ctx, in)
}
func (s *Server) UpdateRouteTable(ctx context.Context, in *resources.UpdateRouteTableRequest) (*resources.RouteTableResource, error) {
	return s.RouteTableService.Service.Update(ctx, in)
}
func (s *Server) DeleteRouteTable(ctx context.Context, in *resources.DeleteRouteTableRequest) (*common.Empty, error) {
	return s.RouteTableService.Service.Delete(ctx, in)
}

func (s *Server) CreateRouteTableAssociation(ctx context.Context, in *resources.CreateRouteTableAssociationRequest) (*resources.RouteTableAssociationResource, error) {
	return s.RouteTableAssociationService.Service.Create(ctx, in)
}
func (s *Server) ReadRouteTableAssociation(ctx context.Context, in *resources.ReadRouteTableAssociationRequest) (*resources.RouteTableAssociationResource, error) {
	return s.RouteTableAssociationService.Service.Read(ctx, in)
}
func (s *Server) UpdateRouteTableAssociation(ctx context.Context, in *resources.UpdateRouteTableAssociationRequest) (*resources.RouteTableAssociationResource, error) {
	return s.RouteTableAssociationService.Service.Update(ctx, in)
}
func (s *Server) DeleteRouteTableAssociation(ctx context.Context, in *resources.DeleteRouteTableAssociationRequest) (*common.Empty, error) {
	return s.RouteTableService.Service.Delete(ctx, in)
}

func (s *Server) CreateNetworkSecurityGroup(ctx context.Context, in *resources.CreateNetworkSecurityGroupRequest) (*resources.NetworkSecurityGroupResource, error) {
	return s.NetworkSecurityGroupService.Service.Create(ctx, in)
}
func (s *Server) ReadNetworkSecurityGroup(ctx context.Context, in *resources.ReadNetworkSecurityGroupRequest) (*resources.NetworkSecurityGroupResource, error) {
	return s.NetworkSecurityGroupService.Service.Read(ctx, in)
}
func (s *Server) UpdateNetworkSecurityGroup(ctx context.Context, in *resources.UpdateNetworkSecurityGroupRequest) (*resources.NetworkSecurityGroupResource, error) {
	return s.NetworkSecurityGroupService.Service.Update(ctx, in)
}
func (s *Server) DeleteNetworkSecurityGroup(ctx context.Context, in *resources.DeleteNetworkSecurityGroupRequest) (*common.Empty, error) {
	return s.NetworkSecurityGroupService.Service.Delete(ctx, in)
}

func (s *Server) CreateDatabase(ctx context.Context, in *resources.CreateDatabaseRequest) (*resources.DatabaseResource, error) {
	return s.DatabaseService.Service.Create(ctx, in)
}
func (s *Server) ReadDatabase(ctx context.Context, in *resources.ReadDatabaseRequest) (*resources.DatabaseResource, error) {
	return s.DatabaseService.Service.Read(ctx, in)
}
func (s *Server) UpdateDatabase(ctx context.Context, in *resources.UpdateDatabaseRequest) (*resources.DatabaseResource, error) {
	return s.DatabaseService.Service.Update(ctx, in)
}
func (s *Server) DeleteDatabase(ctx context.Context, in *resources.DeleteDatabaseRequest) (*common.Empty, error) {
	return s.DatabaseService.Service.Delete(ctx, in)
}

func (s *Server) CreateObjectStorage(ctx context.Context, in *resources.CreateObjectStorageRequest) (*resources.ObjectStorageResource, error) {
	return s.ObjectStorageService.Service.Create(ctx, in)
}
func (s *Server) ReadObjectStorage(ctx context.Context, in *resources.ReadObjectStorageRequest) (*resources.ObjectStorageResource, error) {
	return s.ObjectStorageService.Service.Read(ctx, in)
}
func (s *Server) UpdateObjectStorage(ctx context.Context, in *resources.UpdateObjectStorageRequest) (*resources.ObjectStorageResource, error) {
	return s.ObjectStorageService.Service.Update(ctx, in)
}
func (s *Server) DeleteObjectStorage(ctx context.Context, in *resources.DeleteObjectStorageRequest) (*common.Empty, error) {
	return s.ObjectStorageService.Service.Delete(ctx, in)
}

func (s *Server) CreateObjectStorageObject(ctx context.Context, in *resources.CreateObjectStorageObjectRequest) (*resources.ObjectStorageObjectResource, error) {
	return s.ObjectStorageObjectService.Service.Create(ctx, in)
}
func (s *Server) ReadObjectStorageObject(ctx context.Context, in *resources.ReadObjectStorageObjectRequest) (*resources.ObjectStorageObjectResource, error) {
	return s.ObjectStorageObjectService.Service.Read(ctx, in)
}
func (s *Server) UpdateObjectStorageObject(ctx context.Context, in *resources.UpdateObjectStorageObjectRequest) (*resources.ObjectStorageObjectResource, error) {
	return s.ObjectStorageObjectService.Service.Update(ctx, in)
}
func (s *Server) DeleteObjectStorageObject(ctx context.Context, in *resources.DeleteObjectStorageObjectRequest) (*common.Empty, error) {
	return s.ObjectStorageObjectService.Service.Delete(ctx, in)
}

func (s *Server) CreatePublicIp(ctx context.Context, in *resources.CreatePublicIpRequest) (*resources.PublicIpResource, error) {
	return s.PublicIpService.Service.Create(ctx, in)
}
func (s *Server) ReadPublicIp(ctx context.Context, in *resources.ReadPublicIpRequest) (*resources.PublicIpResource, error) {
	return s.PublicIpService.Service.Read(ctx, in)
}
func (s *Server) UpdatePublicIp(ctx context.Context, in *resources.UpdatePublicIpRequest) (*resources.PublicIpResource, error) {
	return s.PublicIpService.Service.Update(ctx, in)
}
func (s *Server) DeletePublicIp(ctx context.Context, in *resources.DeletePublicIpRequest) (*common.Empty, error) {
	return s.PublicIpService.Service.Delete(ctx, in)
}

func (s *Server) CreateKubernetesCluster(ctx context.Context, in *resources.CreateKubernetesClusterRequest) (*resources.KubernetesClusterResource, error) {
	return s.KubernetesClusterService.Service.Create(ctx, in)
}
func (s *Server) ReadKubernetesCluster(ctx context.Context, in *resources.ReadKubernetesClusterRequest) (*resources.KubernetesClusterResource, error) {
	return s.KubernetesClusterService.Service.Read(ctx, in)
}
func (s *Server) UpdateKubernetesCluster(ctx context.Context, in *resources.UpdateKubernetesClusterRequest) (*resources.KubernetesClusterResource, error) {
	return s.KubernetesClusterService.Service.Update(ctx, in)
}
func (s *Server) DeleteKubernetesCluster(ctx context.Context, in *resources.DeleteKubernetesClusterRequest) (*common.Empty, error) {
	return s.KubernetesClusterService.Service.Delete(ctx, in)
}

func (s *Server) CreateKubernetesNodePool(ctx context.Context, in *resources.CreateKubernetesNodePoolRequest) (*resources.KubernetesNodePoolResource, error) {
	return s.KubernetesNodePoolService.Service.Create(ctx, in)
}
func (s *Server) ReadKubernetesNodePool(ctx context.Context, in *resources.ReadKubernetesNodePoolRequest) (*resources.KubernetesNodePoolResource, error) {
	return s.KubernetesNodePoolService.Service.Read(ctx, in)
}
func (s *Server) UpdateKubernetesNodePool(ctx context.Context, in *resources.UpdateKubernetesNodePoolRequest) (*resources.KubernetesNodePoolResource, error) {
	return s.KubernetesNodePoolService.Service.Update(ctx, in)
}
func (s *Server) DeleteKubernetesNodePool(ctx context.Context, in *resources.DeleteKubernetesNodePoolRequest) (*common.Empty, error) {
	return s.KubernetesNodePoolService.Service.Delete(ctx, in)
}

func (s *Server) CreateLambda(ctx context.Context, in *resources.CreateLambdaRequest) (*resources.LambdaResource, error) {
	return s.LambdaService.Service.Create(ctx, in)
}
func (s *Server) ReadLambda(ctx context.Context, in *resources.ReadLambdaRequest) (*resources.LambdaResource, error) {
	return s.LambdaService.Service.Read(ctx, in)
}
func (s *Server) UpdateLambda(ctx context.Context, in *resources.UpdateLambdaRequest) (*resources.LambdaResource, error) {
	return s.LambdaService.Service.Update(ctx, in)
}
func (s *Server) DeleteLambda(ctx context.Context, in *resources.DeleteLambdaRequest) (*common.Empty, error) {
	return s.LambdaService.Service.Delete(ctx, in)
}

func (s *Server) CreateVault(ctx context.Context, in *resources.CreateVaultRequest) (*resources.VaultResource, error) {
	return s.VaultService.Service.Create(ctx, in)
}
func (s *Server) ReadVault(ctx context.Context, in *resources.ReadVaultRequest) (*resources.VaultResource, error) {
	return s.VaultService.Service.Read(ctx, in)
}
func (s *Server) UpdateVault(ctx context.Context, in *resources.UpdateVaultRequest) (*resources.VaultResource, error) {
	return s.VaultService.Service.Update(ctx, in)
}
func (s *Server) DeleteVault(ctx context.Context, in *resources.DeleteVaultRequest) (*common.Empty, error) {
	return s.VaultService.Service.Delete(ctx, in)
}

func (s *Server) CreateVaultSecret(ctx context.Context, in *resources.CreateVaultSecretRequest) (*resources.VaultSecretResource, error) {
	return s.VaultSecretService.Service.Create(ctx, in)
}
func (s *Server) ReadVaultSecret(ctx context.Context, in *resources.ReadVaultSecretRequest) (*resources.VaultSecretResource, error) {
	return s.VaultSecretService.Service.Read(ctx, in)
}
func (s *Server) UpdateVaultSecret(ctx context.Context, in *resources.UpdateVaultSecretRequest) (*resources.VaultSecretResource, error) {
	return s.VaultSecretService.Service.Update(ctx, in)
}
func (s *Server) DeleteVaultSecret(ctx context.Context, in *resources.DeleteVaultSecretRequest) (*common.Empty, error) {
	return s.VaultSecretService.Service.Delete(ctx, in)
}

func (s *Server) CreateVaultAccessPolicy(ctx context.Context, in *resources.CreateVaultAccessPolicyRequest) (*resources.VaultAccessPolicyResource, error) {
	return s.VaultAccessPolicyService.Service.Create(ctx, in)
}
func (s *Server) ReadVaultAccessPolicy(ctx context.Context, in *resources.ReadVaultAccessPolicyRequest) (*resources.VaultAccessPolicyResource, error) {
	return s.VaultAccessPolicyService.Service.Read(ctx, in)
}
func (s *Server) UpdateVaultAccessPolicy(ctx context.Context, in *resources.UpdateVaultAccessPolicyRequest) (*resources.VaultAccessPolicyResource, error) {
	return s.VaultAccessPolicyService.Service.Update(ctx, in)
}
func (s *Server) DeleteVaultAccessPolicy(ctx context.Context, in *resources.DeleteVaultAccessPolicyRequest) (*common.Empty, error) {
	return s.VaultAccessPolicyService.Service.Delete(ctx, in)
}

func (s *Server) CreateVirtualMachine(ctx context.Context, in *resources.CreateVirtualMachineRequest) (*resources.VirtualMachineResource, error) {
	return s.VirtualMachineService.Service.Create(ctx, in)
}
func (s *Server) ReadVirtualMachine(ctx context.Context, in *resources.ReadVirtualMachineRequest) (*resources.VirtualMachineResource, error) {
	return s.VirtualMachineService.Service.Read(ctx, in)
}
func (s *Server) UpdateVirtualMachine(ctx context.Context, in *resources.UpdateVirtualMachineRequest) (*resources.VirtualMachineResource, error) {
	return s.VirtualMachineService.Service.Update(ctx, in)
}
func (s *Server) DeleteVirtualMachine(ctx context.Context, in *resources.DeleteVirtualMachineRequest) (*common.Empty, error) {
	return s.VirtualMachineService.Service.Delete(ctx, in)
}
