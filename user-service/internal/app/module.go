package app

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	userv1 "github.com/suprt/trading/proto/user/v1"
	"github.com/suprt/trading/shared/bootstrap"
	grpchand "github.com/suprt/trading/user-service/internal/adapters/inbound/grpc"
	"github.com/suprt/trading/user-service/internal/adapters/inbound/interceptors"
	"github.com/suprt/trading/user-service/internal/adapters/outbound/crypto"
	"github.com/suprt/trading/user-service/internal/adapters/outbound/idgen"
	"github.com/suprt/trading/user-service/internal/adapters/outbound/inmemory"
	"github.com/suprt/trading/user-service/internal/adapters/outbound/jwt"
	"github.com/suprt/trading/user-service/internal/application"
	"github.com/suprt/trading/user-service/internal/core/ports/inbound"
	"github.com/suprt/trading/user-service/internal/core/ports/outbound"
	userconfig "github.com/suprt/trading/user-service/internal/infrastructure/config"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Module(cfg userconfig.Config) fx.Option {
	return fx.Module("user",
		fx.Provide(func() userconfig.Config {
			return cfg
		}),

		fx.Provide(
			func(cfg userconfig.Config) jwt.Config {
				return jwt.Config{
					Secret: []byte(cfg.Secret),
					TTL:    cfg.JWT.AccessTTL,
				}

			},

			func(cfg userconfig.Config) application.LoginConfig {
				return application.LoginConfig{
					RefreshTTL: cfg.JWT.RefreshTTL,
					AccessTTL:  cfg.JWT.AccessTTL,
				}
			},

			func(cfg userconfig.Config) application.RefreshTokenConfig {
				return application.RefreshTokenConfig{
					RefreshTTL: cfg.JWT.RefreshTTL,
					AccessTTL:  cfg.JWT.AccessTTL,
				}
			},
		),

		fx.Provide(
			fx.Annotate(
				interceptors.UnaryValidation,
				fx.ResultTags(`group:"grpc_server_interceptor"`),
			),
			fx.Annotate(
				func(cfg userconfig.Config, tokens outbound.TokenIssuer) grpc.UnaryServerInterceptor {
					return interceptors.UnaryAuth(tokens, cfg.Auth.Whitelist)
				},

				fx.ResultTags(`group:"grpc_server_interceptor"`),
			),
		),

		fx.Provide(
			fx.Annotate(inmemory.NewUserRepo, fx.As(new(outbound.UserRepository))),
			fx.Annotate(inmemory.NewCredentialsRepo, fx.As(new(outbound.CredentialsRepository))),
			fx.Annotate(inmemory.NewSessionRepo, fx.As(new(outbound.SessionRepository))),
			fx.Annotate(crypto.NewBcryptHasher, fx.As(new(outbound.PasswordHasher))),
			fx.Annotate(idgen.NewUUIDGen, fx.As(new(outbound.IDGenerator))),
			fx.Annotate(jwt.NewJWTIssuer, fx.As(new(outbound.TokenIssuer))),
		),

		//---- use cases ----
		fx.Provide(
			application.NewRegisterUseCase,
			application.NewLoginUseCase,
			application.NewRefreshTokenUseCase,
			application.NewLogoutUseCase,
			application.NewGetSelfProfileUseCase,
			application.NewGetUserByIDUseCase,
			application.NewAddRoleUseCase,
			application.NewRemoveRoleUseCase,
			application.NewService,
		),

		fx.Provide(
			fx.Annotate(
				func(s *application.Service) inbound.UserService { return s },
			),
		),

		fx.Provide(grpchand.New),
		fx.Invoke(registerGRPCHandlers),

		fx.Provide(application.NewSeed),
		fx.Invoke(RunSeeder),
	)
}
func registerGRPCHandlers(srv *grpc.Server, h *grpchand.Handler) {
	userv1.RegisterUserServiceServer(srv, h)
	reflection.Register(srv)
}

func Start() {
	_ = godotenv.Load()
	yamlPath := os.Getenv("CONFIG_YAML_PATH")
	if yamlPath == "" {
		yamlPath = "config.yaml"
	}
	cfg, err := userconfig.Load(yamlPath)
	if err != nil {
		log.Fatal(err)
	}

	bootstrap.NewApp(bootstrap.Options{
		Addr:        cfg.GRPC.Addr,
		Modules:     []fx.Option{Module(*cfg)},
		StopTimeout: cfg.GRPC.StopTimeout,
	}).Run()
}
