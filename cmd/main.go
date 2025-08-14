package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

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
	orderService, err := service.NewOrderService(orderRepo, 5)
	if err != nil {
		log.Fatal(err)
	}
	//TODO  вот это все в env запихнуть
	consumer := kafka.NewOrderConsumer("localhost:29092", "orders", "order-consumer-group", orderService)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go consumer.Start(ctx)

	handler := handlers.NewHandler(orderService)
	router := handlers.RoutesInit(handler)

	fmt.Println("Запускаем api на порту :8081")
	err = http.ListenAndServe(":8081", router)
	if err != nil {
		log.Fatal("Неудачный запуск api: ", err)
	}
}
