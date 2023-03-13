package svc

import (
	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"go-red-envelope-rain/rain/api/internal/config"
	"go-red-envelope-rain/rain/model"
	"go-red-envelope-rain/user/rpc/userclient"
)

type ServiceContext struct {
	Config      config.Config
	RainModel   model.RainModel
	UserRpc     userclient.User
	RedisClient *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:      c,
		RainModel:   model.NewRainModel(conn, c.CacheRedis),
		UserRpc:     userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		RedisClient: redis.NewClient(&redis.Options{Addr: c.CacheRedis[0].Host}),
	}
}
