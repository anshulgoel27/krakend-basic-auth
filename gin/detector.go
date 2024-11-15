package gin

import (
	"errors"
	"net/http"

	auth "github.com/anshulgoel27/krakend-basic-auth"
	"github.com/gin-gonic/gin"
	"github.com/luraproject/lura/v2/config"
	"github.com/luraproject/lura/v2/logging"
	"github.com/luraproject/lura/v2/proxy"
	krakendgin "github.com/luraproject/lura/v2/router/gin"
)

const logPrefix = "[SERVICE: Gin][basic-auth]"

// Register checks the configuration and, if required, registers a bot detector middleware at the gin engine
func Register(cfg config.ServiceConfig, l logging.Logger, engine *gin.Engine) {
	detectorCfg, err := auth.ParseConfig(cfg.ExtraConfig)
	if err == auth.ErrNoConfig {
		return
	}
	if err != nil {
		l.Warning(logPrefix, err.Error())
		return
	}
	d := auth.New(detectorCfg)
	engine.Use(middleware(d, l))
}

// New checks the configuration and, if required, wraps the handler factory with a bot detector middleware
func New(hf krakendgin.HandlerFactory, l logging.Logger) krakendgin.HandlerFactory {
	return func(cfg *config.EndpointConfig, p proxy.Proxy) gin.HandlerFunc {
		next := hf(cfg, p)
		logPrefix := "[ENDPOINT: " + cfg.Endpoint + "][basic-auth]"

		detectorCfg, err := auth.ParseConfig(cfg.ExtraConfig)
		if err == auth.ErrNoConfig {
			return next
		}
		if err != nil {
			l.Warning(logPrefix, err.Error())
			return next
		}

		d := auth.New(detectorCfg)
		return handler(d, next, l)
	}
}

func middleware(f auth.AuthFunc, l logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !f(c.Request) {
			l.Error(logPrefix, errBasicAuthRejected)
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}

func handler(f auth.AuthFunc, next gin.HandlerFunc, l logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !f(c.Request) {
			l.Error(logPrefix, errBasicAuthRejected)
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		next(c)
	}
}

var errBasicAuthRejected = errors.New("basic auth rejected")
