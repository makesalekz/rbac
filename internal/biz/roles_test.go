package biz_test

import (
	"context"
	"testing"

	v1 "gitlab.calendaria.team/services/rbac/api/rbac/v1"
	"gitlab.calendaria.team/services/rbac/ent"
	"gitlab.calendaria.team/services/rbac/internal/biz"
	"gitlab.calendaria.team/services/rbac/internal/data"
	"gitlab.calendaria.team/services/rbac/internal/data/mock"
	"gitlab.calendaria.team/services/utils/v2/zap"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestRolesUsecase_GetRoleById(t *testing.T) {
	logger := zap.NewZapLogger(true)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rolesRepo := mock.NewMockRoleRepo(ctrl)
	uc, err := biz.NewRolesUsecase(logger, rolesRepo)
	require.NoError(t, err)

	ctx := context.Background()
	tenantId := int64(1)
	roleId := int64(1)

	role := &ent.Role{
		ID:          roleId,
		TenantID:    tenantId,
		Name:        "testName",
		Description: "testDesc",
	}
	rolesRepo.EXPECT().GetRoleById(ctx, tenantId, roleId).Return(role, nil)
	rolesRepo.EXPECT().GetRoleById(ctx, gomock.Any(), gomock.Not(roleId)).Return(nil, &ent.NotFoundError{})
	rolesRepo.EXPECT().GetRoleById(ctx, gomock.Not(tenantId), gomock.Any()).Return(nil, &ent.NotFoundError{})

	role1, err := uc.GetRoleById(ctx, tenantId, roleId)
	require.NoError(t, err)
	require.Equal(t, role, role1)

	role2, err := uc.GetRoleById(ctx, 2, roleId)
	require.Error(t, err)
	require.Equal(t, v1.ErrorNotFound("role not found"), err)
	require.Nil(t, role2)

	role3, err := uc.GetRoleById(ctx, tenantId, 3)
	require.Error(t, err)
	require.Equal(t, v1.ErrorNotFound("role not found"), err)
	require.Nil(t, role3)
}

func TestRolesUsecase_UpdateRole(t *testing.T) {
	logger := zap.NewZapLogger(true)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rolesRepo := mock.NewMockRoleRepo(ctrl)
	uc, err := biz.NewRolesUsecase(logger, rolesRepo)
	require.NoError(t, err)

	ctx := context.Background()
	tenantId := int64(1)
	roleId := int64(1)

	role := &ent.Role{
		ID:          roleId,
		TenantID:    tenantId,
		Name:        "testName",
		Description: "testDesc",
	}
	systemRole := &ent.Role{
		ID:          roleId,
		TenantID:    tenantId,
		Name:        "testName",
		Description: "testDesc",
		IsSystem:    true,
	}
	dto := data.UpdateRoleDto{
		Name:        "updName",
		Description: "updDesc",
	}
	roleUpdated := &ent.Role{
		ID:          roleId,
		TenantID:    tenantId,
		Name:        "updName",
		Description: "updDesc",
	}
	rolesRepo.EXPECT().UpdateRole(ctx, role, dto).Return(roleUpdated, nil)

	role1, err := uc.UpdateRole(ctx, role, dto)
	require.NoError(t, err)
	require.Equal(t, roleUpdated, role1)

	role2, err := uc.UpdateRole(ctx, systemRole, dto)
	require.Error(t, err)
	require.Equal(t, v1.ErrorForbidden("unable to edit system role"), err)
	require.Nil(t, role2)
}

func TestRolesUsecase_DeleteRole(t *testing.T) {
	logger := zap.NewZapLogger(true)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rolesRepo := mock.NewMockRoleRepo(ctrl)
	uc, err := biz.NewRolesUsecase(logger, rolesRepo)
	require.NoError(t, err)

	ctx := context.Background()
	tenantId := int64(1)
	roleId := int64(1)

	role := &ent.Role{
		ID:          roleId,
		TenantID:    tenantId,
		Name:        "testName",
		Description: "testDesc",
	}
	systemRole := &ent.Role{
		ID:          roleId,
		TenantID:    tenantId,
		Name:        "testName",
		Description: "testDesc",
		IsSystem:    true,
	}
	rolesRepo.EXPECT().DeleteRole(ctx, role).Return(nil)

	err = uc.DeleteRole(ctx, role)
	require.NoError(t, err)

	err = uc.DeleteRole(ctx, systemRole)
	require.Error(t, err)
	require.Equal(t, v1.ErrorForbidden("unable to delete system role"), err)
}

func TestRolesUsecase_GetRoles(t *testing.T) {
	logger := zap.NewZapLogger(true)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rolesRepo := mock.NewMockRoleRepo(ctrl)
	uc, err := biz.NewRolesUsecase(logger, rolesRepo)
	require.NoError(t, err)

	ctx := context.Background()
	tenantId := int64(1)
	search := "test"

	roles := []*ent.Role{
		{
			ID:          1,
			TenantID:    tenantId,
			Name:        "testName",
			Description: "testDesc",
		},
		{
			ID:          2,
			TenantID:    tenantId,
			Name:        "testName2",
			Description: "testDesc2",
		},
	}
	rolesRepo.EXPECT().GetRolesList(ctx, tenantId, search).Return(roles, nil)

	roles1, err := uc.GetRoles(ctx, tenantId, search)
	require.NoError(t, err)
	require.Equal(t, roles, roles1)
}

func TestRolesUsecase_CreateRole(t *testing.T) {
	logger := zap.NewZapLogger(true)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rolesRepo := mock.NewMockRoleRepo(ctrl)
	uc, err := biz.NewRolesUsecase(logger, rolesRepo)
	require.NoError(t, err)

	ctx := context.Background()
	createRoleDto := data.CreateRoleDto{
		TenantId:    1,
		Name:        "testName",
		Description: "testDesc",
	}
	role := &ent.Role{
		ID:          1,
		TenantID:    createRoleDto.TenantId,
		Name:        createRoleDto.Name,
		Description: createRoleDto.Description,
	}
	rolesRepo.EXPECT().CreateRole(ctx, createRoleDto).Return(role, nil)

	role1, err := uc.CreateRole(ctx, createRoleDto)
	require.NoError(t, err)
	require.Equal(t, role, role1)

	createRoleDto.IsSystem = true
	role2, err := uc.CreateRole(ctx, createRoleDto)
	require.Error(t, err)
	require.Equal(t, v1.ErrorForbidden("unable to create system role"), err)
	require.Nil(t, role2)
}

func TestRolesUsecase_SetRolePermission(t *testing.T) {
	logger := zap.NewZapLogger(true)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rolesRepo := mock.NewMockRoleRepo(ctrl)
	uc, err := biz.NewRolesUsecase(logger, rolesRepo)
	require.NoError(t, err)

	ctx := context.Background()
	tenantId := int64(1)
	roleId := int64(1)
	permissionId := "some.permission"

	role := &ent.Role{
		ID:          roleId,
		TenantID:    tenantId,
		Name:        "testName",
		Description: "testDesc",
	}
	systemRole := &ent.Role{
		ID:          roleId,
		TenantID:    tenantId,
		Name:        "testName",
		Description: "testDesc",
		IsSystem:    true,
	}
	permission := &ent.Permission{
		ID:          permissionId,
		Name:        "testName",
		Description: "testDesc",
		Fields:      []string{"field1", "field2"},
	}
	dto := data.CreateRolePermissionDto{
		Fields: []string{"field1"},
	}
	dto2 := data.CreateRolePermissionDto{
		Fields: []string{"field3"},
	}
	rolesRepo.EXPECT().SetRolePermission(ctx, role, permission, dto).Return(nil)

	err = uc.SetRolePermission(ctx, role, permission, dto)
	require.NoError(t, err)

	err = uc.SetRolePermission(ctx, systemRole, permission, dto)
	require.Error(t, err)
	require.Equal(t, v1.ErrorForbidden("unable to edit system role"), err)

	err = uc.SetRolePermission(ctx, role, permission, dto2)
	require.Error(t, err)
	require.Equal(t, v1.ErrorBadRequest("fields not valid"), err)
}

func TestRolesUsecase_RemovePermissionFromRole(t *testing.T) {
	logger := zap.NewZapLogger(true)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rolesRepo := mock.NewMockRoleRepo(ctrl)
	uc, err := biz.NewRolesUsecase(logger, rolesRepo)
	require.NoError(t, err)

	ctx := context.Background()
	tenantId := int64(1)
	roleId := int64(1)
	permissionId := "some.permission"

	role := &ent.Role{
		ID:          roleId,
		TenantID:    tenantId,
		Name:        "testName",
		Description: "testDesc",
	}
	systemRole := &ent.Role{
		ID:          roleId,
		TenantID:    tenantId,
		Name:        "testName",
		Description: "testDesc",
		IsSystem:    true,
	}
	permission := &ent.Permission{
		ID:          permissionId,
		Name:        "testName",
		Description: "testDesc",
		Fields:      []string{"field1", "field2"},
	}
	rolesRepo.EXPECT().RemovePermissionFromRole(ctx, role, permission).Return(nil)

	err = uc.RemovePermissionFromRole(ctx, role, permission)
	require.NoError(t, err)

	err = uc.RemovePermissionFromRole(ctx, systemRole, permission)
	require.Error(t, err)
	require.Equal(t, v1.ErrorForbidden("unable to edit system role"), err)
}

func TestRolesUsecase_ListRolePermissions(t *testing.T) {
	logger := zap.NewZapLogger(true)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rolesRepo := mock.NewMockRoleRepo(ctrl)
	uc, err := biz.NewRolesUsecase(logger, rolesRepo)
	require.NoError(t, err)

	ctx := context.Background()
	tenantId := int64(1)
	roleId := int64(1)

	role := &ent.Role{
		ID:          roleId,
		TenantID:    tenantId,
		Name:        "testName",
		Description: "testDesc",
	}
	permissions := []*ent.RolePermission{
		{
			ID:           1,
			TenantID:     tenantId,
			RoleID:       roleId,
			PermissionID: "some.permission",
			Fields:       []string{"field1", "field2"},
		},
		{
			ID:           2,
			TenantID:     tenantId,
			RoleID:       roleId,
			PermissionID: "some.permission2",
			Fields:       []string{"field3", "field4"},
		},
	}
	rolesRepo.EXPECT().ListRolePermissions(ctx, role).Return(permissions, nil)

	permissions1, err := uc.ListRolePermissions(ctx, role)
	require.NoError(t, err)
	require.Equal(t, permissions, permissions1)
}
