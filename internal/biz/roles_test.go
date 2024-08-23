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
	tenantID := int64(1)
	roleID := int64(1)

	role := &ent.Role{
		ID:          roleID,
		TenantID:    tenantID,
		Name:        "testName",
		Description: "testDesc",
	}
	rolesRepo.EXPECT().GetRoleByID(ctx, tenantID, roleID).Return(role, nil)
	rolesRepo.EXPECT().GetRoleByID(ctx, gomock.Any(), gomock.Not(roleID)).Return(nil, &ent.NotFoundError{})
	rolesRepo.EXPECT().GetRoleByID(ctx, gomock.Not(tenantID), gomock.Any()).Return(nil, &ent.NotFoundError{})

	role1, err := uc.GetRoleByID(ctx, tenantID, roleID)
	require.NoError(t, err)
	require.Equal(t, role, role1)

	role2, err := uc.GetRoleByID(ctx, 2, roleID)
	require.Error(t, err)
	require.Equal(t, v1.ErrorNotFound("role not found"), err)
	require.Nil(t, role2)

	role3, err := uc.GetRoleByID(ctx, tenantID, 3)
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
	tenantID := int64(1)
	roleID := int64(1)

	dto := data.UpdateRoleDto{
		Name:        "updName",
		Description: "updDesc",
	}
	roleUpdated := &ent.Role{
		ID:          roleID,
		TenantID:    tenantID,
		Name:        "updName",
		Description: "updDesc",
	}
	rolesRepo.EXPECT().UpdateRole(ctx, tenantID, roleID, dto).Return(roleUpdated, nil)

	role1, err := uc.UpdateRole(ctx, tenantID, roleID, dto)
	require.NoError(t, err)
	require.Equal(t, roleUpdated, role1)
}

func TestRolesUsecase_DeleteRole(t *testing.T) {
	logger := zap.NewZapLogger(true)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rolesRepo := mock.NewMockRoleRepo(ctrl)
	uc, err := biz.NewRolesUsecase(logger, rolesRepo)
	require.NoError(t, err)

	ctx := context.Background()
	tenantID := int64(1)
	roleID := int64(1)

	rolesRepo.EXPECT().DeleteRole(ctx, tenantID, roleID).Return(nil)

	err = uc.DeleteRole(ctx, tenantID, roleID)
	require.NoError(t, err)
}

func TestRolesUsecase_GetRoles(t *testing.T) {
	logger := zap.NewZapLogger(true)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rolesRepo := mock.NewMockRoleRepo(ctrl)
	uc, err := biz.NewRolesUsecase(logger, rolesRepo)
	require.NoError(t, err)

	ctx := context.Background()
	tenantID := int64(1)
	search := "test"

	roles := []*ent.Role{
		{
			ID:          1,
			TenantID:    tenantID,
			Name:        "testName",
			Description: "testDesc",
		},
		{
			ID:          2,
			TenantID:    tenantID,
			Name:        "testName2",
			Description: "testDesc2",
		},
	}
	rolesRepo.EXPECT().GetRolesList(ctx, tenantID, search, false).Return(roles, nil)

	roles1, err := uc.GetRoles(ctx, tenantID, search, "")
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
		TenantID:    1,
		Name:        "testName",
		Description: "testDesc",
	}
	role := &ent.Role{
		ID:          1,
		TenantID:    createRoleDto.TenantID,
		Name:        createRoleDto.Name,
		Description: createRoleDto.Description,
	}
	rolesRepo.EXPECT().CreateRole(ctx, createRoleDto).Return(role, nil)

	role1, err := uc.CreateRole(ctx, createRoleDto)
	require.NoError(t, err)
	require.Equal(t, role, role1)
}

func TestRolesUsecase_SetRolePermission(t *testing.T) {
	logger := zap.NewZapLogger(true)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rolesRepo := mock.NewMockRoleRepo(ctrl)
	uc, err := biz.NewRolesUsecase(logger, rolesRepo)
	require.NoError(t, err)

	ctx := context.Background()
	tenantID := int64(1)
	roleID := int64(1)
	permissionID := "some.permission"

	permission := &ent.Permission{
		ID:          permissionID,
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
	rolesRepo.EXPECT().SetRolePermission(ctx, tenantID, roleID, permission, dto).Return(nil)

	err = uc.SetRolePermission(ctx, tenantID, roleID, permission, dto)
	require.NoError(t, err)

	err = uc.SetRolePermission(ctx, tenantID, roleID, permission, dto2)
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
	tenantID := int64(1)
	roleID := int64(1)
	permissionID := "some.permission"

	permission := &ent.Permission{
		ID:          permissionID,
		Name:        "testName",
		Description: "testDesc",
		Fields:      []string{"field1", "field2"},
	}
	rolesRepo.EXPECT().RemovePermissionFromRole(ctx, tenantID, roleID, permission).Return(nil)

	err = uc.RemovePermissionFromRole(ctx, tenantID, roleID, permission)
	require.NoError(t, err)
}

func TestRolesUsecase_ListRolePermissions(t *testing.T) {
	logger := zap.NewZapLogger(true)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	rolesRepo := mock.NewMockRoleRepo(ctrl)
	uc, err := biz.NewRolesUsecase(logger, rolesRepo)
	require.NoError(t, err)

	ctx := context.Background()
	tenantID := int64(1)
	roleID := int64(1)

	permissions := []*ent.RolePermission{
		{
			ID:           1,
			TenantID:     tenantID,
			RoleID:       roleID,
			PermissionID: "some.permission",
			Fields:       []string{"field1", "field2"},
		},
		{
			ID:           2,
			TenantID:     tenantID,
			RoleID:       roleID,
			PermissionID: "some.permission2",
			Fields:       []string{"field3", "field4"},
		},
	}
	rolesRepo.EXPECT().ListRolePermissions(ctx, tenantID, roleID).Return(permissions, nil)

	permissions1, err := uc.ListRolePermissions(ctx, tenantID, roleID)
	require.NoError(t, err)
	require.Equal(t, permissions, permissions1)
}
