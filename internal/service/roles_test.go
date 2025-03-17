package service_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	v1 "gitlab.calendaria.team/services/rbac/api/rbac/v1"
	"gitlab.calendaria.team/services/rbac/ent"
	"gitlab.calendaria.team/services/rbac/internal/biz"
	"gitlab.calendaria.team/services/rbac/internal/data"
	"gitlab.calendaria.team/services/rbac/internal/data/mock"
	"gitlab.calendaria.team/services/rbac/internal/service"
	u_nats "gitlab.calendaria.team/services/utils/v2/nats"
	u_nats_mock "gitlab.calendaria.team/services/utils/v2/nats/mock"
	u_zap "gitlab.calendaria.team/services/utils/v2/zap"

	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func mockTenantServerContext(tenantID int64) context.Context {
	md := metadata.Metadata{
		"x-md-global-tenant-id":  []string{strconv.FormatInt(tenantID, 10)},
		"x-md-global-actor-id":   []string{"1234567"},
		"x-md-global-identities": []string{"identity1", "identity2"},
		"x-md-global-app-id":     []string{"app-id"},
	}
	return metadata.NewServerContext(context.Background(), md)
}

func createRolesService(
	t *testing.T,
	qm u_nats.IQueueManager,
	permissionRepo data.PermissionRepo,
	assignedRepo data.AssignedRolesRepo,
	roleRepo data.RoleRepo,
	teamsRepo data.TeamsRepo,
) *service.RolesService {
	logger := u_zap.NewZapLogger(true)

	ru, err := biz.NewRolesUsecase(logger, roleRepo)
	require.NoError(t, err)

	pu, err := biz.NewPermissionsUsecase(logger, permissionRepo, roleRepo, assignedRepo)
	require.NoError(t, err)

	au, err := biz.NewAssignedRolesUsecase(logger, assignedRepo, roleRepo, teamsRepo, qm)
	require.NoError(t, err)

	check, err := biz.NewCheckPermissionsUsecase(assignedRepo, roleRepo, teamsRepo)
	require.NoError(t, err)

	sh := service.NewServiceHelper(check)

	return service.NewRolesService(sh, ru, pu, au, check)
}

// ------------------ CreateRole ------------------------

func TestRolesService_CreateRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	qm := u_nats_mock.NewMockIQueueManager(ctrl)
	qm.EXPECT().AddConsumer(biz.QueueRoleAssignHandler, gomock.Any())
	qm.EXPECT().AddConsumer(biz.QueueRoleUnassignHandler, gomock.Any())
	permissionRepo := mock.NewMockPermissionRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	service := createRolesService(t, qm, permissionRepo, assignedRepo, roleRepo, teamsRepo)

	tenantID := int64(1234)
	ctx := mockTenantServerContext(tenantID)

	req := &v1.CreateRoleRequest{
		Name: "testName",
	}
	dto := data.CreateRoleDto{
		TenantID: tenantID,
		Name:     req.GetName(),
	}
	role := &ent.Role{
		Name:      req.GetName(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	roleRepo.EXPECT().CreateRole(ctx, dto).Return(role, nil)

	expect := &v1.Role{
		Name:      role.Name,
		CreatedAt: role.CreatedAt.String(),
		UpdatedAt: role.UpdatedAt.String(),
	}

	reply, err := service.CreateRole(ctx, req)
	require.NoError(t, err)
	require.Equal(t, expect, reply.GetRole())
}

func TestRolesService_CreateRoleEmptyName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	qm := u_nats_mock.NewMockIQueueManager(ctrl)
	qm.EXPECT().AddConsumer(biz.QueueRoleAssignHandler, gomock.Any())
	qm.EXPECT().AddConsumer(biz.QueueRoleUnassignHandler, gomock.Any())
	permissionRepo := mock.NewMockPermissionRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	service := createRolesService(t, qm, permissionRepo, assignedRepo, roleRepo, teamsRepo)

	ctx := mockServerContext()
	req := &v1.CreateRoleRequest{}

	reply, err := service.CreateRole(ctx, req)
	require.Error(t, err)
	require.Equal(t, v1.ErrorBadRequest("empty name"), err)
	require.Nil(t, reply)
}

// ------------------ UpdateRole ------------------------

func TestRolesService_UpdateRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	qm := u_nats_mock.NewMockIQueueManager(ctrl)
	qm.EXPECT().AddConsumer(biz.QueueRoleAssignHandler, gomock.Any())
	qm.EXPECT().AddConsumer(biz.QueueRoleUnassignHandler, gomock.Any())
	permissionRepo := mock.NewMockPermissionRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	service := createRolesService(t, qm, permissionRepo, assignedRepo, roleRepo, teamsRepo)

	tenantID := int64(123456)
	ctx := mockTenantServerContext(tenantID)

	req := &v1.UpdateRoleRequest{
		RoleId: 1,
		Name:   "testName",
	}
	dto := data.UpdateRoleDto{
		Name: req.GetName(),
	}
	role := &ent.Role{
		ID:        req.GetRoleId(),
		Name:      req.GetName(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	roleRepo.EXPECT().UpdateRole(ctx, tenantID, req.GetRoleId(), dto).Return(role, nil)

	expect := &v1.Role{
		Id:        role.ID,
		Name:      role.Name,
		CreatedAt: role.CreatedAt.String(),
		UpdatedAt: role.UpdatedAt.String(),
	}

	reply, err := service.UpdateRole(ctx, req)
	require.NoError(t, err)
	require.Equal(t, expect, reply.GetRole())
}

func TestRolesService_UpdateRoleEmptyId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	qm := u_nats_mock.NewMockIQueueManager(ctrl)
	qm.EXPECT().AddConsumer(biz.QueueRoleAssignHandler, gomock.Any())
	qm.EXPECT().AddConsumer(biz.QueueRoleUnassignHandler, gomock.Any())
	permissionRepo := mock.NewMockPermissionRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	service := createRolesService(t, qm, permissionRepo, assignedRepo, roleRepo, teamsRepo)

	tenantID := int64(1234)
	ctx := mockTenantServerContext(tenantID)

	req := &v1.UpdateRoleRequest{
		Name: "testName",
	}

	reply, err := service.UpdateRole(ctx, req)
	require.Error(t, err)
	require.Equal(t, v1.ErrorBadRequest("empty role id"), err)
	require.Nil(t, reply)
}

// ------------------ DeleteRole ------------------------

func TestRolesService_DeleteRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	qm := u_nats_mock.NewMockIQueueManager(ctrl)
	qm.EXPECT().AddConsumer(biz.QueueRoleAssignHandler, gomock.Any())
	qm.EXPECT().AddConsumer(biz.QueueRoleUnassignHandler, gomock.Any())
	permissionRepo := mock.NewMockPermissionRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	service := createRolesService(t, qm, permissionRepo, assignedRepo, roleRepo, teamsRepo)

	tenantID := int64(1234)
	ctx := mockTenantServerContext(tenantID)

	req := &v1.RoleRequest{
		RoleId: 1,
	}
	roleRepo.EXPECT().DeleteRole(ctx, tenantID, req.GetRoleId()).Return(nil)

	_, err := service.DeleteRole(ctx, req)
	require.NoError(t, err)
}

func TestRolesService_DeleteRoleEmptyId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	qm := u_nats_mock.NewMockIQueueManager(ctrl)
	qm.EXPECT().AddConsumer(biz.QueueRoleAssignHandler, gomock.Any())
	qm.EXPECT().AddConsumer(biz.QueueRoleUnassignHandler, gomock.Any())
	permissionRepo := mock.NewMockPermissionRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	service := createRolesService(t, qm, permissionRepo, assignedRepo, roleRepo, teamsRepo)

	tenantID := int64(1234)
	ctx := mockTenantServerContext(tenantID)

	req := &v1.RoleRequest{}

	reply, err := service.DeleteRole(ctx, req)
	require.Error(t, err)
	require.Equal(t, v1.ErrorBadRequest("empty role id"), err)
	require.Nil(t, reply)
}

// ------------------ GetRole ---------------------------

func TestRolesService_GetRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	qm := u_nats_mock.NewMockIQueueManager(ctrl)
	qm.EXPECT().AddConsumer(biz.QueueRoleAssignHandler, gomock.Any())
	qm.EXPECT().AddConsumer(biz.QueueRoleUnassignHandler, gomock.Any())
	permissionRepo := mock.NewMockPermissionRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	service := createRolesService(t, qm, permissionRepo, assignedRepo, roleRepo, teamsRepo)

	tenantID := int64(1234)
	ctx := mockTenantServerContext(tenantID)

	req := &v1.RoleRequest{
		RoleId: 1,
	}
	role := &ent.Role{
		ID:          req.GetRoleId(),
		Name:        "testName",
		Description: "testDescription",
		IsSystem:    false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
	}
	roleRepo.EXPECT().GetRoleByID(ctx, tenantID, req.GetRoleId()).Return(role, nil)

	expect := &v1.Role{
		Id:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		IsSystem:    role.IsSystem,
		CreatedAt:   role.CreatedAt.String(),
		UpdatedAt:   role.UpdatedAt.String(),
		DeletedAt:   "",
	}

	reply, err := service.GetRole(ctx, req)
	require.NoError(t, err)
	require.Equal(t, expect, reply.GetRole())
}

func TestRolesService_GetRoleEmptyId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	qm := u_nats_mock.NewMockIQueueManager(ctrl)
	qm.EXPECT().AddConsumer(biz.QueueRoleAssignHandler, gomock.Any())
	qm.EXPECT().AddConsumer(biz.QueueRoleUnassignHandler, gomock.Any())
	permissionRepo := mock.NewMockPermissionRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	service := createRolesService(t, qm, permissionRepo, assignedRepo, roleRepo, teamsRepo)

	tenantID := int64(1234)
	ctx := mockTenantServerContext(tenantID)

	req := &v1.RoleRequest{}

	reply, err := service.GetRole(ctx, req)
	require.Error(t, err)
	require.Equal(t, v1.ErrorBadRequest("empty role id"), err)
	require.Nil(t, reply)
}

// ------------------ ListRoles -------------------------

func TestRolesService_ListRoles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	qm := u_nats_mock.NewMockIQueueManager(ctrl)
	qm.EXPECT().AddConsumer(biz.QueueRoleAssignHandler, gomock.Any())
	qm.EXPECT().AddConsumer(biz.QueueRoleUnassignHandler, gomock.Any())
	permissionRepo := mock.NewMockPermissionRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	rolesService := createRolesService(t, qm, permissionRepo, assignedRepo, roleRepo, teamsRepo)

	tenantID := int64(1234)
	ctx := mockTenantServerContext(tenantID)

	req := &v1.ListRolesRequest{
		Search: "testSearch",
	}
	role := &ent.Role{
		ID:          1,
		Name:        "testName",
		Description: "testDescription",
		IsSystem:    false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
	}
	roles := []*ent.Role{
		role,
		{
			ID:          2,
			Name:        "testName2",
			Description: "testDescription2",
			IsSystem:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			DeletedAt:   nil,
		},
	}
	roleRepo.EXPECT().GetRolesList(ctx, tenantID, req.GetSearch(), false).Return(roles, nil)

	expect := &v1.Role{
		Id:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		IsSystem:    role.IsSystem,
		CreatedAt:   role.CreatedAt.String(),
		UpdatedAt:   role.UpdatedAt.String(),
		DeletedAt:   "",
	}

	reply, err := rolesService.ListRoles(ctx, req)
	require.NoError(t, err)
	require.Len(t, reply.GetRoles(), 2)
	require.Equal(t, expect, reply.GetRoles()[0])
}

// ------------------ AddPermissionToRole ---------------

func TestRolesService_AddPermissionToRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	qm := u_nats_mock.NewMockIQueueManager(ctrl)
	qm.EXPECT().AddConsumer(biz.QueueRoleAssignHandler, gomock.Any())
	qm.EXPECT().AddConsumer(biz.QueueRoleUnassignHandler, gomock.Any())
	permissionRepo := mock.NewMockPermissionRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	service := createRolesService(t, qm, permissionRepo, assignedRepo, roleRepo, teamsRepo)

	tenantID := int64(1234)
	ctx := mockTenantServerContext(tenantID)

	req := &v1.AddPermissionToRoleRequest{
		RoleId:       1,
		PermissionId: "some.group.permission",
	}
	dto := data.CreateRolePermissionDto{}
	permission := &ent.Permission{
		ID:          req.GetPermissionId(),
		GroupID:     "some.group",
		AppID:       "app-id",
		Name:        "testName",
		Description: "testDesc",
		Fields:      []string{"field1", "field2"},
	}
	permissionRepo.EXPECT().GetPermissionByID(ctx, req.GetPermissionId()).Return(permission, nil)
	roleRepo.EXPECT().SetRolePermission(ctx, tenantID, req.GetRoleId(), permission, dto).Return(nil)

	_, err := service.AddPermissionToRole(ctx, req)
	require.NoError(t, err)
}

func TestRolesService_AddPermissionToRoleEmptyRoleId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	qm := u_nats_mock.NewMockIQueueManager(ctrl)
	qm.EXPECT().AddConsumer(biz.QueueRoleAssignHandler, gomock.Any())
	qm.EXPECT().AddConsumer(biz.QueueRoleUnassignHandler, gomock.Any())
	permissionRepo := mock.NewMockPermissionRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	service := createRolesService(t, qm, permissionRepo, assignedRepo, roleRepo, teamsRepo)

	tenantID := int64(1234)
	ctx := mockTenantServerContext(tenantID)

	req := &v1.AddPermissionToRoleRequest{
		PermissionId: "some.group.permission",
	}

	reply, err := service.AddPermissionToRole(ctx, req)
	require.Error(t, err)
	require.Equal(t, v1.ErrorBadRequest("empty role id"), err)
	require.Nil(t, reply)
}

func TestRolesService_AddPermissionToRoleEmptyPermissionId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	qm := u_nats_mock.NewMockIQueueManager(ctrl)
	qm.EXPECT().AddConsumer(biz.QueueRoleAssignHandler, gomock.Any())
	qm.EXPECT().AddConsumer(biz.QueueRoleUnassignHandler, gomock.Any())
	permissionRepo := mock.NewMockPermissionRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	service := createRolesService(t, qm, permissionRepo, assignedRepo, roleRepo, teamsRepo)

	tenantID := int64(1234)
	ctx := mockTenantServerContext(tenantID)

	req := &v1.AddPermissionToRoleRequest{
		RoleId: 1,
	}

	reply, err := service.AddPermissionToRole(ctx, req)
	require.Error(t, err)
	require.Equal(t, v1.ErrorBadRequest("empty permission id"), err)
	require.Nil(t, reply)
}

// ------------------ RemovePermissionFromRole ----------

func TestRolesService_RemovePermissionFromRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	qm := u_nats_mock.NewMockIQueueManager(ctrl)
	qm.EXPECT().AddConsumer(biz.QueueRoleAssignHandler, gomock.Any())
	qm.EXPECT().AddConsumer(biz.QueueRoleUnassignHandler, gomock.Any())
	permissionRepo := mock.NewMockPermissionRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	service := createRolesService(t, qm, permissionRepo, assignedRepo, roleRepo, teamsRepo)

	tenantID := int64(1234)
	ctx := mockTenantServerContext(tenantID)

	req := &v1.RemovePermissionFromRoleRequest{
		RoleId:       1,
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
	permissionRepo.EXPECT().GetPermissionByID(ctx, req.GetPermissionId()).Return(permission, nil)
	roleRepo.EXPECT().RemovePermissionFromRole(ctx, tenantID, req.GetRoleId(), permission).Return(nil)

	_, err := service.RemovePermissionFromRole(ctx, req)
	require.NoError(t, err)
}

func TestRolesService_RemovePermissionFromRoleEmptyRoleId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	qm := u_nats_mock.NewMockIQueueManager(ctrl)
	qm.EXPECT().AddConsumer(biz.QueueRoleAssignHandler, gomock.Any())
	qm.EXPECT().AddConsumer(biz.QueueRoleUnassignHandler, gomock.Any())
	permissionRepo := mock.NewMockPermissionRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	service := createRolesService(t, qm, permissionRepo, assignedRepo, roleRepo, teamsRepo)

	tenantID := int64(1234)
	ctx := mockTenantServerContext(tenantID)

	req := &v1.RemovePermissionFromRoleRequest{
		PermissionId: "some.group.permission",
	}

	reply, err := service.RemovePermissionFromRole(ctx, req)
	require.Error(t, err)
	require.Equal(t, v1.ErrorBadRequest("empty role id"), err)
	require.Nil(t, reply)
}

func TestRolesService_RemovePermissionFromRoleEmptyPermissionId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	qm := u_nats_mock.NewMockIQueueManager(ctrl)
	qm.EXPECT().AddConsumer(biz.QueueRoleAssignHandler, gomock.Any())
	qm.EXPECT().AddConsumer(biz.QueueRoleUnassignHandler, gomock.Any())
	permissionRepo := mock.NewMockPermissionRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	service := createRolesService(t, qm, permissionRepo, assignedRepo, roleRepo, teamsRepo)

	tenantID := int64(1234)
	ctx := mockTenantServerContext(tenantID)

	req := &v1.RemovePermissionFromRoleRequest{
		RoleId: 1,
	}

	reply, err := service.RemovePermissionFromRole(ctx, req)
	require.Error(t, err)
	require.Equal(t, v1.ErrorBadRequest("empty permission id"), err)
	require.Nil(t, reply)
}

// ------------------ ListRolePermissions --------------

func TestRolesService_ListRolePermissions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	qm := u_nats_mock.NewMockIQueueManager(ctrl)
	qm.EXPECT().AddConsumer(biz.QueueRoleAssignHandler, gomock.Any())
	qm.EXPECT().AddConsumer(biz.QueueRoleUnassignHandler, gomock.Any())
	permissionRepo := mock.NewMockPermissionRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	service := createRolesService(t, qm, permissionRepo, assignedRepo, roleRepo, teamsRepo)

	tenantID := int64(1234)
	appID := "app-id"

	ctx := mockTenantServerContext(tenantID)

	req := &v1.RoleRequest{
		RoleId: 1,
	}
	permission := &ent.RolePermission{
		PermissionID: "some.group.permission",
		Deny:         true,
		Fields:       []string{"field1", "field2"},
	}
	roles := []*ent.RolePermission{
		permission,
		{
			PermissionID: "some.group.permission2",
		},
	}
	roleRepo.EXPECT().ListRolePermissions(
		ctx, tenantID, req.GetRoleId(),
		[]string{appID, "common", "admin"},
	).Return(roles, nil)

	expect := &v1.RolePermission{
		Id:     permission.PermissionID,
		Deny:   permission.Deny,
		Fields: permission.Fields,
	}

	reply, err := service.ListRolePermissions(ctx, req)
	require.NoError(t, err)
	require.Len(t, reply.GetPermissions(), 2)
	require.Equal(t, expect, reply.GetPermissions()[0])
}
