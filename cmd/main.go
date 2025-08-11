package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/MAPiryazev/Wildberries_L0/internal/db"
	"github.com/MAPiryazev/Wildberries_L0/internal/handlers"
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

	handler := handlers.NewHandler(orderService)
	router := handlers.RoutesInit(handler)

	fmt.Println("Starting server on :8081")
	err = http.ListenAndServe(":8081", router)
	if err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
