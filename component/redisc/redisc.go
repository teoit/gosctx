package redisc

import (
	"context"
	"encoding/json"
	"time"

	"github.com/teoit/gosctx"
	"github.com/teoit/gosctx/component/errs"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/teoit/gosctx/configs"
)

type RedisComponent interface {
	GetClient() *redis.Client
	SetCacheRedis(ctx context.Context, data *string, key string, sessionLifetime int) error
	GetCacheRedis(ctx context.Context, key string) (*string, error)
	DeleteCacheRedis(ctx context.Context, key string) error
}

type redisc struct {
	id        string
	client    *redis.Client
	logger    gosctx.Logger
	redisUri  string
	maxActive int
	maxIde    int
}

func NewRedisc(id string) *redisc {
	return &redisc{id: id}
}

func (r *redisc) ID() string {
	return r.id
}

func (r *redisc) InitFlags() {
	r.redisUri = configs.RedisUri
	r.maxActive = configs.MaxActive
	r.maxIde = configs.MaxIde
}

func (r *redisc) Activate(sc gosctx.ServiceContext) error {
	r.logger = gosctx.GlobalLogger().GetLogger(r.id)
	r.logger.Info("Connecting to cache Redis at ", r.redisUri, "...")

	opt, err := redis.ParseURL(r.redisUri)

	if err != nil {
		r.logger.Error("Cannot parse cache Redis ", err.Error())
		return err
	}

	opt.PoolSize = r.maxActive
	opt.MinIdleConns = r.maxIde

	client := redis.NewClient(opt)

	// Ping to test Redis connection
	if err = client.Ping(context.Background()).Err(); err != nil {
		r.logger.Error("Cannot connect cache Redis. ", err.Error())
		return err
	}

	// Enable tracing instrumentation.
	if err = redisotel.InstrumentTracing(client); err != nil {
		panic(err)
	}

	// Enable metrics instrumentation.
	if err = redisotel.InstrumentMetrics(client); err != nil {
		panic(err)
	}

	// Connect successfully, assign client to goRedisDB
	r.client = client
	return nil
}

func (r *redisc) Stop() error {
	if err := r.client.Close(); err != nil {
		return err
	}

	return nil
}

func (r *redisc) GetClient() *redis.Client {
	return r.client
}

func (r *redisc) SetCacheRedis(ctx context.Context, data *string, key string, sessionLifetime int) error {
	if key == "" {
		return errs.ErrKeyCacheRedisNotEmpty
	}
	out, err := json.Marshal(*data)
	if err != nil {
		return err
	}
	exp := time.Duration(sessionLifetime)
	if err = r.client.Set(ctx, key, out, exp).Err(); err != nil {
		return err
	}
	return nil
}

func (r *redisc) GetCacheRedis(ctx context.Context, key string) (*string, error) {
	if key == "" {
		return nil, errs.ErrKeyCacheRedisNotEmpty
	}
	result, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, errs.ErrDataNotFound
	} else if err != nil {
		return nil, err
	}

	var data string
	err = json.Unmarshal([]byte(result), &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// delete redis
func (r *redisc) DeleteCacheRedis(ctx context.Context, key string) error {
	if key == "" {
		return errs.ErrKeyCacheRedisNotEmpty
	}
	return r.client.Del(ctx, key).Err()
}
