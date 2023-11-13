package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	jwtv4 "github.com/golang-jwt/jwt/v4"
	v1 "gitlab.calendaria.team/services/rbac/api/rbac/v1"
	"gitlab.calendaria.team/services/rbac/internal/conf"
	"gitlab.calendaria.team/services/rbac/internal/data"
	"gitlab.calendaria.team/services/rbac/internal/service"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Bootstrap, logger log.Logger, jwtp *data.JwtProcessor,
	roleSrvc *service.RolesService,
	permissionsSrvc *service.PermissionsService,
	teamSrvc *service.TeamsService,
	teamIdentityRolesSrvc *service.TeamIdentityRoleService,
	checkPermissionSrvc *service.CheckPermissionsService,
) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			metadata.Server(),
			jwt.Server(func(token *jwtv4.Token) (interface{}, error) {
				return jwtp.GetSecret(), nil
			}, jwt.WithSigningMethod(jwtv4.SigningMethodHS256)),
		),
	}
	if c.Server.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Server.Grpc.Network))
	}
	if c.Server.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Server.Grpc.Addr))
	}
	if c.Server.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Server.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)

	v1.RegisterRolesServer(srv, roleSrvc)
	v1.RegisterPermissionsServer(srv, permissionsSrvc)
	v1.RegisterTeamsServer(srv, teamSrvc)
	v1.RegisterTeamIdentityRoleServer(srv, teamIdentityRolesSrvc)
	v1.RegisterCheckPermissionsServer(srv, checkPermissionSrvc)

	return srv
}
