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

func TestPermissionsUsecase_GetPermissionById(t *testing.T) {
	logger := zap.NewZapLogger(true)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	permissionRepo := mock.NewMockPermissionRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	uc, err := biz.NewPermissionsUsecase(logger, permissionRepo, roleRepo, assignedRepo)
	require.NoError(t, err)

	ctx := context.Background()
	permissionId := "some.permission"

	permission := &ent.Permission{
		ID:          permissionId,
		AppID:       "app-id",
		Name:        "testName",
		Description: "testDesc",
	}
	permissionRepo.EXPECT().GetPermissionById(ctx, permissionId).Return(permission, nil)
	permissionRepo.EXPECT().GetPermissionById(ctx, gomock.Not(permissionId)).Return(nil, &ent.NotFoundError{})

	permission1, err := uc.GetPermissionById(ctx, permissionId)
	require.NoError(t, err)
	require.Equal(t, permission, permission1)

	permission2, err := uc.GetPermissionById(ctx, "other.permission")
	require.Error(t, err)
	require.Equal(t, v1.ErrorNotFound("permission not found"), err)
	require.Nil(t, permission2)
}

func TestPermissionsUsecase_CreatePermission(t *testing.T) {
	logger := zap.NewZapLogger(true)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	permissionRepo := mock.NewMockPermissionRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	uc, err := biz.NewPermissionsUsecase(logger, permissionRepo, roleRepo, assignedRepo)
	require.NoError(t, err)

	ctx := context.Background()
	dto := data.CreatePermissionDto{
		Id:          "some.permission",
		Name:        "testName",
		Description: "testDesc",
		AppId:       "app-id",
	}
	permission := &ent.Permission{
		ID:          "some.permission",
		AppID:       dto.AppId,
		Name:        dto.Name,
		Description: dto.Description,
	}
	permissionRepo.EXPECT().CreatePermission(ctx, dto).Return(permission, nil)

	permission1, err := uc.CreatePermission(ctx, dto)
	require.NoError(t, err)
	require.Equal(t, permission, permission1)
}

func TestPermissionsUsecase_UpdatePermission(t *testing.T) {
	logger := zap.NewZapLogger(true)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	permissionRepo := mock.NewMockPermissionRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	uc, err := biz.NewPermissionsUsecase(logger, permissionRepo, roleRepo, assignedRepo)
	require.NoError(t, err)

	ctx := context.Background()
	permissionId := "some.permission"
	dto := data.UpdatePermissionDto{
		Name:        "testName",
		Description: "testDesc",
	}
	permission := &ent.Permission{
		ID:          permissionId,
		AppID:       "app-id",
		Name:        dto.Name,
		Description: dto.Description,
	}
	permissionRepo.EXPECT().UpdatePermission(ctx, permissionId, dto).Return(permission, nil)

	permission1, err := uc.UpdatePermission(ctx, permissionId, dto)
	require.NoError(t, err)
	require.Equal(t, permission, permission1)
}

func TestPermissionsUsecase_DeletePermission(t *testing.T) {
	logger := zap.NewZapLogger(true)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	permissionRepo := mock.NewMockPermissionRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	uc, err := biz.NewPermissionsUsecase(logger, permissionRepo, roleRepo, assignedRepo)
	require.NoError(t, err)

	ctx := context.Background()
	permissionId := "some.permission"
	permissionRepo.EXPECT().DeletePermission(ctx, permissionId).Return(nil)

	err = uc.DeletePermission(ctx, permissionId)
	require.NoError(t, err)
}

func TestPermissionsUsecase_GetPermissions(t *testing.T) {
	logger := zap.NewZapLogger(true)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	permissionRepo := mock.NewMockPermissionRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	uc, err := biz.NewPermissionsUsecase(logger, permissionRepo, roleRepo, assignedRepo)
	require.NoError(t, err)

	ctx := context.Background()
	appId := "app-id"
	permissionIds := []string{"first.permission", "second.permission"}

	permissions := []*ent.Permission{
		{
			ID:          "first.permission",
			AppID:       appId,
			Name:        "firstName",
			Description: "firstDesc",
		},
		{
			ID:          "second.permission",
			AppID:       appId,
			Name:        "secondName",
			Description: "secondDesc",
		},
	}
	permissionRepo.EXPECT().GetPermissions(ctx, appId, permissionIds).Return(permissions, nil)

	permissions1, err := uc.GetPermissions(ctx, appId, permissionIds)
	require.NoError(t, err)
	require.Equal(t, permissions, permissions1)
}

func TestPermissionsUsecase_GetGroupedPermissions(t *testing.T) {
	logger := zap.NewZapLogger(true)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	permissionRepo := mock.NewMockPermissionRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	uc, err := biz.NewPermissionsUsecase(logger, permissionRepo, roleRepo, assignedRepo)
	require.NoError(t, err)

	ctx := context.Background()
	tenantId := int64(1)
	permission1 := &ent.Permission{
		ID:   "permission.one",
		Name: "Permission One",
	}
	permission2 := &ent.Permission{
		ID:   "permission.two",
		Name: "Permission Two",
	}
	permission3 := &ent.Permission{
		ID:   "permission.three",
		Name: "Permission Three",
	}
	permission4 := &ent.Permission{
		ID:   "permission.four",
		Name: "Permission Four",
	}
	permissionGroups := []*ent.PermissionGroup{
		{
			ID:   "permission.group.one",
			Name: "group1",
			Edges: ent.PermissionGroupEdges{
				Permissions: []*ent.Permission{permission1, permission2},
			},
		},
		{
			ID:   "permission.group.two",
			Name: "group2",
			Edges: ent.PermissionGroupEdges{
				Permissions: []*ent.Permission{permission3, permission4},
			},
		},
	}
	filter := data.FilterPermissions{
		WithDenied: true,
	}
	filter2 := data.FilterPermissions{}
	identities := []string{"identity1", "identity2"}

	assignedRoles := []*ent.ResourceAccess{
		{
			ID:         1,
			TenantID:   tenantId,
			IdentityID: "identity1",
			RoleID:     1,
		},
		{
			ID:         2,
			TenantID:   tenantId,
			IdentityID: "identity2",
			RoleID:     2,
		},
	}

	filterRolePermissions := data.FilterRolePermissions{
		TenantId:   tenantId,
		RolesIds:   []int64{1, 2},
		DeniedOnly: true,
	}
	permissionGroups2 := []*ent.PermissionGroup{
		{
			ID:   "permission.group.one",
			Name: "group1",
			Edges: ent.PermissionGroupEdges{
				Permissions: []*ent.Permission{permission2},
			},
		},
		{
			ID:   "permission.group.two",
			Name: "group2",
			Edges: ent.PermissionGroupEdges{
				Permissions: []*ent.Permission{permission3, permission4},
			},
		},
	}

	permissionRepo.EXPECT().GetGroupedPermissions(ctx, filter).Return(permissionGroups, nil)
	permissionRepo.EXPECT().GetGroupedPermissions(ctx, filter2).Return(permissionGroups, nil)
	assignedRepo.EXPECT().ListAssignedRoles(ctx, tenantId, identities, nil).Return(assignedRoles, nil)
	roleRepo.EXPECT().ListRolesPermissions(ctx, filterRolePermissions).Return([]*ent.RolePermission{
		{
			ID:           1,
			TenantID:     tenantId,
			RoleID:       1,
			PermissionID: permission1.ID,
		},
	}, nil)

	pg, err := uc.GetGroupedPermissions(ctx, tenantId, identities, filter)
	require.NoError(t, err)
	require.Equal(t, permissionGroups, pg)

	pg, err = uc.GetGroupedPermissions(ctx, tenantId, identities, filter2)
	require.NoError(t, err)
	require.Equal(t, permissionGroups2, pg)
}
