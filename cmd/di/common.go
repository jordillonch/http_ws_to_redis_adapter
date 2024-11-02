package di

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/jordillonch/http_ws_to_redis_adapter/configs"
	"github.com/jordillonch/http_ws_to_redis_adapter/internal/app/http-ws-server/domain/topic-message-processor/commands"
	"github.com/jordillonch/http_ws_to_redis_adapter/internal/app/http-ws-server/domain/topic-message-processor/events"
	"go.uber.org/zap"
	"time"
)

type CommonDi struct {
	Services             *CommonServices
	Env                  config.Config
	MessageProcessors    *MessageProcessors
	HttpWsServerServices *HttpWsServerModule
}

type CommonServices struct {
	Logger      *zap.Logger
	RedisClient *redis.Client
}

type MessageProcessors struct {
	Commands *commands.Processor
	Events   *events.Processor
}

func Init() *CommonDi {
	return setUp()
}

func InitWithEnvFile(envFiles ...string) *CommonDi {
	err := godotenv.Overload(envFiles...)
	if err != nil {
		panic(err)
	}

	return setUp()
}

func setUp() *CommonDi {
	cnf := config.BuildConfig()
	l := buildLogger(cnf)
	redisClient := buildRedisClient(cnf)

	services := &CommonServices{
		Logger:      l,
		RedisClient: redisClient,
	}

	processors := &MessageProcessors{
		Commands: commands.NewProcessor(l),
		Events:   events.NewProcessor(l),
	}

	httpWsServerDi := &CommonDi{
		Services:             services,
		Env:                  cnf,
		HttpWsServerServices: NewHttpWsServerModule(services, cnf, processors),
	}

	return httpWsServerDi
}

func buildLogger(config config.Config) *zap.Logger {
	var lCnf zap.Config
	if config.IsDevelopment {
		lCnf = zap.NewDevelopmentConfig()
	} else {
		lCnf = zap.NewProductionConfig()
	}

	switch loggerLevel := config.LoggerLevel; loggerLevel {
	case "debug":
		lCnf.Level.SetLevel(zap.DebugLevel)
	case "info":
		lCnf.Level.SetLevel(zap.InfoLevel)
	case "warn":
		lCnf.Level.SetLevel(zap.WarnLevel)
	case "error":
		lCnf.Level.SetLevel(zap.ErrorLevel)
	default:
		lCnf.Level.SetLevel(zap.WarnLevel)
	}

	logger, err := lCnf.Build()
	if err != nil {
		panic(err)
	}

	return logger
}

func buildRedisClient(cnf config.Config) *redis.Client {
	redisOptions := &redis.Options{
		Addr:        fmt.Sprintf("%s:%d", cnf.RedisHost, cnf.RedisPort),
		IdleTimeout: time.Duration(cnf.RedisIdleTimeout),
		PoolSize:    cnf.RedisPoolSize,
		PoolTimeout: time.Duration(cnf.RedisPoolTimeout),
	}
	return redis.NewClient(redisOptions)
}
