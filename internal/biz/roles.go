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
	repo data.RoleRepo
}

// NewRolesUsecase .
func NewRolesUsecase(
	logger log.Logger,
	repo data.RoleRepo,
) (*RolesUsecase, error) {
	return &RolesUsecase{
		repo: repo,
	}, nil
}

func (uc *RolesUsecase) GetRoleByID(ctx context.Context, tenantID, roleID int64) (*ent.Role, error) {
	role, err := uc.repo.GetRoleByID(ctx, tenantID, roleID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, v1.ErrorNotFound("role not found")
		}
		return nil, err
	}

	return role, nil
}

func (uc *RolesUsecase) UpdateRole(
	ctx context.Context,
	tenantID, roleID int64,
	dto data.UpdateRoleDto,
) (*ent.Role, error) {
	return uc.repo.UpdateRole(ctx, tenantID, roleID, dto)
}

func (uc *RolesUsecase) DeleteRole(ctx context.Context, tenantID, roleID int64) error {
	return uc.repo.DeleteRole(ctx, tenantID, roleID)
}

func (uc *RolesUsecase) GetRoles(ctx context.Context, tenantID int64, search, appID string) ([]*ent.Role, error) {
	var isSystem = false
	if appID == PmsAppID {
		isSystem = true
	}

	return uc.repo.GetRolesList(ctx, tenantID, search, isSystem)
}

func (uc *RolesUsecase) CreateRole(ctx context.Context, dto data.CreateRoleDto) (*ent.Role, error) {
	return uc.repo.CreateRole(ctx, dto)
}

func (uc *RolesUsecase) SetRolePermission(
	ctx context.Context,
	tenantID, roleID int64,
	permission *ent.Permission,
	dto data.CreateRolePermissionDto,
) error {
	if !validateFields(permission.Fields, dto.Fields) {
		return v1.ErrorBadRequest("fields not valid")
	}

	return uc.repo.SetRolePermission(ctx, tenantID, roleID, permission, dto)
}

func (uc *RolesUsecase) RemovePermissionFromRole(
	ctx context.Context,
	tenantID, roleID int64,
	permission *ent.Permission,
) error {
	return uc.repo.RemovePermissionFromRole(ctx, tenantID, roleID, permission)
}

func (uc *RolesUsecase) ListRolePermissions(
	ctx context.Context,
	tenantID, roleID int64,
) ([]*ent.RolePermission, error) {
	return uc.repo.ListRolePermissions(ctx, tenantID, roleID)
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
