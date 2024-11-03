package router

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/TakeAway-Inc/platform/logger"
	"github.com/TakeAway-Inc/platform/metrics"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type Router struct {
	r *gin.Engine
}

func New(addr string, log *logger.Logger) *Router {
	if env := os.Getenv("APP_ENV"); env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := &Router{r: gin.New()}

	r.r.Use(otelgin.Middleware(addr))
	r.r.Use(logMiddleware(log))
	r.r.Use(metricMiddleware())

	r.r.GET("/status", r.status)
	r.r.GET("/metrics", r.metrics)
	r.r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return r
}

func (r *Router) Router() *gin.Engine {
	return r.r
}

// @Summary Metrics
// @Description Prometheus metrics
// @Tags Status
// @Produce application/json
// @Success 200
// @Router /metrics [get]
func (r *Router) metrics(ctx *gin.Context) {
	promhttp.Handler().ServeHTTP(ctx.Writer, ctx.Request)
}

// @Summary Status Check
// @Description Checking status of backend
// @Tags Status
// @Produce application/json
// @Success 200
// @Router /status [get]
func (r *Router) status(ctx *gin.Context) {
	ctx.String(http.StatusOK, "ok\n")
}

func logMiddleware(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Info("got http request", slog.String("method", c.Request.Method), slog.String("path", c.Request.URL.Path))

		c.Next()
	}
}

func metricMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		metrics.HTTPRequestsCount.With(map[string]string{
			"uri":    c.Request.URL.Path,
			"method": c.Request.Method,
		}).Inc()

		c.Next()
	}
}
