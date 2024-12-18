package service

import (
	"context"

	v1 "gitlab.calendaria.team/services/rbac/api/rbac/v1"
	"gitlab.calendaria.team/services/rbac/internal/biz"
	"gitlab.calendaria.team/services/utils/v2/auth"
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

func (s *CheckPermissionsService) CheckPermissions(
	ctx context.Context, req *v1.CheckPermissionsRequest,
) (*v1.CheckPermissionsReply, error) {
	appID := auth.GetAppIdFromContext(ctx)
	if appID == "" {
		return nil, v1.ErrorEmptyAppId("empty app id")
	}

	identities := req.GetIdentities()
	tenantID := req.GetTenantId()
	// use context if request does not have tenantId and identities
	if tenantID == 0 || len(identities) == 0 {
		tenantID = auth.GetTenantIdFromContext(ctx)
		if tenantID == 0 {
			return nil, v1.ErrorEmptyActorId("empty tenant id")
		}

		identities = auth.GetIdentitiesFromContext(ctx)
		if len(identities) == 0 {
			return nil, v1.ErrorEmptyActorId("empty identities")
		}
	}

	permissionsMap, err := s.uc.CheckPermissions(ctx, tenantID, appID,
		identities, req.GetPermissions(), req.GetResources(),
		req.GetValue())
	if err != nil {
		return nil, err
	}

	return &v1.CheckPermissionsReply{
		Permissions: permissionsMap.GetPermissions(),
		Metadata:    permissionsMap.GetMetadata(),
	}, nil
}
