package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/MAPiryazev/Wildberries_L0/internal/config"
	"github.com/MAPiryazev/Wildberries_L0/internal/db"
	"github.com/MAPiryazev/Wildberries_L0/internal/handlers"
	"github.com/MAPiryazev/Wildberries_L0/internal/kafka"
	"github.com/MAPiryazev/Wildberries_L0/internal/repository"
	"github.com/MAPiryazev/Wildberries_L0/internal/service"
	"github.com/MAPiryazev/Wildberries_L0/internal/shutdown"
	kgo "github.com/segmentio/kafka-go"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	DBConfig, err := config.LoadDBConfig()
	if err != nil {
		switch {
		case errors.Is(err, config.ErrEnvNotFound):
			log.Panicln(err)
		case errors.Is(err, config.ErrParamNotFound):
			log.Panicln("Не все параметры подключения к БД были распознаны", err)
		default:
			log.Panicln("Неизвестная ошибка при получении конфига БД: ", err)
		}
	}

	psqlDB, err := db.InitPsqlDB(DBConfig)
	if err != nil {
		log.Fatal("Ошибка инициализации БД: ", err)
	}
	defer psqlDB.Close()

	orderRepo := repository.NewOrderRepo(psqlDB)

	dlqWriter := &kgo.Writer{
		Addr:     kgo.TCP("localhost:29092"),
		Topic:    "orders_DLQ",
		Balancer: &kgo.LeastBytes{},
	}

	orderService, err := service.NewOrderService(orderRepo, 100, dlqWriter)
	if err != nil {
		log.Fatal(err)
	}

	consumer, err := kafka.NewOrderConsumer(orderService)
	if err != nil {
		switch {
		case errors.Is(err, config.ErrEnvNotFound):
			log.Panicln(err)
		case errors.Is(err, config.ErrKafkaParamNotFound):
			log.Panicln("Не все параметры косьюмера кафки были распознаны", err)
		default:
			log.Println("Неизвестная ошибка при получении конфига консьюмера кафки: ", err)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go consumer.Start(ctx)

	handler := handlers.NewHandler(orderService)
	router := handlers.RoutesInit(handler)

	APIConfig, err := config.LoadAPIConfig()
	if err != nil {
		switch {
		case errors.Is(err, config.ErrEnvNotFound):
			log.Fatal("Фатальная ошибка, .env файл отсутствует:", err)
		default:
			fmt.Println(err)
		}
	}
	APIPort := APIConfig.APIPort

	srv := &http.Server{
		Addr:    ":" + APIPort,
		Handler: router,
	}
	go func() {
		fmt.Println("Запускаем API на порту", APIPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Panicln("Ошибка запуска API: ", err)
		}
	}()

	rootContext, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-rootContext.Done()
	shutdown.GracefulShutdown(cancel, srv, psqlDB)
}
