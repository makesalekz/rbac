package server

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"gitlab.calendaria.team/services/rbac/internal/biz"
)

type BackgroundServer struct {
	log *log.Helper
	pc  *biz.PaidContent
}

// NewCronServer
func NewBackgroundServer(
	logger log.Logger,
	pc *biz.PaidContent,
) *BackgroundServer {
	cs := &BackgroundServer{
		log: log.NewHelper(log.With(logger, "module", "server/background")),
		pc:  pc,
	}

	return cs
}

func (cs *BackgroundServer) Start(ctx context.Context) error {
	cs.log.Info("background server started")

	return nil
}

func (cs *BackgroundServer) Stop(ctx context.Context) error {
	cs.log.Info("background server stopped")

	return nil
}
