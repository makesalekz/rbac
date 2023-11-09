package service

import (
	"context"
	pb "rbac/api/rbac/v1"
	"rbac/internal/biz"
	"rbac/internal/data"
)

type PermissionsService struct {
	pb.UnimplementedPermissionsServer

	uc *biz.PermissionsUsecase
}

func NewPermissionsService(uc *biz.PermissionsUsecase) *PermissionsService {
	return &PermissionsService{
		uc: uc,
	}
}

func (s *PermissionsService) CreatePermission(ctx context.Context, req *pb.CreatePermissionRequest) (*pb.PermissionReply, error) {
	permission, err := s.uc.CreatePermission(ctx, data.CreatePermissionDto{
		Id:          req.Id,
		AppId:       req.AppId,
		Name:        req.Name,
		Description: req.Description,
		Fields:      req.Fields,
	})
	if err != nil {
		return nil, pb.ErrorDatabaseQuery(err.Error())
	}
	return &pb.PermissionReply{
		Id:          permission.ID,
		AppId:       permission.AppID,
		Name:        permission.Name,
		Description: permission.Description,
		Fields:      permission.Fields,
	}, nil
}
func (s *PermissionsService) UpdatePermission(ctx context.Context, req *pb.UpdatePermissionRequest) (*pb.PermissionReply, error) {
	permission, err := s.uc.UpdatePermission(ctx, req.PermissionId, data.UpdatePermissionDto{
		Name:        req.Name,
		Description: req.Description,
		Fields:      req.Fields,
	})
	if err != nil {
		return nil, err
	}
	return &pb.PermissionReply{
		Id:          permission.ID,
		AppId:       permission.AppID,
		Name:        permission.Name,
		Description: permission.Description,
		Fields:      permission.Fields,
	}, nil
}
func (s *PermissionsService) DeletePermission(ctx context.Context, req *pb.DeletePermissionRequest) (*pb.EmptyReply, error) {
	err := s.uc.DeletePermission(ctx, req.PermissionId)
	if err != nil {
		return nil, err
	}
	return &pb.EmptyReply{}, nil
}
func (s *PermissionsService) GetPermission(ctx context.Context, req *pb.GetPermissionRequest) (*pb.PermissionReply, error) {
	permission, err := s.uc.GetPermissionById(ctx, req.PermissionId)
	if err != nil {
		return nil, err
	}
	return &pb.PermissionReply{
		Id:          permission.ID,
		AppId:       permission.AppID,
		Name:        permission.Name,
		Description: permission.Description,
		Fields:      permission.Fields,
	}, nil
}
func (s *PermissionsService) ListPermissions(ctx context.Context, req *pb.ListPermissionsRequest) (*pb.ListPermissionsReply, error) {
	permissions, err := s.uc.GetPermissions(ctx, req.AppId, req.Ids)
	if err != nil {
		return nil, err
	}
	permissionsReply := make([]*pb.PermissionReply, len(permissions))
	for i, permission := range permissions {
		permissionsReply[i] = &pb.PermissionReply{
			Id:          permission.ID,
			AppId:       permission.AppID,
			Name:        permission.Name,
			Description: permission.Description,
			Fields:      permission.Fields,
		}
	}
	return &pb.ListPermissionsReply{
		Permissions: permissionsReply,
	}, nil
}
