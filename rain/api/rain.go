package main

import (
	"flag"
	"fmt"
	"go-red-envelope-rain/rain/api/internal/config"
	"go-red-envelope-rain/rain/api/internal/handler"
	"go-red-envelope-rain/rain/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "/etc/rain-api.yaml", "the config file")
var sysMoney = flag.String("balance", "10000000000", "balance, unit=penny")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)

	err := ctx.RedisClient.Set("rain:balance", *sysMoney, 0).Err()
	if err != nil {
		panic(err)
	}

	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
