package main

import (
	"context"
	"ethgo/cmd/chainFinder/app"
	"ethgo/util/logx"
	"os"
	"os/signal"
	"syscall"
)

func run(c *app.Config) {
	app := app.New(c)
	ctx, cancel := context.WithCancel(context.Background())
	if err := app.Init(ctx); err != nil {
		panic(err)
	}

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
		<-ch
		defer cancel()
	}()

	if err := app.Run(ctx); err != nil {
		panic(err)
	}
}

// @title chainFinder 服务
// @version 1.0
// @description chainFinder 服务的目标是提供区块链浏览器的服务。
// @description
func main() {
	c, err := app.NewConfig("./chainFinder.toml")
	if err != nil {
		panic(err)
	}

	var log = logx.New(c.Logger)
	app.SetLogger(log)

	log.Info("启动")
	defer log.Info("退出")

	run(c)
}
