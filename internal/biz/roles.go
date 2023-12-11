package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	v1 "gitlab.calendaria.team/services/rbac/api/rbac/v1"
	"gitlab.calendaria.team/services/rbac/ent"
	"gitlab.calendaria.team/services/rbac/internal/data"
	"gitlab.calendaria.team/services/utils/v1/jwt"
)

// RolesUsecase .
type RolesUsecase struct {
	log      *log.Helper
	jwt      *jwt.JwtProcessor
	roleRepo data.RoleRepo
}

// NewRolesUsecase .
func NewRolesUsecase(logger log.Logger, jwt *jwt.JwtProcessor, usersRepo data.RoleRepo) (*RolesUsecase, error) {
	return &RolesUsecase{
		log:      log.NewHelper(logger),
		jwt:      jwt,
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

func (uc *RolesUsecase) UpdateRole(ctx context.Context, role *ent.Role, dto data.UpdateRoleDto) (*ent.Role, error) {
	if role.IsSystem {
		return nil, v1.ErrorForbidden("unable to edit system role")
	}

	return uc.roleRepo.UpdateRole(ctx, role, dto)
}

func (uc *RolesUsecase) DeleteRole(ctx context.Context, role *ent.Role) error {
	if role.IsSystem {
		return v1.ErrorForbidden("unable to edit system role")
	}

	return uc.roleRepo.DeleteRole(ctx, role)
}

func (uc *RolesUsecase) GetRoles(ctx context.Context, tenantId int64, search string) ([]*ent.Role, error) {
	return uc.roleRepo.GetRolesList(ctx, tenantId, search)
}

func (uc *RolesUsecase) CreateRole(ctx context.Context, createRoleDto data.CreateRoleDto) (*ent.Role, error) {
	return uc.roleRepo.CreateRole(ctx, createRoleDto)
}

func (uc *RolesUsecase) SetRolePermission(ctx context.Context, role *ent.Role, permission *ent.Permission, dto data.CreateRolePermissionDto) error {
	if role.IsSystem {
		return v1.ErrorForbidden("unable to edit system role")
	}

	if !validateFields(permission.Fields, dto.Fields) {
		return v1.ErrorInvalidRequest("invalid fields")
	}

	return uc.roleRepo.SetRolePermission(ctx, role, permission, dto)
}

func (uc *RolesUsecase) RemovePermissionFromRole(ctx context.Context, role *ent.Role, permission *ent.Permission) error {
	if role.IsSystem {
		return v1.ErrorForbidden("unable to edit system role")
	}

	return uc.roleRepo.RemovePermissionFromRole(ctx, role, permission)
}

func (uc *RolesUsecase) ListRolePermissions(ctx context.Context, role *ent.Role) ([]*ent.RolePermission, error) {
	return uc.roleRepo.ListRolePermissions(ctx, role)
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
