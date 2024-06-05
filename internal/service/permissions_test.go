package service_test

import (
	"context"
	"testing"

	v1 "gitlab.calendaria.team/services/rbac/api/rbac/v1"
	"gitlab.calendaria.team/services/rbac/ent"
	"gitlab.calendaria.team/services/rbac/internal/biz"
	"gitlab.calendaria.team/services/rbac/internal/data"
	"gitlab.calendaria.team/services/rbac/internal/data/mock"
	"gitlab.calendaria.team/services/rbac/internal/service"
	"gitlab.calendaria.team/services/utils/v2/zap"

	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func mockServerContext() context.Context {
	md := metadata.Metadata{
		"x-md-global-tenant-id":  []string{"12345"},
		"x-md-global-actor-id":   []string{"1234567"},
		"x-md-global-identities": []string{"identity1", "identity2"},
		"x-md-global-app-id":     []string{"app-id"},
	}
	return metadata.NewServerContext(context.Background(), md)
}

func mockRolePermissions(ids ...string) []*ent.RolePermission {
	var roles []*ent.RolePermission
	for _, id := range ids {
		roles = append(roles, &ent.RolePermission{
			PermissionID: id,
			Fields:       []string{},
		})
	}
	return roles
}

func mockAccess() []*ent.ResourceAccess {
	return []*ent.ResourceAccess{
		{
			RoleID: 1,
		},
		{
			RoleID: 2,
		},
	}
}

// ------------------------ CreatePermission ------------------------

func TestPermissionsService_CreatePermission(t *testing.T) {
	logger := zap.NewZapLogger(true)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	permissionRepo := mock.NewMockPermissionRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	uc, err := biz.NewPermissionsUsecase(logger, permissionRepo, roleRepo, assignedRepo)
	require.NoError(t, err)

	check, err := biz.NewCheckPermissionsUsecase(assignedRepo, roleRepo, teamsRepo)
	require.NoError(t, err)

	sh := service.NewServiceHelper(check)

	service := service.NewPermissionsService(sh, uc, check)

	ctx := mockServerContext()

	req := &v1.CreatePermissionRequest{
		Id:      "some.group.permission",
		GroupId: "some.group",
		AppId:   "app-id",
		Name:    "testName",
	}
	permission := &ent.Permission{
		ID:          req.GetId(),
		GroupID:     req.GetGroupId(),
		AppID:       req.GetAppId(),
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Fields:      req.GetFields(),
	}
	assignedRepo.EXPECT().CheckRoles(ctx, gomock.Any()).Return(mockAccess(), nil)
	roleRepo.EXPECT().ListRolesPermissions(ctx, gomock.Any()).Return(mockRolePermissions("admin.permission.create"), nil)
	permissionRepo.EXPECT().CreatePermission(ctx, gomock.Any()).Return(permission, nil)

	expect := &v1.Permission{
		Id:          permission.ID,
		GroupId:     permission.GroupID,
		AppId:       permission.AppID,
		Name:        permission.Name,
		Description: permission.Description,
		Fields:      permission.Fields,
	}

	reply, err := service.CreatePermission(ctx, req)
	require.NoError(t, err)
	require.Equal(t, expect, reply.GetPermission())
}

func TestPermissionsService_CreatePermissionAccessDenied(t *testing.T) {
	logger := zap.NewZapLogger(true)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	permissionRepo := mock.NewMockPermissionRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	uc, err := biz.NewPermissionsUsecase(logger, permissionRepo, roleRepo, assignedRepo)
	require.NoError(t, err)

	check, err := biz.NewCheckPermissionsUsecase(assignedRepo, roleRepo, teamsRepo)
	require.NoError(t, err)

	sh := service.NewServiceHelper(check)

	service := service.NewPermissionsService(sh, uc, check)

	ctx := mockServerContext()

	req := &v1.CreatePermissionRequest{
		Id:          "some.group.permission",
		GroupId:     "some.group",
		AppId:       "app-id",
		Name:        "testName",
		Description: "testDesc",
		Fields:      []string{"field1", "field2"},
	}
	assignedRepo.EXPECT().CheckRoles(ctx, gomock.Any()).Return(mockAccess(), nil)
	roleRepo.EXPECT().ListRolesPermissions(ctx, gomock.Any()).Return(mockRolePermissions("wrong.permission"), nil)

	reply, err := service.CreatePermission(ctx, req)
	require.Error(t, err)
	require.Equal(t, v1.ErrorForbidden("has no permission"), err)
	require.Nil(t, reply)
}

func TestPermissionsService_CreatePermissionEmptyID(t *testing.T) {
	logger := zap.NewZapLogger(true)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	permissionRepo := mock.NewMockPermissionRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	uc, err := biz.NewPermissionsUsecase(logger, permissionRepo, roleRepo, assignedRepo)
	require.NoError(t, err)

	check, err := biz.NewCheckPermissionsUsecase(assignedRepo, roleRepo, teamsRepo)
	require.NoError(t, err)

	sh := service.NewServiceHelper(check)

	service := service.NewPermissionsService(sh, uc, check)

	ctx := mockServerContext()

	req := &v1.CreatePermissionRequest{
		GroupId:     "some.group",
		AppId:       "app-id",
		Name:        "testName",
		Description: "testDesc",
		Fields:      []string{"field1", "field2"},
	}
	assignedRepo.EXPECT().CheckRoles(ctx, gomock.Any()).Return(mockAccess(), nil)
	roleRepo.EXPECT().ListRolesPermissions(ctx, gomock.Any()).Return(mockRolePermissions("admin.permission.create"), nil)

	reply, err := service.CreatePermission(ctx, req)
	require.Error(t, err)
	require.Equal(t, v1.ErrorBadRequest("empty id"), err)
	require.Nil(t, reply)
}

func TestPermissionsService_CreatePermissionEmptyName(t *testing.T) {
	logger := zap.NewZapLogger(true)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	permissionRepo := mock.NewMockPermissionRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	uc, err := biz.NewPermissionsUsecase(logger, permissionRepo, roleRepo, assignedRepo)
	require.NoError(t, err)

	check, err := biz.NewCheckPermissionsUsecase(assignedRepo, roleRepo, teamsRepo)
	require.NoError(t, err)

	sh := service.NewServiceHelper(check)

	service := service.NewPermissionsService(sh, uc, check)

	ctx := mockServerContext()

	req := &v1.CreatePermissionRequest{
		Id:          "some.group.permission",
		GroupId:     "some.group",
		AppId:       "app-id",
		Description: "testDesc",
		Fields:      []string{"field1", "field2"},
	}
	assignedRepo.EXPECT().CheckRoles(ctx, gomock.Any()).Return(mockAccess(), nil)
	roleRepo.EXPECT().ListRolesPermissions(ctx, gomock.Any()).Return(mockRolePermissions("admin.permission.create"), nil)

	reply, err := service.CreatePermission(ctx, req)
	require.Error(t, err)
	require.Equal(t, v1.ErrorBadRequest("empty name"), err)
	require.Nil(t, reply)
}

func TestPermissionsService_CreatePermissionEmptyGroup(t *testing.T) {
	logger := zap.NewZapLogger(true)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	permissionRepo := mock.NewMockPermissionRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	uc, err := biz.NewPermissionsUsecase(logger, permissionRepo, roleRepo, assignedRepo)
	require.NoError(t, err)

	check, err := biz.NewCheckPermissionsUsecase(assignedRepo, roleRepo, teamsRepo)
	require.NoError(t, err)

	sh := service.NewServiceHelper(check)

	service := service.NewPermissionsService(sh, uc, check)

	ctx := mockServerContext()

	req := &v1.CreatePermissionRequest{
		Id:          "some.group.permission",
		AppId:       "app-id",
		Name:        "testName",
		Description: "testDesc",
		Fields:      []string{"field1", "field2"},
	}
	assignedRepo.EXPECT().CheckRoles(ctx, gomock.Any()).Return(mockAccess(), nil)
	roleRepo.EXPECT().ListRolesPermissions(ctx, gomock.Any()).Return(mockRolePermissions("admin.permission.create"), nil)

	reply, err := service.CreatePermission(ctx, req)
	require.Error(t, err)
	require.Equal(t, v1.ErrorBadRequest("empty group"), err)
	require.Nil(t, reply)
}

func TestPermissionsService_CreatePermissionInvalidGroup(t *testing.T) {
	logger := zap.NewZapLogger(true)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	permissionRepo := mock.NewMockPermissionRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	uc, err := biz.NewPermissionsUsecase(logger, permissionRepo, roleRepo, assignedRepo)
	require.NoError(t, err)

	check, err := biz.NewCheckPermissionsUsecase(assignedRepo, roleRepo, teamsRepo)
	require.NoError(t, err)

	sh := service.NewServiceHelper(check)

	service := service.NewPermissionsService(sh, uc, check)

	ctx := mockServerContext()

	req := &v1.CreatePermissionRequest{
		Id:          "some.group.permission",
		GroupId:     "other.group",
		AppId:       "app-id",
		Name:        "testName",
		Description: "testDesc",
		Fields:      []string{"field1", "field2"},
	}
	assignedRepo.EXPECT().CheckRoles(ctx, gomock.Any()).Return(mockAccess(), nil)
	roleRepo.EXPECT().ListRolesPermissions(ctx, gomock.Any()).Return(mockRolePermissions("admin.permission.create"), nil)

	reply, err := service.CreatePermission(ctx, req)
	require.Error(t, err)
	require.Equal(t, v1.ErrorBadRequest("id must start with group id"), err)
	require.Nil(t, reply)
}

func TestPermissionsService_CreatePermissionEmptyAppID(t *testing.T) {
	logger := zap.NewZapLogger(true)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	permissionRepo := mock.NewMockPermissionRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	uc, err := biz.NewPermissionsUsecase(logger, permissionRepo, roleRepo, assignedRepo)
	require.NoError(t, err)

	check, err := biz.NewCheckPermissionsUsecase(assignedRepo, roleRepo, teamsRepo)
	require.NoError(t, err)

	sh := service.NewServiceHelper(check)

	service := service.NewPermissionsService(sh, uc, check)

	ctx := mockServerContext()

	req := &v1.CreatePermissionRequest{
		Id:          "some.group.permission",
		GroupId:     "some.group",
		Name:        "testName",
		Description: "testDesc",
		Fields:      []string{"field1", "field2"},
	}
	assignedRepo.EXPECT().CheckRoles(ctx, gomock.Any()).Return(mockAccess(), nil)
	roleRepo.EXPECT().ListRolesPermissions(ctx, gomock.Any()).Return(mockRolePermissions("admin.permission.create"), nil)

	reply, err := service.CreatePermission(ctx, req)
	require.Error(t, err)
	require.Equal(t, v1.ErrorBadRequest("empty app id"), err)
	require.Nil(t, reply)
}

// ------------------------ UpdatePermission ------------------------

func TestPermissionsService_UpdatePermission(t *testing.T) {
	logger := zap.NewZapLogger(true)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	permissionRepo := mock.NewMockPermissionRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	uc, err := biz.NewPermissionsUsecase(logger, permissionRepo, roleRepo, assignedRepo)
	require.NoError(t, err)

	check, err := biz.NewCheckPermissionsUsecase(assignedRepo, roleRepo, teamsRepo)
	require.NoError(t, err)

	sh := service.NewServiceHelper(check)

	service := service.NewPermissionsService(sh, uc, check)

	ctx := mockServerContext()

	req := &v1.UpdatePermissionRequest{
		PermissionId: "some.group.permission",
		Name:         "testNewName",
	}
	permission := &ent.Permission{
		ID:          req.GetPermissionId(),
		Name:        req.GetName(),
		GroupID:     "some.group",
		AppID:       "app-id",
		Description: "testDesc",
		Fields:      []string{"field1", "field2"},
	}
	assignedRepo.EXPECT().CheckRoles(ctx, gomock.Any()).Return(mockAccess(), nil)
	roleRepo.EXPECT().ListRolesPermissions(ctx, gomock.Any()).Return(mockRolePermissions("admin.permission.update"), nil)
	permissionRepo.EXPECT().UpdatePermission(ctx, req.GetPermissionId(), gomock.Any()).Return(permission, nil)

	expect := &v1.Permission{
		Id:          permission.ID,
		GroupId:     permission.GroupID,
		AppId:       permission.AppID,
		Name:        permission.Name,
		Description: permission.Description,
		Fields:      permission.Fields,
	}

	reply, err := service.UpdatePermission(ctx, req)
	require.NoError(t, err)
	require.Equal(t, expect, reply.GetPermission())
}

func TestPermissionsService_UpdatePermissionAccessDenied(t *testing.T) {
	logger := zap.NewZapLogger(true)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	permissionRepo := mock.NewMockPermissionRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	uc, err := biz.NewPermissionsUsecase(logger, permissionRepo, roleRepo, assignedRepo)
	require.NoError(t, err)

	check, err := biz.NewCheckPermissionsUsecase(assignedRepo, roleRepo, teamsRepo)
	require.NoError(t, err)

	sh := service.NewServiceHelper(check)

	service := service.NewPermissionsService(sh, uc, check)

	ctx := mockServerContext()

	req := &v1.UpdatePermissionRequest{
		PermissionId: "some.group.permission",
		Name:         "testNewName",
	}
	assignedRepo.EXPECT().CheckRoles(ctx, gomock.Any()).Return(mockAccess(), nil)
	roleRepo.EXPECT().ListRolesPermissions(ctx, gomock.Any()).Return(mockRolePermissions("wrong.permission"), nil)

	reply, err := service.UpdatePermission(ctx, req)
	require.Error(t, err)
	require.Equal(t, v1.ErrorForbidden("has no permission"), err)
	require.Nil(t, reply)
}

func TestPermissionsService_UpdatePermissionEmptyID(t *testing.T) {
	logger := zap.NewZapLogger(true)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	permissionRepo := mock.NewMockPermissionRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	uc, err := biz.NewPermissionsUsecase(logger, permissionRepo, roleRepo, assignedRepo)
	require.NoError(t, err)

	check, err := biz.NewCheckPermissionsUsecase(assignedRepo, roleRepo, teamsRepo)
	require.NoError(t, err)

	sh := service.NewServiceHelper(check)

	service := service.NewPermissionsService(sh, uc, check)

	ctx := mockServerContext()

	req := &v1.UpdatePermissionRequest{
		Name: "testNewName",
	}
	assignedRepo.EXPECT().CheckRoles(ctx, gomock.Any()).Return(mockAccess(), nil)
	roleRepo.EXPECT().ListRolesPermissions(ctx, gomock.Any()).Return(mockRolePermissions("admin.permission.update"), nil)

	reply, err := service.UpdatePermission(ctx, req)
	require.Error(t, err)
	require.Equal(t, v1.ErrorBadRequest("empty permission id"), err)
	require.Nil(t, reply)
}

// ------------------------ DeletePermission ------------------------

func TestPermissionsService_DeletePermission(t *testing.T) {
	logger := zap.NewZapLogger(true)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	permissionRepo := mock.NewMockPermissionRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	uc, err := biz.NewPermissionsUsecase(logger, permissionRepo, roleRepo, assignedRepo)
	require.NoError(t, err)

	check, err := biz.NewCheckPermissionsUsecase(assignedRepo, roleRepo, teamsRepo)
	require.NoError(t, err)

	sh := service.NewServiceHelper(check)

	service := service.NewPermissionsService(sh, uc, check)

	ctx := mockServerContext()

	req := &v1.PermissionRequest{
		PermissionId: "some.group.permission",
	}
	assignedRepo.EXPECT().CheckRoles(ctx, gomock.Any()).Return(mockAccess(), nil)
	roleRepo.EXPECT().ListRolesPermissions(ctx, gomock.Any()).Return(mockRolePermissions("admin.permission.delete"), nil)
	permissionRepo.EXPECT().DeletePermission(ctx, req.GetPermissionId()).Return(nil)

	_, err = service.DeletePermission(ctx, req)
	require.NoError(t, err)
}

func TestPermissionsService_DeletePermissionAccessDenied(t *testing.T) {
	logger := zap.NewZapLogger(true)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	permissionRepo := mock.NewMockPermissionRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	uc, err := biz.NewPermissionsUsecase(logger, permissionRepo, roleRepo, assignedRepo)
	require.NoError(t, err)

	check, err := biz.NewCheckPermissionsUsecase(assignedRepo, roleRepo, teamsRepo)
	require.NoError(t, err)

	sh := service.NewServiceHelper(check)

	service := service.NewPermissionsService(sh, uc, check)

	ctx := mockServerContext()

	req := &v1.PermissionRequest{
		PermissionId: "some.group.permission",
	}
	assignedRepo.EXPECT().CheckRoles(ctx, gomock.Any()).Return(mockAccess(), nil)
	roleRepo.EXPECT().ListRolesPermissions(ctx, gomock.Any()).Return(mockRolePermissions("wrong.permission"), nil)

	_, err = service.DeletePermission(ctx, req)
	require.Error(t, err)
	require.Equal(t, v1.ErrorForbidden("has no permission"), err)
}

// ------------------------ GetPermission ---------------------------

func TestPermissionsService_GetPermission(t *testing.T) {
	logger := zap.NewZapLogger(true)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	permissionRepo := mock.NewMockPermissionRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	uc, err := biz.NewPermissionsUsecase(logger, permissionRepo, roleRepo, assignedRepo)
	require.NoError(t, err)

	check, err := biz.NewCheckPermissionsUsecase(assignedRepo, roleRepo, teamsRepo)
	require.NoError(t, err)

	sh := service.NewServiceHelper(check)

	service := service.NewPermissionsService(sh, uc, check)

	ctx := mockServerContext()

	req := &v1.PermissionRequest{
		PermissionId: "some.group.permission",
	}
	permission := &ent.Permission{
		ID:          req.GetPermissionId(),
		GroupID:     "some.group",
		AppID:       "app-id",
		Name:        "testName",
		Description: "testDesc",
		Fields:      []string{"field1", "field2"},
	}
	assignedRepo.EXPECT().CheckRoles(ctx, gomock.Any()).Return(mockAccess(), nil)
	roleRepo.EXPECT().ListRolesPermissions(ctx, gomock.Any()).Return(mockRolePermissions("admin.permission.read"), nil)
	permissionRepo.EXPECT().GetPermissionByID(ctx, req.GetPermissionId()).Return(permission, nil)

	expect := &v1.Permission{
		Id:          permission.ID,
		GroupId:     permission.GroupID,
		AppId:       permission.AppID,
		Name:        permission.Name,
		Description: permission.Description,
		Fields:      permission.Fields,
	}

	reply, err := service.GetPermission(ctx, req)
	require.NoError(t, err)
	require.Equal(t, expect, reply.GetPermission())
}

func TestPermissionsService_GetPermissionAccessDenied(t *testing.T) {
	logger := zap.NewZapLogger(true)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	permissionRepo := mock.NewMockPermissionRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	uc, err := biz.NewPermissionsUsecase(logger, permissionRepo, roleRepo, assignedRepo)
	require.NoError(t, err)

	check, err := biz.NewCheckPermissionsUsecase(assignedRepo, roleRepo, teamsRepo)
	require.NoError(t, err)

	sh := service.NewServiceHelper(check)

	service := service.NewPermissionsService(sh, uc, check)

	ctx := mockServerContext()

	req := &v1.PermissionRequest{
		PermissionId: "some.group.permission",
	}
	assignedRepo.EXPECT().CheckRoles(ctx, gomock.Any()).Return(mockAccess(), nil)
	roleRepo.EXPECT().ListRolesPermissions(ctx, gomock.Any()).Return(mockRolePermissions("wrong.permission"), nil)

	reply, err := service.GetPermission(ctx, req)
	require.Error(t, err)
	require.Equal(t, v1.ErrorForbidden("has no permission"), err)
	require.Nil(t, reply)
}

// ------------------------ ListPermissions -------------------------

func TestPermissionsService_ListPermissions(t *testing.T) {
	logger := zap.NewZapLogger(true)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	permissionRepo := mock.NewMockPermissionRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	uc, err := biz.NewPermissionsUsecase(logger, permissionRepo, roleRepo, assignedRepo)
	require.NoError(t, err)

	check, err := biz.NewCheckPermissionsUsecase(assignedRepo, roleRepo, teamsRepo)
	require.NoError(t, err)

	sh := service.NewServiceHelper(check)

	service := service.NewPermissionsService(sh, uc, check)

	ctx := mockServerContext()

	req := &v1.ListPermissionsRequest{
		AppsIds: []string{"app-id"},
	}
	filter := data.FilterPermissions{
		AppsIDs: req.GetAppsIds(),
	}
	groups := []*ent.PermissionGroup{
		{
			ID:    "some.group",
			Name:  "Some Group",
			AppID: "app-id",
			Edges: ent.PermissionGroupEdges{
				Permissions: []*ent.Permission{
					{
						ID:          "some.group.permission",
						GroupID:     "some.group",
						AppID:       "app-id",
						Name:        "testName",
						Description: "testDesc",
						Fields:      []string{"field1", "field2"},
					},
					{
						ID:          "some.group.permission2",
						GroupID:     "some.group",
						AppID:       "app-id",
						Name:        "testName2",
						Description: "testDesc2",
					},
				},
			},
		},
		{
			ID:    "some.group2",
			Name:  "Some Group 2",
			AppID: "app-id",
			Edges: ent.PermissionGroupEdges{
				Permissions: []*ent.Permission{
					{
						ID:          "some.group2.permission",
						GroupID:     "some.group2",
						AppID:       "app-id",
						Name:        "testName3",
						Description: "testDesc3",
					},
					{
						ID:          "some.group2.permission2",
						GroupID:     "some.group2",
						AppID:       "app-id",
						Name:        "testName4",
						Description: "testDesc4",
					},
				},
			},
		},
	}
	assignedRepo.EXPECT().CheckRoles(ctx, gomock.Any()).Return(mockAccess(), nil)
	roleRepo.EXPECT().ListRolesPermissions(ctx, gomock.Any()).Return(mockRolePermissions("admin.permission.read"), nil)
	permissionRepo.EXPECT().GetGroupedPermissions(ctx, filter).Return(groups, nil)
	assignedRepo.EXPECT().ListAssignedRoles(ctx, gomock.Any()).Return(mockAccess(), nil)
	roleRepo.EXPECT().ListRolesPermissions(ctx, gomock.Any()).Return(nil, nil)

	expectGroup := &v1.Group{
		Id:    "some.group",
		AppId: "app-id",
		Name:  "Some Group",
		Permissions: []*v1.Permission{
			{
				Id:          "some.group.permission",
				GroupId:     "some.group",
				AppId:       "app-id",
				Name:        "testName",
				Description: "testDesc",
				Fields:      []string{"field1", "field2"},
			},
			{
				Id:          "some.group.permission2",
				GroupId:     "some.group",
				AppId:       "app-id",
				Name:        "testName2",
				Description: "testDesc2",
			},
		},
	}

	reply, err := service.ListPermissions(ctx, req)
	require.NoError(t, err)
	require.Len(t, reply.GetGroups(), 2)
	require.Equal(t, expectGroup, reply.GetGroups()[0])
}

func TestPermissionsService_ListPermissionsAccessDenied(t *testing.T) {
	logger := zap.NewZapLogger(true)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	permissionRepo := mock.NewMockPermissionRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	uc, err := biz.NewPermissionsUsecase(logger, permissionRepo, roleRepo, assignedRepo)
	require.NoError(t, err)

	check, err := biz.NewCheckPermissionsUsecase(assignedRepo, roleRepo, teamsRepo)
	require.NoError(t, err)

	sh := service.NewServiceHelper(check)

	service := service.NewPermissionsService(sh, uc, check)

	ctx := mockServerContext()

	req := &v1.ListPermissionsRequest{
		AppsIds: []string{"app-id"},
	}
	assignedRepo.EXPECT().CheckRoles(ctx, gomock.Any()).Return(mockAccess(), nil)
	roleRepo.EXPECT().ListRolesPermissions(ctx, gomock.Any()).Return(mockRolePermissions("wrong.permission"), nil)

	reply, err := service.ListPermissions(ctx, req)
	require.Error(t, err)
	require.Equal(t, v1.ErrorForbidden("has no permission"), err)
	require.Nil(t, reply)
}
