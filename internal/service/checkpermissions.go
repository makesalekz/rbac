package service

import (
	"context"

	v1 "gitlab.calendaria.team/services/rbac/api/rbac/v1"
	"gitlab.calendaria.team/services/rbac/internal/biz"
)

type CheckPermissionsService struct {
	v1.UnimplementedCheckPermissionsServer

	sh *ServiceHelper
	uc *biz.TeamIdentityUsecase
}

func NewCheckPermissionsService(
	sh *ServiceHelper,
	uc *biz.TeamIdentityUsecase,
) *CheckPermissionsService {
	return &CheckPermissionsService{
		sh: sh,
		uc: uc,
	}
}

func (s *CheckPermissionsService) CheckPermissions(ctx context.Context, req *v1.CheckPermissionsRequest) (*v1.CheckPermissionsReply, error) {
	tenantId, err := s.sh.GetTenantId(ctx, req.TenantId)
	if err != nil {
		return nil, v1.ErrorUnauthorized("invalid token")
	}

	identities, err := s.sh.GetIdentities(ctx, req.Identities)
	if err != nil {
		return nil, v1.ErrorUnauthorized("invalid token")
	}

	permissionsMap, err := s.uc.CheckPermissions(ctx, tenantId, identities, req.TeamId, req.Permissions)
	if err != nil {
		return nil, err
	}
	return &v1.CheckPermissionsReply{
		Permissions: permissionsMap,
	}, nil
}
