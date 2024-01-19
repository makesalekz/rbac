package service

import (
	"context"

	v1 "gitlab.calendaria.team/services/rbac/api/rbac/v1"
	"gitlab.calendaria.team/services/utils/v1/jwt"
)

type ServiceHelper struct {
	jwt *jwt.JwtProcessor
}

func NewServiceHelper(
	jwt *jwt.JwtProcessor,
) *ServiceHelper {
	return &ServiceHelper{
		jwt: jwt,
	}
}

func (s *ServiceHelper) GetTenantId(ctx context.Context, reqTenantId int64) (int64, error) {
	// TODO: remove getting from context
	claims, ok := s.jwt.GetClaimsFromContext(ctx)
	if ok && claims.IsUserTenantRequest() {
		return claims.GetTenantId(), nil
	}

	if reqTenantId != 0 {
		return reqTenantId, nil
	}
	return 0, v1.ErrorUnauthorized("invalid token")
}

func (s *ServiceHelper) GetIdentities(ctx context.Context, reqIdentities []string) ([]string, error) {
	// TODO: remove getting from context
	claims, ok := s.jwt.GetClaimsFromContext(ctx)
	if ok && claims.IsUserTenantRequest() {
		return claims.GetIdentities(), nil
	}

	if reqIdentities != nil {
		return reqIdentities, nil
	}
	return nil, v1.ErrorUnauthorized("invalid token")
}
