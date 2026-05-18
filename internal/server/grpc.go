package server

import (
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	v1 "github.com/makesalekz/rbac/api/rbac/v1"
	"github.com/makesalekz/rbac/internal/conf"
	"github.com/makesalekz/rbac/internal/service"
	"github.com/makesalekz/utils/v1/middlewares/metrics"
	"github.com/makesalekz/utils/v2/jwt"
	"github.com/makesalekz/utils/v2/middlewares/auth"
	u_tracing "github.com/makesalekz/utils/v2/tracing"

	prom "github.com/go-kratos/kratos/contrib/metrics/prometheus/v2"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(
	c *conf.Bootstrap,
	jwtp jwt.IJwtProcessor,
	tracer *u_tracing.Tracer,
	roleSrvc *service.RolesService,
	permissionsSrvc *service.PermissionsService,
	teamSrvc *service.TeamsService,
	teamIdentityRolesSrvc *service.AssignsService,
	checkPermissionSrvc *service.CheckPermissionsService,
) *grpc.Server {
	err := tracer.Initialize()
	if err != nil {
		panic(err)
	}

	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			metadata.Server(),
			tracing.Server(),
			auth.Server(jwtp),
			metrics.Server(
				metrics.WithSeconds(prom.NewHistogram(_metricSeconds)),
				metrics.WithRequests(prom.NewCounter(_metricRequests)),
				metrics.WithGauge(prom.NewGauge(_activeRequests)),
			),
		),
	}
	if c.GetServer().GetGrpc().GetNetwork() != "" {
		opts = append(opts, grpc.Network(c.GetServer().GetGrpc().GetNetwork()))
	}
	if c.GetServer().GetGrpc().GetAddr() != "" {
		opts = append(opts, grpc.Address(c.GetServer().GetGrpc().GetAddr()))
	}
	if c.GetServer().GetGrpc().GetTimeout() != nil {
		opts = append(opts, grpc.Timeout(c.GetServer().GetGrpc().GetTimeout().AsDuration()))
	}
	srv := grpc.NewServer(opts...)

	v1.RegisterRolesServer(srv, roleSrvc)
	v1.RegisterPermissionsServer(srv, permissionsSrvc)
	v1.RegisterTeamsServer(srv, teamSrvc)
	v1.RegisterAssignsServer(srv, teamIdentityRolesSrvc)
	v1.RegisterCheckPermissionsServer(srv, checkPermissionSrvc)

	return srv
}
