package server

import (
	"context"
	"net"
	"net/http"

	"github.com/Ssnakerss/practicum-metrics/internal/app"
	"github.com/Ssnakerss/practicum-metrics/internal/dtadapter"
	"github.com/Ssnakerss/practicum-metrics/internal/encrypt"
	"github.com/Ssnakerss/practicum-metrics/internal/mygrpc"
	"github.com/Ssnakerss/practicum-metrics/proto"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type Server struct {
	c *app.ServerConfig
	l *zap.SugaredLogger
	A *dtadapter.Adapter
	s *http.Server
	e *encrypt.Coder //Кодировка

	g *grpc.Server
}

func New(ctx context.Context, l *zap.SugaredLogger) (*Server, error) {
	//Initialize and read server config.
	c := app.MakeServerConfig()

	//Initialize adapter with server config.
	a, err := InitAdapter(ctx, c)
	if err != nil {
		return nil, err
	}

	//create encoder if cryptokey is set
	var e *encrypt.Coder
	if c.CryptoKey != "" {
		e = &encrypt.Coder{}
		err := e.LoadPrivateKey(c.CryptoKey)
		if err == nil {
			l.Info("Private key loaded", e)
		} else {
			l.Warnw("Can't load private key: ", "error", err)
			e = nil
		}
	}

	//Initialize router.
	r := NewRouter(a, c, e)

	s := &http.Server{
		Addr:    c.Address,
		Handler: r,
	}

	//creating grpc serveris address is set not empty
	ms := mygrpc.Server{
		Da: a,
		C:  c,
		E:  e,
		L:  l,
	}
	g := grpc.NewServer()
	proto.RegisterMetricsServer(g, &ms)

	return &Server{c, l, a, s, e, g}, nil
}

func (s *Server) Run(ctx context.Context) {
	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		s.l.Infow("startup", "config", s.c)
		s.l.Infow("starting http server at ", "address", s.c.Address)
		return s.s.ListenAndServe() //Запускаем сервер
	})
	//starting grpc server to listen on grpc address
	g.Go(func() error {
		listen, err := net.Listen("tcp", s.c.GrpcAddress)
		if err != nil {
			return err
		}
		s.l.Infow("starting grpc server at ", "address", s.c.GrpcAddress)
		return s.g.Serve(listen)
	})

	g.Go(func() error {
		<-gCtx.Done() //Ожидаем завершения контекста
		s.l.Info("shutting server down")
		//stopping grpc server
		s.g.Stop()
		return s.s.Shutdown(context.Background()) //Завершаем сервер
	})

	g.Wait()

	//Saving state to sync storage
	s.A.DoSync()
	s.A.Ds.Close()

	s.l.Warn("server stopped")
}
