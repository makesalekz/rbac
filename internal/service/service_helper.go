package service

import (
	"context"

	v1 "github.com/makesalekz/rbac/api/rbac/v1"
	"github.com/makesalekz/rbac/internal/biz"
	"github.com/makesalekz/utils/v2/auth"
)

type ServiceHelper struct {
	uc *biz.CheckPermissionsUsecase
}

func NewServiceHelper(
	uc *biz.CheckPermissionsUsecase,
) *ServiceHelper {
	return &ServiceHelper{
		uc: uc,
	}
}

func (s *ServiceHelper) HasPermission(ctx context.Context, permission string) (int64, *v1.ListOfFields, error) {
	tenantID := auth.GetTenantIdFromContext(ctx)
	if tenantID == 0 {
		return 0, nil, v1.ErrorEmptyActorId("empty tenant id")
	}

	appID := auth.GetAppIdFromContext(ctx)
	if appID == "" {
		return 0, nil, v1.ErrorEmptyAppId("empty app id")
	}

	identities := auth.GetIdentitiesFromContext(ctx)
	if len(identities) == 0 {
		return 0, nil, v1.ErrorEmptyActorId("empty identities")
	}

	fields, err := s.uc.HasPermission(ctx, tenantID, appID, identities, permission)
	if err != nil {
		return 0, nil, err
	}
	if fields == nil {
		return 0, nil, v1.ErrorForbidden("has no permission")
	}

	return tenantID, fields, nil
}
