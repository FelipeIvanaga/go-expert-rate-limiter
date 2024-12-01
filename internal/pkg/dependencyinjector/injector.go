package dependencyinjector

import (
	"time"

	"github.com/felipeivanaga/go-expert-rate-limiter/config"
	"github.com/felipeivanaga/go-expert-rate-limiter/internal/infra/database"
	"github.com/felipeivanaga/go-expert-rate-limiter/internal/infra/web"
	"github.com/felipeivanaga/go-expert-rate-limiter/internal/infra/web/handlers"
	"github.com/felipeivanaga/go-expert-rate-limiter/internal/infra/web/middlewares"
	"github.com/felipeivanaga/go-expert-rate-limiter/internal/pkg/ratelimiter"
	ratelimiter_strategies "github.com/felipeivanaga/go-expert-rate-limiter/internal/pkg/ratelimiter/strategies"
	"github.com/felipeivanaga/go-expert-rate-limiter/internal/pkg/responsehandler"
)

type DependencyInjectorInterface interface {
	Inject() (*Dependencies, error)
}

type DependencyInjector struct {
	Config *config.Conf
}

type Dependencies struct {
	ResponseHandler       responsehandler.WebResponseHandlerInterface
	HelloWebHandler       handlers.HelloWebHandlerInterface
	RateLimiterMiddleware middlewares.RateLimiterMiddlewareInterface
	WebServer             web.WebServerInterface
	RedisDatabase         database.RedisDatabaseInterface
	RateLimiter           ratelimiter.RateLimiterInterface
	RedisLimiterStrategy  ratelimiter_strategies.LimiterStrategyInterface
}

func NewDependencyInjector(c *config.Conf) *DependencyInjector {
	return &DependencyInjector{
		Config: c,
	}
}

func (di *DependencyInjector) Inject() (*Dependencies, error) {
	responseHandler := responsehandler.NewWebResponseHandler()

	redisDB, err := database.NewRedisDatabase(*di.Config)
	if err != nil {
		return nil, err
	}

	redisLimiterStrategy := ratelimiter_strategies.NewRedisLimiterStrategy(
		redisDB.Client,
		time.Now,
	)

	limiter := ratelimiter.NewRateLimiter(
		redisLimiterStrategy,
		di.Config.IPMaxRequests,
		di.Config.TokenMaxRequests,
		di.Config.TimeWindowMilliseconds,
	)

	helloWebHandler := handlers.NewHelloWebHandler(responseHandler)
	rateLimiterMiddleware := middlewares.NewRateLimiterMiddleware(responseHandler, limiter)

	webRouter := web.NewWebRouter(helloWebHandler, rateLimiterMiddleware)
	webServer := web.NewWebServer(
		di.Config.ServerPort,
		webRouter.Build(),
		webRouter.BuildMiddlewares(),
	)

	return &Dependencies{
		ResponseHandler:       responseHandler,
		HelloWebHandler:       helloWebHandler,
		RateLimiterMiddleware: rateLimiterMiddleware,
		WebServer:             webServer,
		RedisDatabase:         redisDB,
		RateLimiter:           limiter,
		RedisLimiterStrategy:  redisLimiterStrategy,
	}, nil
}
