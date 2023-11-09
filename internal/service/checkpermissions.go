package service

import (
	"context"
	"rbac/internal/biz"

	pb "rbac/api/rbac/v1"
)

type CheckPermissionsService struct {
	pb.UnimplementedCheckPermissionsServer

	uc *biz.TeamIdentityUsecase
}

func NewCheckPermissionsService(uc *biz.TeamIdentityUsecase) *CheckPermissionsService {
	return &CheckPermissionsService{
		uc: uc,
	}
}

func (s *CheckPermissionsService) CheckPermissions(ctx context.Context, req *pb.CheckPermissionsRequest) (*pb.CheckPermissionsReply, error) {
	permissionsMap, err := s.uc.CheckPermissions(ctx, req.TeamId, req.Permissions)
	if err != nil {
		return nil, err
	}
	return &pb.CheckPermissionsReply{
		Permissions: permissionsMap,
	}, nil
}
