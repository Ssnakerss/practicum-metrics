package server

import (
	"context"
	"net/http"

	"github.com/Ssnakerss/practicum-metrics/internal/app"
	"github.com/Ssnakerss/practicum-metrics/internal/dtadapter"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	c *app.ServerConfig
	l *zap.SugaredLogger
	A *dtadapter.Adapter
	s *http.Server
}

func New(ctx context.Context, l *zap.SugaredLogger) (*Server, error) {
	//Initialize and read server config.
	c := &app.ServerConfig{}
	if err := c.Read(); err != nil {
		return nil, err
	}

	//Initialize adapter with server config.
	a, err := InitAdapter(ctx, c)
	if err != nil {
		return nil, err
	}

	//Initialize router.
	r := NewRouter(a, c)

	s := &http.Server{
		Addr:    c.EndPointAddress,
		Handler: r,
	}

	return &Server{c, l, a, s}, nil
}

func (s *Server) Run(ctx context.Context) {
	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		s.l.Infow("startup", "config", s.c)
		s.l.Infow("starting server at ", "address", s.c.EndPointAddress)
		return s.s.ListenAndServe() //Запускаем сервер
	})
	g.Go(func() error {
		<-gCtx.Done() //Ожидаем завершения контекста
		s.l.Info("shutting server down")
		return s.s.Shutdown(context.Background()) //Завершаем сервер
	})

	if err := g.Wait(); err != nil {
		s.l.Info("performing pre-shutdown tasks")
		s.A.DoSync()
		s.A.Ds.Close()
		s.l.Warnw("server stopped", "error", err)
	}
}
