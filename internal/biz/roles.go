package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	v1 "gitlab.calendaria.team/services/rbac/api/rbac/v1"
	"gitlab.calendaria.team/services/rbac/ent"
	"gitlab.calendaria.team/services/rbac/internal/data"
)

// RolesUsecase .
type RolesUsecase struct {
	roleRepo data.RoleRepo
}

// NewRolesUsecase .
func NewRolesUsecase(
	logger log.Logger,
	usersRepo data.RoleRepo,
) (*RolesUsecase, error) {
	return &RolesUsecase{
		roleRepo: usersRepo,
	}, nil
}

func (uc *RolesUsecase) GetRoleById(ctx context.Context, tenantId, roleId int64) (*ent.Role, error) {
	role, err := uc.roleRepo.GetRoleById(ctx, tenantId, roleId)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, v1.ErrorNotFound("role not found")
		}
		return nil, err
	}

	return role, nil
}

func (uc *RolesUsecase) UpdateRole(ctx context.Context, tenantId, roleId int64, dto data.UpdateRoleDto) (*ent.Role, error) {
	return uc.roleRepo.UpdateRole(ctx, tenantId, roleId, dto)
}

func (uc *RolesUsecase) DeleteRole(ctx context.Context, tenantId, roleId int64) error {
	return uc.roleRepo.DeleteRole(ctx, tenantId, roleId)
}

func (uc *RolesUsecase) GetRoles(ctx context.Context, tenantId int64, search string) ([]*ent.Role, error) {
	return uc.roleRepo.GetRolesList(ctx, tenantId, search)
}

func (uc *RolesUsecase) CreateRole(ctx context.Context, dto data.CreateRoleDto) (*ent.Role, error) {
	return uc.roleRepo.CreateRole(ctx, dto)
}

func (uc *RolesUsecase) SetRolePermission(ctx context.Context, tenantId, roleId int64, permission *ent.Permission, dto data.CreateRolePermissionDto) error {
	if !validateFields(permission.Fields, dto.Fields) {
		return v1.ErrorBadRequest("fields not valid")
	}

	return uc.roleRepo.SetRolePermission(ctx, tenantId, roleId, permission, dto)
}

func (uc *RolesUsecase) RemovePermissionFromRole(ctx context.Context, tenantId, roleId int64, permission *ent.Permission) error {
	return uc.roleRepo.RemovePermissionFromRole(ctx, tenantId, roleId, permission)
}

func (uc *RolesUsecase) ListRolePermissions(ctx context.Context, tenantId, roleId int64) ([]*ent.RolePermission, error) {
	return uc.roleRepo.ListRolePermissions(ctx, tenantId, roleId)
}

func validateFields(sourceFields []string, targetFields []string) bool {
	if sourceFields == nil {
		return true
	}
	if targetFields == nil {
		return true
	}
	elementMap := make(map[string]bool)

	// Заполняем карту элементами исходного массива
	for _, element := range sourceFields {
		elementMap[element] = true
	}

	// Проверяем каждый элемент второго массива
	for _, element := range targetFields {
		if _, exists := elementMap[element]; !exists {
			return false
		}
	}

	return true
}
