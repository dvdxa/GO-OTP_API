package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"otp_api/conf"
	"otp_api/internal/app/adapter"
	"otp_api/internal/app/controller"
	"otp_api/internal/app/pkg/db/redis"
	"otp_api/internal/app/pkg/logger"
	"otp_api/internal/app/pkg/server"
	"otp_api/internal/app/storage"
	"otp_api/internal/app/usecase"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatalf("failed to init env %s", err.Error())
	}
}

func Run() {
	cfg := conf.InitConfigs()

	log := logger.InitLogger(&cfg.Logger)

	redisClient := redis.NewRedisClient(cfg, log)

	storage := storage.NewStorage(log, redisClient)

	adapter := adapter.NewAdapter(log, cfg)

	service := usecase.NewService(storage, adapter, log)

	handler := controller.NewHandler(cfg, log, service)

	srv := server.NewServer(cfg, handler)

	var wg sync.WaitGroup

	wg.Add(2)

	ctx, cancel := context.WithCancel(context.Background())

	go startServer(srv, ctx, &wg)
	waitForInterrupt(cancel, &wg)

	wg.Wait()
}

func startServer(srv *server.Server, ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	go func() {
		err := srv.Run()
		if err != nil {
			log.Fatalf("failed to run server: %v\n", err)
		}
	}()

	select {
	case <-ctx.Done():
		fmt.Println("shutting down server gracefully")

		//5 sec to gracefully shutdown server
		shutDownCtx, shutDowncancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutDowncancel()

		err := srv.Shutdown(shutDownCtx)
		if err != nil {
			fmt.Printf("error shutting 	down server %v\n", err.Error())
		}
	}
}

func waitForInterrupt(cancel context.CancelFunc, wg *sync.WaitGroup) {

	ch := make(chan os.Signal, 1)

	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		defer wg.Done()
		s := <-ch
		fmt.Printf("received signal %v\n", s)
		cancel()
	}()
}
