package modules

import (
	"os"
	"strings"
	"time"

	"github.com/AlphaCodinggroup/alpha_auth-api/pkg/metrics"
	"github.com/gin-contrib/cors"
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

// CORSMiddleware configures CORS headers for all routes.
// Allowed origins are read from env var ALLOWED_ORIGINS (comma separated).
// If ALLOWED_ORIGINS is "*", credentials are disabled to comply with CORS spec.
func CORSMiddleware() SharedMiddleware {
	raw := os.Getenv("ALLOWED_ORIGINS")
	var origins []string
	if raw == "" {
		// Default: allow all for local/dev unless overridden
		origins = []string{"*"}
	} else {
		origins = strings.Split(raw, ",")
	}

	allowCredentials := true
	for _, o := range origins {
		if strings.TrimSpace(o) == "*" {
			allowCredentials = false
			break
		}
	}

	cfg := cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: allowCredentials,
		MaxAge:           12 * time.Hour,
	}

	return SharedMiddleware{Middleware: cors.New(cfg)}
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
