package app

import (
	"context"

	"github.com/suprt/trading/user-service/internal/application"
	"github.com/suprt/trading/user-service/internal/core/domain"
	userconfig "github.com/suprt/trading/user-service/internal/infrastructure/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type seederParams struct {
	fx.In

	LC     fx.Lifecycle
	Cfg    userconfig.Config
	Seed   *application.Seed
	Logger *zap.Logger
}

func RunSeeder(p seederParams) {
	if !p.Cfg.Seed.Enabled {
		p.Logger.Info("seeder disabled")
		return
	}

	p.LC.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			for _, su := range p.Cfg.Seed.Users {
				roles := make([]domain.Role, len(su.Roles))
				for i, r := range su.Roles {
					roles[i] = domain.Role(r)
				}

				if err := p.Seed.SeedUser(ctx, application.SeedUserInput{
					Email:    su.Email,
					Password: su.Password,
					Roles:    roles,
				}); err != nil {
					p.Logger.Error("seed user failed", zap.String("email", su.Email), zap.Error(err))
					return err
				}

			}
			p.Logger.Info("seed completed")
			return nil
		},
	})
}
