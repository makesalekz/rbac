package service_test

import (
	"context"
	"strconv"
	"testing"

	v1 "gitlab.calendaria.team/services/rbac/api/rbac/v1"
	"gitlab.calendaria.team/services/rbac/ent"
	"gitlab.calendaria.team/services/rbac/internal/biz"
	"gitlab.calendaria.team/services/rbac/internal/data"
	"gitlab.calendaria.team/services/rbac/internal/data/mock"
	"gitlab.calendaria.team/services/rbac/internal/service"

	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgtype"
	"github.com/stretchr/testify/require"
)

func TestRolesService_CheckPermissions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	uc, err := biz.NewCheckPermissionsUsecase(assignedRepo, roleRepo, teamsRepo)
	require.NoError(t, err)

	service := service.NewCheckPermissionsService(uc)

	tenantID := int64(1234)
	ctx := mockTenantServerContext(tenantID)

	req := &v1.CheckPermissionsRequest{
		TenantId:    tenantID,
		Identities:  []string{"identity1"},
		Permissions: []string{"permission1", "permission2"},
	}
	dtoAssigned := data.ListRolesDto{
		TenantID:    tenantID,
		IdentityIDs: req.GetIdentities(),
	}
	assignedRoles := []*ent.ResourceAccess{
		{RoleID: 1},
		{RoleID: 2},
	}
	assignedRepo.EXPECT().CheckRoles(ctx, dtoAssigned).Return(assignedRoles, nil)

	dtoPermissions := data.FilterRolePermissions{
		TenantID:    tenantID,
		RoleIDs:     []int64{1, 2},
		Permissions: req.GetPermissions(),
	}
	rp1 := &ent.RolePermission{
		ID:           1,
		RoleID:       1,
		PermissionID: "permission1",
		Fields:       []string{"field1", "field2"},
	}
	rp2 := &ent.RolePermission{
		ID:           2,
		RoleID:       2,
		PermissionID: "permission2",
		Fields:       []string{},
	}
	rolesPermissions := []*ent.RolePermission{rp1, rp2}
	roleRepo.EXPECT().ListRolesPermissions(ctx, dtoPermissions).Return(rolesPermissions, nil)

	expect := map[string]*v1.ListOfFields{
		"permission1": {
			Fields: rp1.Fields,
		},
		"permission2": {
			Fields: rp2.Fields,
		},
	}

	reply, err := service.CheckPermissions(ctx, req)
	require.NoError(t, err)
	require.Len(t, reply.GetPermissions(), 2)
	require.Equal(t, expect, reply.GetPermissions())
}

func TestRolesService_CheckPermissionsResources(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	uc, err := biz.NewCheckPermissionsUsecase(assignedRepo, roleRepo, teamsRepo)
	require.NoError(t, err)

	service := service.NewCheckPermissionsService(uc)

	tenantID := int64(1234)
	ctx := mockTenantServerContext(tenantID)

	parentIDs := []pgtype.Int8{
		{Int: 11, Status: pgtype.Present},
		{Int: 22, Status: pgtype.Present},
	}
	team := &ent.Team{
		ID:         111,
		TenantID:   tenantID,
		ParentsIds: &pgtype.Int8Array{Elements: parentIDs, Status: pgtype.Present},
	}
	resources := []*v1.Resource{
		{Id: 333, Type: "project"},
		{Id: team.ID, Type: "team"},
	}
	req := &v1.CheckPermissionsRequest{
		TenantId:    tenantID,
		Identities:  []string{"identity1"},
		Permissions: []string{"permission1", "permission2"},
		Resources:   resources,
	}

	teamsRepo.EXPECT().GetTeams(ctx, tenantID, []int64{team.ID}).Return([]*ent.Team{team}, nil)

	resources = append(resources, &v1.Resource{Id: 11, Type: "team"}, &v1.Resource{Id: 22, Type: "team"})

	dtoAssigned := data.ListRolesDto{
		TenantID:    tenantID,
		IdentityIDs: req.GetIdentities(),
		Resources:   resources,
	}
	assignedRoles := []*ent.ResourceAccess{
		{RoleID: 1},
		{RoleID: 2},
	}
	assignedRepo.EXPECT().CheckRoles(ctx, dtoAssigned).Return(assignedRoles, nil)

	dtoPermissions := data.FilterRolePermissions{
		TenantID:    tenantID,
		RoleIDs:     []int64{1, 2},
		Permissions: req.GetPermissions(),
	}
	rp1 := &ent.RolePermission{
		ID:           1,
		RoleID:       1,
		PermissionID: "permission1",
		Fields:       []string{"field1", "field2"},
	}
	rp2 := &ent.RolePermission{
		ID:           2,
		RoleID:       2,
		PermissionID: "permission2",
		Fields:       []string{},
	}
	rolesPermissions := []*ent.RolePermission{rp1, rp2}
	roleRepo.EXPECT().ListRolesPermissions(ctx, dtoPermissions).Return(rolesPermissions, nil)

	expect := map[string]*v1.ListOfFields{
		"permission1": {
			Fields: rp1.Fields,
		},
		"permission2": {
			Fields: rp2.Fields,
		},
	}

	reply, err := service.CheckPermissions(ctx, req)
	require.NoError(t, err)
	require.Len(t, reply.GetPermissions(), 2)
	require.Equal(t, expect, reply.GetPermissions())
}

func TestRolesService_CheckPermissionsMeta(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	roleRepo := mock.NewMockRoleRepo(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	teamsRepo := mock.NewMockTeamsRepo(ctrl)

	uc, err := biz.NewCheckPermissionsUsecase(assignedRepo, roleRepo, teamsRepo)
	require.NoError(t, err)

	service := service.NewCheckPermissionsService(uc)

	tenantID := int64(1234)
	identities := []string{"identity1"}
	md := metadata.Metadata{
		"x-md-global-tenant-id":  []string{strconv.FormatInt(tenantID, 10)},
		"x-md-global-actor-id":   []string{"1234567"},
		"x-md-global-identities": identities,
		"x-md-global-app-id":     []string{"app-id"},
	}
	ctx := metadata.NewServerContext(context.Background(), md)

	req := &v1.CheckPermissionsRequest{
		Permissions: []string{"permission1", "permission2"},
	}
	dtoAssigned := data.ListRolesDto{
		TenantID:    tenantID,
		IdentityIDs: identities,
	}
	assignedRoles := []*ent.ResourceAccess{
		{RoleID: 1},
		{RoleID: 2},
	}
	assignedRepo.EXPECT().CheckRoles(ctx, dtoAssigned).Return(assignedRoles, nil)

	dtoPermissions := data.FilterRolePermissions{
		TenantID:    tenantID,
		RoleIDs:     []int64{1, 2},
		Permissions: req.GetPermissions(),
	}
	rp1 := &ent.RolePermission{
		ID:           1,
		RoleID:       1,
		PermissionID: "permission1",
		Fields:       []string{"field1", "field2"},
	}
	rp2 := &ent.RolePermission{
		ID:           2,
		RoleID:       2,
		PermissionID: "permission2",
		Fields:       []string{},
	}
	rolesPermissions := []*ent.RolePermission{rp1, rp2}
	roleRepo.EXPECT().ListRolesPermissions(ctx, dtoPermissions).Return(rolesPermissions, nil)

	expect := map[string]*v1.ListOfFields{
		"permission1": {
			Fields: rp1.Fields,
		},
		"permission2": {
			Fields: rp2.Fields,
		},
	}

	reply, err := service.CheckPermissions(ctx, req)
	require.NoError(t, err)
	require.Len(t, reply.GetPermissions(), 2)
	require.Equal(t, expect, reply.GetPermissions())
}
