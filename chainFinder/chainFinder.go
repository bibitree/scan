package chainFinder

import (
	"context"
	"sync"
	"time"
)

type ChainFinder struct {
	conf                  *Config
	suggestNonceKeepalive time.Duration
}

func New(conf *Config) *ChainFinder {
	return &ChainFinder{
		conf:                  conf,
		suggestNonceKeepalive: 0,
	}
}

func (t *ChainFinder) Run(ctx context.Context) error {
	return t.run(ctx)
}

func (t *ChainFinder) run(ctx context.Context) error {
	log.Info("开始侦听")
	defer log.Info("结束侦听")

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		t.WatchTransactionStorage(ctx)
	}()

	wg.Wait()
	return nil
}
