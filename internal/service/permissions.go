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

	uc *biz.PermissionsUsecase
	au *biz.TeamIdentityUsecase
}

func NewPermissionsService(uc *biz.PermissionsUsecase, au *biz.TeamIdentityUsecase) *PermissionsService {
	return &PermissionsService{
		uc: uc,
		au: au,
	}
}

func (s *PermissionsService) CreatePermission(ctx context.Context, req *v1.CreatePermissionRequest) (*v1.PermissionReply, error) {
	fields, err := s.au.HasPermission(ctx, "admin.permission.create")
	if err != nil {
		return nil, err
	}
	if fields == nil {
		return nil, v1.ErrorForbidden("has no permission")
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
	fields, err := s.au.HasPermission(ctx, "admin.permission.update")
	if err != nil {
		return nil, err
	}
	if fields == nil {
		return nil, v1.ErrorForbidden("has no permission")
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

func (s *PermissionsService) DeletePermission(ctx context.Context, req *v1.DeletePermissionRequest) (*utils_v1.EmptyReply, error) {
	fields, err := s.au.HasPermission(ctx, "admin.permission.delete")
	if err != nil {
		return nil, err
	}
	if fields == nil {
		return nil, v1.ErrorForbidden("has no permission")
	}

	err = s.uc.DeletePermission(ctx, req.PermissionId)
	if err != nil {
		return nil, err
	}
	return &utils_v1.EmptyReply{}, nil
}

func (s *PermissionsService) GetPermission(ctx context.Context, req *v1.GetPermissionRequest) (*v1.PermissionReply, error) {
	fields, err := s.au.HasPermission(ctx, "admin.permission.read")
	if err != nil {
		return nil, err
	}
	if fields == nil {
		return nil, v1.ErrorForbidden("has no permission")
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
	fields, err := s.au.HasPermission(ctx, "admin.permission.read")
	if err != nil {
		return nil, err
	}
	if fields == nil {
		return nil, v1.ErrorForbidden("has no permission")
	}

	groups, err := s.uc.GetGroupedPermissions(ctx, data.FilterPermissions{
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
