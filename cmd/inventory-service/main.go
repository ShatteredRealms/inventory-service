package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/ShatteredRealms/go-common-service/pkg/log"
	"github.com/ShatteredRealms/go-common-service/pkg/pb"
	"github.com/ShatteredRealms/go-common-service/pkg/srv"
	"github.com/ShatteredRealms/go-common-service/pkg/telemetry"
	"github.com/ShatteredRealms/go-common-service/pkg/util"
	"github.com/ShatteredRealms/inventory-service/pkg/config"
	"github.com/WilSimpson/gocloak/v13"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	tracer := otel.Tracer("InventoryService")
	ctx, span := tracer.Start(ctx, "main")
	defer span.End()

	cfg, err := config.NewInventoryConfig(ctx)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("loading config: %v", err)
		return
	}

	otelShutdown, err := telemetry.SetupOTelSDK(ctx, "inventory", config.Version, cfg.OpenTelemtryAddress)
	defer func() {
		log.Logger.Infof("Shutting down")
		err = otelShutdown(context.Background())
		if err != nil {
			log.Logger.Warnf("Error shutting down: %v", err)
		}
	}()

	if err != nil {
		log.Logger.WithContext(ctx).Errorf("connecting to otel: %v", err)
		return
	}

	keycloakClient := gocloak.NewClient(cfg.Keycloak.BaseURL)
	grpcServer, gwmux := util.InitServerDefaults(keycloakClient, cfg.Keycloak.Realm)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	pb.RegisterHealthServiceServer(grpcServer, srv.NewHealthServiceServer())
	err = pb.RegisterHealthServiceHandlerFromEndpoint(ctx, gwmux, cfg.Server.Address(), opts)
	if err != nil {
		log.Logger.WithContext(ctx).Errorf("register health service handler endpoint: %v", err)
		return
	}

	span.End()
	srvErr := make(chan error, 1)
	go func() {
		srvErr <- util.StartServer(ctx, grpcServer, gwmux, cfg.Server.Address())
	}()

	select {
	case err = <-srvErr:
		log.Logger.Errorf("listen server: %v", err)

	case <-ctx.Done():
		log.Logger.Info("Server canceled by user input.")
		stop()
	}

}
