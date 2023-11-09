package biz

import "github.com/google/wire"

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	NewQueueManager,
	NewPermissionUsecase,
	NewRolesUsecase,
	NewTeamIdentityUsecase,
	NewTeamsUsecase,
)
