package app

import (
	"context"
	"ethgo/eth"
	"ethgo/model"
	"ethgo/tyche"
	"net/http"
	"time"

	_ "ethgo/cmd/tyche/docs"

	"github.com/gin-gonic/gin"
)

type App struct {
	backend eth.Backend
	base    *tyche.Tyche
	conf    *Config
}

func New(conf *Config) App {
	return App{conf: conf}
}

func (app *App) Init(ctx context.Context) error {
	err := model.Init(app.conf.Redis)
	if err != nil {
		return err
	}

	if app.backend, err = eth.New(app.conf.Backend); err != nil {
		return err
	}

	app.base = tyche.New(app.backend, app.conf.Tyche)
	return app.base.Init(ctx)
}

func (app *App) Run(ctx context.Context) error {
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(gin.Logger())
	engine.RedirectTrailingSlash = false
	app.Router(engine)

	var c = app.conf.Tyche
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

		if app.backend != nil {
			app.backend.Close()
		}
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
	return app.base.Run(ctx)
}
