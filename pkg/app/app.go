package app

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"github.com/1024casts/snake/pkg/conf"
	"github.com/1024casts/snake/pkg/log"
)

type App struct {
	c      *conf.Config
	opts   options
	ctx    context.Context
	cancel func()
	log    log.Logger
}

func New(c *conf.Config, opts ...Option) *App {
	options := options{
		ctx:    context.Background(),
		logger: log.GetLogger(),
		sigs:   []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
	}
	if id, err := uuid.NewUUID(); err == nil {
		options.id = id.String()
	}
	for _, o := range opts {
		o(&options)
	}

	ctx, cancel := context.WithCancel(options.ctx)
	return &App{
		c:      c,
		opts:   options,
		ctx:    ctx,
		log:    log.GetLogger(),
		cancel: cancel,
	}
}

// Run start app
func (a *App) Run() error {
	a.log.Infof("app_id: %s, app_name: %s, version: %s",
		a.opts.id,
		a.opts.name,
		a.opts.version,
	)
	g, ctx := errgroup.WithContext(a.ctx)

	// start server
	for _, srv := range a.opts.servers {
		srv := srv
		g.Go(func() error {
			<-ctx.Done()
			return srv.Stop()
		})
		g.Go(func() error {
			return srv.Start()
		})
	}

	// watch signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, a.opts.sigs...)
	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case s := <-quit:
				a.log.Infof("receive a quit signal: %s", s.String())
				return a.Stop()
			}
		}
	})
	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}

// Stop stops the application gracefully.
func (a *App) Stop() error {
	if a.cancel != nil {
		a.cancel()
	}
	return nil
}
