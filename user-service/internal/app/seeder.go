package app

import (
	"context"

	"github.com/suprt/trading/user-service/internal/application"
	"github.com/suprt/trading/user-service/internal/core/domain"
	userconfig "github.com/suprt/trading/user-service/internal/infrastructure/config"
	"go.uber.org/fx"
)

type seederParams struct {
	fx.In

	LC   fx.Lifecycle
	Cfg  userconfig.Config
	Seed *application.Seed
}

func RunSeeder(p seederParams) {
	if !p.Cfg.Seed.Enabled {
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
					return err
				}

			}
			return nil
		},
	})
}
