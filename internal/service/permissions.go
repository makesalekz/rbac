package service

import (
	"context"

	v1 "gitlab.calendaria.team/services/rbac/api/rbac/v1"
	"gitlab.calendaria.team/services/rbac/ent"
	"gitlab.calendaria.team/services/rbac/internal/biz"
	"gitlab.calendaria.team/services/rbac/internal/data"
	utils_v1 "gitlab.calendaria.team/services/utils/api/utils/v1"
)

type PermissionsService struct {
	v1.UnimplementedPermissionsServer

	sh *ServiceHelper
	uc *biz.PermissionsUsecase
}

func NewPermissionsService(
	sh *ServiceHelper,
	uc *biz.PermissionsUsecase,
) *PermissionsService {
	return &PermissionsService{
		sh: sh,
		uc: uc,
	}
}

func (s *PermissionsService) CreatePermission(ctx context.Context, req *v1.CreatePermissionRequest) (*v1.PermissionReply, error) {
	_, _, err := s.sh.HasPermission(ctx, "admin.permission.create")
	if err != nil {
		return nil, err
	}

	permission, err := s.uc.CreatePermission(ctx, data.CreatePermissionDto{
		Id:          req.Id,
		AppId:       req.AppId,
		Name:        req.Name,
		Description: req.Description,
		Fields:      req.Fields,
	})
	if err != nil {
		return nil, v1.ErrorDatabaseQuery(err.Error())
	}
	return &v1.PermissionReply{
		Permission: s.permissionReply(permission),
	}, nil
}

func (s *PermissionsService) UpdatePermission(ctx context.Context, req *v1.UpdatePermissionRequest) (*v1.PermissionReply, error) {
	_, _, err := s.sh.HasPermission(ctx, "admin.permission.update")
	if err != nil {
		return nil, err
	}

	permission, err := s.uc.UpdatePermission(ctx, req.PermissionId, data.UpdatePermissionDto{
		Name:        req.Name,
		Description: req.Description,
		Fields:      req.Fields,
	})
	if err != nil {
		return nil, err
	}
	return &v1.PermissionReply{
		Permission: s.permissionReply(permission),
	}, nil
}

func (s *PermissionsService) DeletePermission(ctx context.Context, req *v1.PermissionRequest) (*utils_v1.EmptyReply, error) {
	_, _, err := s.sh.HasPermission(ctx, "admin.permission.delete")
	if err != nil {
		return nil, err
	}

	err = s.uc.DeletePermission(ctx, req.PermissionId)
	if err != nil {
		return nil, err
	}
	return &utils_v1.EmptyReply{}, nil
}

func (s *PermissionsService) GetPermission(ctx context.Context, req *v1.PermissionRequest) (*v1.PermissionReply, error) {
	_, _, err := s.sh.HasPermission(ctx, "admin.permission.read")
	if err != nil {
		return nil, err
	}

	permission, err := s.uc.GetPermissionById(ctx, req.PermissionId)
	if err != nil {
		return nil, err
	}
	return &v1.PermissionReply{
		Permission: s.permissionReply(permission),
	}, nil
}

func (s *PermissionsService) ListPermissions(ctx context.Context, req *v1.ListPermissionsRequest) (*v1.ListPermissionsReply, error) {
	claims, _, err := s.sh.HasPermission(ctx, "admin.permission.read")
	if err != nil {
		return nil, err
	}

	groups, err := s.uc.GetGroupedPermissions(
		ctx,
		claims.GetTenantId(),
		claims.GetIdentities(),
		data.FilterPermissions{
			AppsIds: req.AppsIds,
		})
	if err != nil {
		return nil, err
	}

	groupsReply := make([]*v1.Group, len(groups))
	for i, group := range groups {
		groupsReply[i] = s.groupReply(group)
	}

	return &v1.ListPermissionsReply{
		Groups: groupsReply,
	}, nil
}

func (s *PermissionsService) permissionReply(permission *ent.Permission) *v1.Permission {
	return &v1.Permission{
		Id:          permission.ID,
		AppId:       permission.AppID,
		Name:        permission.Name,
		Description: permission.Description,
		Fields:      permission.Fields,
	}
}

func (s *PermissionsService) groupReply(group *ent.PermissionGroup) *v1.Group {
	permissions := make([]*v1.Permission, len(group.Edges.Permissions))
	for i, permission := range group.Edges.Permissions {
		permissions[i] = s.permissionReply(permission)
	}

	return &v1.Group{
		Id:          group.ID,
		AppId:       group.AppID,
		Name:        group.Name,
		Permissions: permissions,
	}
}
