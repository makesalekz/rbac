package service

import (
	"context"

	pb "rbac/api/rbac/v1"
)

type CheckPermissionsService struct {
	pb.UnimplementedCheckPermissionsServer
}

func NewCheckPermissionsService() *CheckPermissionsService {
	return &CheckPermissionsService{}
}

func (s *CheckPermissionsService) CheckPermissions(ctx context.Context, req *pb.CheckPermissionsRequest) (*pb.CheckPermissionsReply, error) {
	return &pb.CheckPermissionsReply{}, nil
}
