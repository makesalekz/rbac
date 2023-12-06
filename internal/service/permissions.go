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
}

func NewPermissionsService(uc *biz.PermissionsUsecase) *PermissionsService {
	return &PermissionsService{
		uc: uc,
	}
}

func (s *PermissionsService) CreatePermission(ctx context.Context, req *v1.CreatePermissionRequest) (*v1.PermissionReply, error) {
	// TODO check admin permissions

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
	return s.permissionReply(permission), nil
}

func (s *PermissionsService) UpdatePermission(ctx context.Context, req *v1.UpdatePermissionRequest) (*v1.PermissionReply, error) {
	// TODO check admin permissions

	permission, err := s.uc.UpdatePermission(ctx, req.PermissionId, data.UpdatePermissionDto{
		Name:        req.Name,
		Description: req.Description,
		Fields:      req.Fields,
	})
	if err != nil {
		return nil, err
	}
	return s.permissionReply(permission), nil
}

func (s *PermissionsService) DeletePermission(ctx context.Context, req *v1.DeletePermissionRequest) (*utils_v1.EmptyReply, error) {
	// TODO check admin permissions

	err := s.uc.DeletePermission(ctx, req.PermissionId)
	if err != nil {
		return nil, err
	}
	return &utils_v1.EmptyReply{}, nil
}

func (s *PermissionsService) GetPermission(ctx context.Context, req *v1.GetPermissionRequest) (*v1.PermissionReply, error) {
	// TODO check admin permissions

	permission, err := s.uc.GetPermissionById(ctx, req.PermissionId)
	if err != nil {
		return nil, err
	}
	return s.permissionReply(permission), nil
}

func (s *PermissionsService) ListPermissions(ctx context.Context, req *v1.ListPermissionsRequest) (*v1.ListPermissionsReply, error) {
	// TODO check admin permissions

	permissions, err := s.uc.GetPermissions(ctx, req.AppId, req.Ids)
	if err != nil {
		return nil, err
	}

	permissionsReply := make([]*v1.PermissionReply, len(permissions))
	for i, permission := range permissions {
		permissionsReply[i] = s.permissionReply(permission)
	}
	return &v1.ListPermissionsReply{
		Permissions: permissionsReply,
	}, nil
}

func (s *PermissionsService) permissionReply(permission *ent.Permission) *v1.PermissionReply {
	return &v1.PermissionReply{
		Id:          permission.ID,
		AppId:       permission.AppID,
		Name:        permission.Name,
		Description: permission.Description,
		Fields:      permission.Fields,
	}
}
