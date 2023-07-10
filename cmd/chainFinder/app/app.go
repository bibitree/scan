package app

import (
	"context"
	"ethgo/chainFinder"
	"ethgo/model"
	"ethgo/util/ginx"
	"fmt"
	"net/http"
	"time"

	// _ "ethgo/cmd/tyche/docs"

	"github.com/gin-gonic/gin"
)

type App struct {
	base *chainFinder.ChainFinder
	conf *Config
}

func New(conf *Config) App {
	return App{conf: conf}
}

func (app *App) Init(ctx context.Context) error {
	err := model.Init(app.conf.Redis)
	if err != nil {
		return err
	}
	err = model.InitMysql(app.conf.Redis)
	if err != nil {
		return err
	}

	app.base = chainFinder.New(app.conf.ChainFinder)

	return nil
}
func (app *App) Run(ctx context.Context) error {
	engine := gin.New()
	engine.Use(gin.Recovery())
	// engine.Use(gin.Logger())
	engine.RedirectTrailingSlash = false
	engine.Use(ginx.Cors())
	engine.Static("/img", "./img")
	app.Router(engine)

	var c = app.conf.ChainFinder
	fmt.Println(c)
	srv := new(http.Server)

	srv.Addr = c.Listen
	srv.ReadTimeout = time.Duration(c.ReadTimeout) * time.Second
	srv.WriteTimeout = time.Duration(c.WriteTimeout) * time.Second
	srv.MaxHeaderBytes = c.MaxHeaderBytes
	srv.Handler = engine

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
	}()

	go func() {
		var err error
		if c.EnableTLS {
			err = srv.ListenAndServeTLS(c.CertFile, c.KeyFile)
		} else {
			err = srv.ListenAndServe()
		}

		if err != http.ErrServerClosed {
			panic(err)
		}
	}()

	defer model.Dispose()
	go app.setChainDataLoop()
	return app.base.Run(ctx)
}
