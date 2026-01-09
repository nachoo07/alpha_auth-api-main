package modules

import (
	"os"

	"github.com/AlphaCodinggroup/alpha_auth-api/cmd/api/app/handlers"
	"github.com/AlphaCodinggroup/alpha_auth-api/cmd/api/app/middleware"
	"github.com/AlphaCodinggroup/alpha_auth-api/cmd/api/app/router"
	"github.com/AlphaCodinggroup/alpha_auth-api/internal/auth"
	"github.com/AlphaCodinggroup/alpha_auth-api/internal/auth/session"
	"github.com/AlphaCodinggroup/alpha_auth-api/internal/auth/user"
	"github.com/AlphaCodinggroup/alpha_auth-api/internal/platform/environment"
	"github.com/AlphaCodinggroup/alpha_auth-api/internal/platform/httpserver"
	"github.com/AlphaCodinggroup/alpha_auth-api/internal/platform/repository/pg"
	authcontext "github.com/AlphaCodinggroup/alpha_auth-api/pkg/context"
	"github.com/AlphaCodinggroup/alpha_auth-api/pkg/validator"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

func NewApp() *fx.App {
	env := environment.GetFromString(os.Getenv("GO_ENVIRONMENT"))

	if env == environment.Local {
		// Cargar .env solo si existe, sin fallar si no está
		godotenv.Load(".env")
	}

	authcontext.NewLogger()

	options := []fx.Option{
		fx.Provide(
			fx.Annotate(
				session.NewService,
				fx.As(new(auth.TokenGenerator)),
				fx.As(new(middleware.TokenValidator)),
			),
			fx.Annotate(
				middleware.NewAuthMiddleware,
				fx.As(new(router.AuthMiddleware)),
			),
			fx.Annotate(
				user.NewRepository,
				fx.As(new(middleware.UserQuery)),
				fx.As(new(middleware.UserCommand)),
				fx.As(new(auth.UserFinder)),
				fx.As(new(auth.UserTokenManager)),
				fx.As(new(auth.UserAccountCommand)),
				fx.As(new(auth.LoginCommand)),
			),
			pg.GetDB,
			fx.Annotate(
				auth.NewLoginUseCase,
				fx.As(new(handlers.LoginUseCase)),
			),
			fx.Annotate(
				auth.NewRegisterUseCase,
				fx.As(new(handlers.RegisterUseCase)),
			),
			fx.Annotate(
				handlers.NewLoginHandler,
				fx.As(new(router.LoginHandler)),
			),
			fx.Annotate(
				handlers.NewSignupHandler,
				fx.As(new(router.SignupHandler)),
			),
			NewRouter,
			httpserver.NewHTTPGinServer,

			MetricMiddleware,
		),
	}

	return fx.New(
		fx.Options(options...),
		fx.Invoke(RegisterMiddlewares),
		fx.Invoke(validator.RegisterValidation),
		fx.Invoke(router.RegisterRouter),
		fx.Invoke(httpserver.StartServer),
	)
}

func NewRouter() *gin.Engine {
	return gin.Default()
}
