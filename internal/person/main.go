package person

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alukart32/effective-mobile-test-task/internal/person/adapters"
	"github.com/alukart32/effective-mobile-test-task/internal/person/ports"
	"github.com/alukart32/effective-mobile-test-task/internal/person/service/persondata"
	"github.com/alukart32/effective-mobile-test-task/internal/person/storage/persons"
	"github.com/alukart32/effective-mobile-test-task/internal/pkg/ginx"
	"github.com/alukart32/effective-mobile-test-task/internal/pkg/postgres"
	"github.com/alukart32/effective-mobile-test-task/internal/pkg/server"
	"github.com/alukart32/effective-mobile-test-task/internal/pkg/zerologx"
	"github.com/caarlos0/env/v8"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type config struct {
	API struct {
		AgifyService     string `env:"AGIFY_SERVICE_API,notEmpty"`
		GenderizeService string `env:"GENDERIZE_SERVICE_API,notEmpty"`
		Nationalize      string `env:"NATIONALIZE_SERVICE_API,notEmpty"`
	}
	GraphQL struct {
		Path string `env:"GRAPHQL_PATH,notEmpty"`
	}
	Kafka struct {
		ReadTopic string   `env:"KAFKA_READ_TOPIC,notEmpty"`
		ErrTopic  string   `env:"KAFKA_ERROR_TOPIC,notEmpty"`
		Brokers   []string `env:"KAFKA_BROKERS,notEmpty"`
		ReadLimit int      `env:"KAFKA_READ_LIMIT" envDefault:"1"`
	}
	Postgres struct {
		URL string `env:"POSTGRES_URL,notEmpty"`
	}
	Redis struct {
		URL         string        `env:"REDIS_ADDRESS,notEmpty" envDefault:"127.0.0.1:6379"`
		DialTimeout time.Duration `env:"REDIS_DEAL_TIMEOUT" envDefault:"100ms"`
		ReadTimeout time.Duration `env:"REDIS_READ_TIMEOUT" envDefault:"100ms"`
	}
}

func Run() {
	appCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := zerologx.Get()

	var cfg config
	err := env.ParseWithOptions(&cfg, env.Options{RequiredIfNoDef: true})
	if err != nil {
		logger.Fatal().Err(err).Msg("parse env params")
	}

	redisOpt, err := redis.ParseURL(cfg.Redis.URL)
	if err != nil {
		logger.Fatal().Err(err).Msg("prepare redis client")
	}
	cache := redis.NewClient(redisOpt)
	if _, err := cache.Ping(appCtx).Result(); err != nil {
		logger.Fatal().Err(err).Msg("prepare redis client")
	}
	defer func() {
		if err := cache.ShutdownSave(appCtx).Err(); err != nil {
			logger.Err(err).Msg("redis client shutdown")
		}
	}()

	postgresPool, err := postgres.Get(cfg.Postgres.URL)
	if err != nil {
		logger.Fatal().Err(err).Msg("prepare postgres pool")
	}
	defer func() {
		postgresPool.Close()
	}()

	repo, _ := persons.CachedStorage(postgresPool, cache)

	personMetaDataProvider, err := adapters.PersonMetaData(
		cfg.API.AgifyService,
		cfg.API.GenderizeService,
		cfg.API.Nationalize,
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("prepare person metadata provider adapter")
	}

	personManager, err := persondata.Manager(repo, personMetaDataProvider)
	if err != nil {
		logger.Fatal().Err(err).Msg("prepare person manager")
	}

	// Prepare API
	_, err = ports.KafkaFIO(
		appCtx,
		cfg.Kafka.ReadTopic,
		cfg.Kafka.ErrTopic,
		cfg.Kafka.Brokers,
		cfg.Kafka.ReadLimit,
		personManager,
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("prepare kafka FIO message handler")
	}

	err = ports.HttpRoutes(gin.New(), personManager)
	if err != nil {
		logger.Fatal().Err(err).Msg("prepare person REST routes")
	}

	ginRouter, err := ginx.Get()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to prepare: gin")
	}
	ports.HttpRoutes(ginRouter, personManager)
	if err != nil {
		logger.Fatal().Err(err).Msg("prepare HTTP routes")
	}
	ports.Graph(ginRouter, cfg.GraphQL.Path, personManager)
	if err != nil {
		logger.Fatal().Err(err).Msg("prepare GraphQL")
	}

	// Run servers
	httpServer, err := server.Http(server.HttpConfig{}, ginRouter)
	if err != nil {
		logger.Fatal().Err(err).Msg("prepare HTTP server")
	}
	defer func() {
		if err := httpServer.Shutdown(); err != nil {
			logger.Err(err).Msg("shutdown HTTP server")
		}
	}()

	// Waiting signals.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case s := <-interrupt:
		logger.Info().Msg(s.String())
	case err = <-httpServer.Notify():
		logger.Fatal().Err(err).Send()
	}
}
