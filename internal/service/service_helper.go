package service

import (
	"context"

	v1 "gitlab.calendaria.team/services/rbac/api/rbac/v1"
	"gitlab.calendaria.team/services/rbac/internal/biz"
	"gitlab.calendaria.team/services/utils/v2/auth"
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
	tenantId := auth.GetTenantIdFromContext(ctx)
	if tenantId == 0 {
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

	fields, err := s.uc.HasPermission(ctx, tenantId, appID, identities, permission)
	if err != nil {
		return 0, nil, err
	}
	if fields == nil {
		return 0, nil, v1.ErrorForbidden("has no permission")
	}

	return tenantId, fields, nil
}
