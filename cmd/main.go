package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/MAPiryazev/Wildberries_L0/internal/config"
	"github.com/MAPiryazev/Wildberries_L0/internal/db"
	"github.com/MAPiryazev/Wildberries_L0/internal/handlers"
	"github.com/MAPiryazev/Wildberries_L0/internal/kafka"
	"github.com/MAPiryazev/Wildberries_L0/internal/repository"
	"github.com/MAPiryazev/Wildberries_L0/internal/service"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	psqlDB := db.InitPsqlDB()
	defer psqlDB.Close()

	orderRepo := repository.NewOrderRepo(psqlDB)
	orderService, err := service.NewOrderService(orderRepo, 100)
	if err != nil {
		log.Fatal(err)
	}

	consumer := kafka.NewOrderConsumer(orderService)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go consumer.Start(ctx)

	handler := handlers.NewHandler(orderService)
	router := handlers.RoutesInit(handler)

	APIPort := strconv.Itoa(config.LoadAPIConfig().APIPort)
	if APIPort == "" {
		log.Println("Не задан порт для API, выбираем 8081 в качестве дефолтного")
		APIPort = "8081"
	}

	fmt.Println("Запускаем api на порту ", APIPort)
	err = http.ListenAndServe(":"+APIPort, router)
	if err != nil {
		log.Fatal("Неудачный запуск api: ", err)
	}
}
