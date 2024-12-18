package biz

import (
	"github.com/google/wire"
)

//nolint:gochecknoglobals // this global variables are required for queues
var (
	QueueRoleAssign   = "role_assign"
	QueueRoleUnassign = "role_unassign"
)

// ProviderSet is biz providers.
//
//nolint:gochecknoglobals // this global variable is required for wire
var ProviderSet = wire.NewSet(
	NewPermissionsUsecase,
	NewRolesUsecase,
	NewAssignedRolesUsecase,
	NewCheckPermissionsUsecase,
	NewTeamsUsecase,
	NewPaidContent,
)
