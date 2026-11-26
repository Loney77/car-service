package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	cargrpc "example.com/car-service/grpc"
	carhttp "example.com/car-service/http"
	"example.com/car-service/implementation"
	"example.com/car-service/pb"
	"example.com/car-service/repository"
	"example.com/car-service/transport"
)

func main() {
	httpAddr := flag.String("http.addr", ":8080", "HTTP listen address")
	grpcAddr := flag.String("grpc.addr", ":9090", "gRPC listen address")
	flag.Parse()

	// Настройка логгера
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = level.NewFilter(logger, level.AllowDebug())
		logger = log.With(logger,
			"svc", "car",
			"ts", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	level.Info(logger).Log("msg", "service started")
	defer level.Info(logger).Log("msg", "service ended")

	// Создание сервиса
	repo := repository.NewMemoryRepository()
	svc := implementation.NewService(repo, logger)

	// Создание HTTP-обработчика
	endpoints := transport.MakeEndpoints(svc)
	handler := carhttp.NewService(endpoints, logger)

	grpcServer := grpc.NewServer()
	defer grpcServer.Stop()
	pb.RegisterCarServiceServer(grpcServer, cargrpc.NewService(endpoints, logger))
	reflection.Register(grpcServer)

	// Запуск сервера
	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		level.Info(logger).Log("transport", "HTTP", "addr", *httpAddr)
		server := &http.Server{
			Addr:    *httpAddr,
			Handler: handler,
		}
		errs <- server.ListenAndServe()
	}()

	go func() {
		listener, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			errs <- err
			return
		}
		level.Info(logger).Log("transport", "gRPC", "addr", *grpcAddr)
		errs <- grpcServer.Serve(listener)
	}()

	level.Error(logger).Log("exit", <-errs)
}
