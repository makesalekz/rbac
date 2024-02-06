package service

import (
	"context"

	v1 "gitlab.calendaria.team/services/rbac/api/rbac/v1"
	"gitlab.calendaria.team/services/rbac/internal/biz"
	"gitlab.calendaria.team/services/utils/v1/jwt"
)

type ServiceHelper struct {
	jwt *jwt.JwtProcessor
	uc  *biz.CheckPermissionsUsecase
}

func NewServiceHelper(
	jwt *jwt.JwtProcessor,
	uc *biz.CheckPermissionsUsecase,
) *ServiceHelper {
	return &ServiceHelper{
		jwt: jwt,
		uc:  uc,
	}
}

func (s *ServiceHelper) GetClaims(ctx context.Context) (*jwt.TenantClaims, error) {
	claims, ok := s.jwt.GetClaimsFromContext(ctx)
	if ok && claims.IsUserTenantRequest() {
		return claims, nil
	}

	return nil, v1.ErrorUnauthorized("invalid token")
}

func (s *ServiceHelper) HasPermission(ctx context.Context, permission string) (*jwt.TenantClaims, *v1.ListOfFields, error) {
	claims, err := s.GetClaims(ctx)
	if err != nil {
		return nil, nil, err
	}

	fields, err := s.uc.HasPermission(ctx, claims.GetTenantId(), claims.GetIdentities(), permission)
	if err != nil {
		return nil, nil, err
	}
	if fields == nil {
		return nil, nil, v1.ErrorForbidden("has no permission")
	}

	return claims, fields, nil
}
