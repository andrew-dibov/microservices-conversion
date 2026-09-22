package main

import (
	"errors"
	"microservices-conversion/internal/clients"
	"microservices-conversion/internal/configs"
	"microservices-conversion/internal/interceptors"
	"microservices-conversion/internal/loggers"
	"microservices-conversion/internal/servers"
	"microservices-conversion/pkg/api/conversion"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func main() {
	appConfig := configs.NewAppConfig()
	appLogger := loggers.NewAppLogger(appConfig)

	appLogger.Info("config",
		"port", appConfig.App.Port,
		"prod", appConfig.App.Prod)

	/* --- --- --- */

	currencyClient, err := clients.NewCurrencyClient(&appConfig)
	if err != nil {
		appLogger.Error("NewCurrencyClient returned error", "error", err)
		os.Exit(1)
	}
	defer currencyClient.Close()

	/* --- --- --- */

	grpcServer := grpc.NewServer(grpc.MaxRecvMsgSize(4*1024*1024), grpc.MaxSendMsgSize(4*1024*1024),
		grpc.UnaryInterceptor(interceptors.TraceInterceptor),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    appConfig.App.KeepaliveTime,
			Timeout: appConfig.App.KeepaliveTimeout,
		}))

	conversion.RegisterConversionServer(grpcServer, servers.NewAppServer(currencyClient, appLogger))

	appListener, err := net.Listen("tcp", ":"+appConfig.App.Port)
	if err != nil {
		appLogger.Error("appListener returned error", "error", err)
		os.Exit(1)
	}

	go func() {
		if err := grpcServer.Serve(appListener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			appLogger.Error("grpcServer returned error", "error", err)
			os.Exit(1)
		}
	}()

	/* --- --- --- */

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	done := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		appLogger.Info("server stopped")
	case <-time.After(appConfig.App.ShutdownTimeout):
		appLogger.Warn("server forced to stop")
		grpcServer.Stop()
	}
}
