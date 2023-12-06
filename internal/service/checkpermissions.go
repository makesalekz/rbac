package service

import (
	"context"

	v1 "gitlab.calendaria.team/services/rbac/api/rbac/v1"
	"gitlab.calendaria.team/services/rbac/internal/biz"
	"gitlab.calendaria.team/services/utils/v1/jwt"
)

type CheckPermissionsService struct {
	v1.UnimplementedCheckPermissionsServer

	jwt *jwt.JwtProcessor
	uc  *biz.TeamIdentityUsecase
}

func NewCheckPermissionsService(jwt *jwt.JwtProcessor, uc *biz.TeamIdentityUsecase) *CheckPermissionsService {
	return &CheckPermissionsService{
		jwt: jwt,
		uc:  uc,
	}
}

func (s *CheckPermissionsService) CheckPermissions(ctx context.Context, req *v1.CheckPermissionsRequest) (*v1.CheckPermissionsReply, error) {
	permissionsMap, err := s.uc.CheckPermissions(ctx, req.TeamId, req.Permissions)
	if err != nil {
		return nil, err
	}
	return &v1.CheckPermissionsReply{
		Permissions: permissionsMap,
	}, nil
}
