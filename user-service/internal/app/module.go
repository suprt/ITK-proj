package app

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	userv1 "github.com/suprt/trading/proto/user/v1"
	"github.com/suprt/trading/shared/bootstrap"
	"github.com/suprt/trading/shared/db"
	sharedinters "github.com/suprt/trading/shared/interceptors"
	"github.com/suprt/trading/shared/logger"
	shareredis "github.com/suprt/trading/shared/redis"
	grpchand "github.com/suprt/trading/user-service/internal/adapters/inbound/grpc"
	userinters "github.com/suprt/trading/user-service/internal/adapters/inbound/interceptors"
	"github.com/suprt/trading/user-service/internal/adapters/outbound/cached"
	"github.com/suprt/trading/user-service/internal/adapters/outbound/crypto"
	"github.com/suprt/trading/user-service/internal/adapters/outbound/idgen"
	"github.com/suprt/trading/user-service/internal/adapters/outbound/jwt"
	"github.com/suprt/trading/user-service/internal/adapters/outbound/postgres"
	"github.com/suprt/trading/user-service/internal/adapters/outbound/rediscache"
	"github.com/suprt/trading/user-service/internal/adapters/outbound/zaperror"
	"github.com/suprt/trading/user-service/internal/application"
	"github.com/suprt/trading/user-service/internal/core/ports/inbound"
	"github.com/suprt/trading/user-service/internal/core/ports/outbound"
	userconfig "github.com/suprt/trading/user-service/internal/infrastructure/config"
	"github.com/suprt/trading/user-service/internal/infrastructure/migrations"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Module(cfg userconfig.Config) fx.Option {
	return fx.Module("user",
		fx.Provide(func() userconfig.Config {
			return cfg
		}),
		// -- Postgres --
		fx.Provide(
			func(cfg userconfig.Config) db.Config {
				return db.Config{
					DSN:             cfg.Postgres.DSN,
					MaxConns:        cfg.Postgres.MaxConns,
					MinConns:        cfg.Postgres.MinConns,
					MaxConnLifetime: cfg.Postgres.MaxConnLifetime,
				}
			},
			db.NewPostgres,
		),

		// -- Redis --
		fx.Provide(
			func(cfg userconfig.Config) shareredis.Config {
				return shareredis.Config{
					Addr:     cfg.Redis.Addr,
					Username: cfg.Redis.Username,
					Password: cfg.Redis.Password,
					DB:       cfg.Redis.DB,
				}
			},
			shareredis.NewClient,
		),

		// -- Кэш + ErrorReporter--
		fx.Provide(
			func(cfg userconfig.Config) rediscache.UserCacheConfig {
				return rediscache.UserCacheConfig{
					UserTTL:         cfg.Cache.UserTTL,
					UserNotFoundTTL: cfg.Cache.UserNotFoundTTL,
				}
			},
		),
		fx.Provide(
			fx.Annotate(rediscache.NewUserCache, fx.As(new(outbound.UserCache))),
			fx.Annotate(zaperror.NewReporter, fx.As(new(outbound.ErrorReporter))),
		),

		// -- Миграции --
		fx.Invoke(func(lc fx.Lifecycle, cfg userconfig.Config) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					return migrations.Up(ctx, cfg.Postgres.DSN)
				},
			})
		}),

		// --Redis ping--
		fx.Invoke(func(lc fx.Lifecycle, client *redis.Client, logger *zap.Logger) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
					defer cancel()
					if err := client.Ping(pingCtx).Err(); err != nil {
						logger.Warn("redis ping failed, cache disabled", zap.Error(err))
					}
					return nil
				},
			})
		}),

		// -- Config --
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
			func(cfg userconfig.Config) crypto.Config {
				return crypto.Config{
					HashCost: cfg.Crypto.HashCost,
				}
			},
		),

		// -- Interceptors --
		fx.Provide(
			func(
				logger *zap.Logger,
				tokens outbound.TokenIssuer,
				cfg userconfig.Config,
			) grpc.UnaryServerInterceptor {
				return sharedinters.Chain(
					sharedinters.PanicRecovery(logger),
					sharedinters.XRequestID(),
					sharedinters.Logger(logger),
					userinters.UnaryAuth(tokens, cfg.Auth.Whitelist),
					userinters.UnaryValidation(),
				)
			},
		),

		//-- IdempotencyStore --
		fx.Provide(
			func(cfg userconfig.Config) application.RegisterConfig {
				return application.RegisterConfig{
					IdempotencyTTL: cfg.Idempotency.TTL,
				}
			},
		),

		//-- Outbound adapters --
		fx.Provide(
			postgres.NewUserRepo,
			fx.Annotate(
				func(repo *postgres.UserRepo,
					cache outbound.UserCache,
					reporter outbound.ErrorReporter,
				) outbound.UserRepository {
					return cached.NewUserRepo(repo, cache, reporter)
				},
			),

			//fx.Annotate(postgres.NewUserRepo, fx.As(new(outbound.UserRepository))),
			fx.Annotate(postgres.NewCredentialsRepo, fx.As(new(outbound.CredentialsRepository))),
			fx.Annotate(postgres.NewSessionRepo, fx.As(new(outbound.SessionRepository))),
			fx.Annotate(crypto.NewBcryptHasher, fx.As(new(outbound.PasswordHasher))),
			fx.Annotate(idgen.NewUUIDGen, fx.As(new(outbound.IDGenerator))),
			fx.Annotate(jwt.NewJWTIssuer, fx.As(new(outbound.TokenIssuer))),
			fx.Annotate(postgres.NewUoW, fx.As(new(outbound.UnitOfWork))),
			fx.Annotate(rediscache.NewIdempotencyStore, fx.As(new(outbound.IdempotencyStore))),
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
			application.NewGetPublicProfileUseCase,
			application.NewUpdateProfileUseCase,
			application.NewChangePasswordUseCase,
			application.NewAdminUpdateProfileUseCase,
			application.NewAdminResetPasswordUseCase,
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

	zapLogger, err := logger.New(logger.Config{
		Level:  cfg.Log.Level,
		Format: cfg.Log.Format,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer func(zapLogger *zap.Logger) {
		err := zapLogger.Sync()
		if err != nil {
			log.Fatal(err)
		}
	}(zapLogger)

	sugar := zapLogger.Sugar()
	sugar.Infow("service starting",
		"name", "user-service",
		"addr", cfg.GRPC.Addr,
		"log_level", cfg.Log.Level,
		"log_format", cfg.Log.Format,
	)
	bootstrap.NewApp(bootstrap.Options{
		Addr:        cfg.GRPC.Addr,
		Modules:     []fx.Option{Module(*cfg)},
		StopTimeout: cfg.GRPC.StopTimeout,
		Log:         zapLogger,
	}).Run()
	sugar.Info("service stopped")
}
