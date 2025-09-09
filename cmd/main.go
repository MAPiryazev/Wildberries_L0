package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/MAPiryazev/Wildberries_L0/internal/config"
	"github.com/MAPiryazev/Wildberries_L0/internal/db"
	"github.com/MAPiryazev/Wildberries_L0/internal/handlers"
	"github.com/MAPiryazev/Wildberries_L0/internal/kafka"
	"github.com/MAPiryazev/Wildberries_L0/internal/repository"
	"github.com/MAPiryazev/Wildberries_L0/internal/service"
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
	orderService, err := service.NewOrderService(orderRepo, 100)
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

	fmt.Println("Запускаем api на порту ", APIPort)
	err = http.ListenAndServe(":"+APIPort, router)
	if err != nil {
		log.Fatal("Неудачный запуск api: ", err)
	}
}
