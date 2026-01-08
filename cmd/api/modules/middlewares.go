package modules

import (
	"github.com/AlphaCodinggroup/alpha_auth-api/pkg/metrics"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type (
	SharedMiddleware struct {
		fx.Out

		Middleware gin.HandlerFunc `group:"shared"`
	}

	EnabledMiddlewares struct {
		fx.In

		// WebMiddlewares       []web.Middleware `group:"web"`
		// BigQueueMiddlewares  []web.Middleware `group:"bigq"`
		// WorkQueueMiddlewares []web.Middleware `group:"workq"`
		SharedMiddlewares []gin.HandlerFunc `group:"shared"`
	}
)

// func LogRequestMiddleware(app *fury.Application, cfg config.Configuration) SharedMiddleware {
// 	return SharedMiddleware{
// 		Middleware: web.LogRequest(app.Logger, web.LogRequestConfig{
// 			IncludeRequest:  cfg.LogRequests,
// 			IncludeResponse: cfg.LogRequests,
// 		}),
// 	}
// }

// func UserMiddleware() WebMiddleware {
// 	return WebMiddleware{
// 		Middleware: session.Middleware(),
// 	}
// }

func MetricMiddleware() SharedMiddleware {
	return SharedMiddleware{
		Middleware: metrics.Middleware(),
	}
}

func RegisterMiddlewares(router *gin.Engine, params EnabledMiddlewares) {
	for _, mw := range params.SharedMiddlewares {
		router.Use(mw)
	}

	router.Use(func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			c.AbortWithStatusJSON(500, gin.H{"error": "Internal Server Error"})
		}
	})
}
