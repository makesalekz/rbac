package service

import (
	"context"

	v1 "gitlab.calendaria.team/services/rbac/api/rbac/v1"
	"gitlab.calendaria.team/services/rbac/internal/biz"
)

type CheckPermissionsService struct {
	v1.UnimplementedCheckPermissionsServer

	uc *biz.CheckPermissionsUsecase
}

func NewCheckPermissionsService(
	uc *biz.CheckPermissionsUsecase,
) *CheckPermissionsService {
	return &CheckPermissionsService{
		uc: uc,
	}
}

func (s *CheckPermissionsService) CheckPermissions(ctx context.Context, req *v1.CheckPermissionsRequest) (*v1.CheckPermissionsReply, error) {
	permissionsMap, err := s.uc.CheckPermissions(ctx, req.TenantId, req.Identities, req.TeamId, req.Permissions)
	if err != nil {
		return nil, err
	}

	return &v1.CheckPermissionsReply{
		Permissions: permissionsMap,
	}, nil
}
