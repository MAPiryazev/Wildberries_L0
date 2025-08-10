package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/MAPiryazev/Wildberries_L0/internal/db"
	models "github.com/MAPiryazev/Wildberries_L0/internal/model"
	"github.com/MAPiryazev/Wildberries_L0/internal/repository"
)

func main() {
	// Инициализируем БД
	psqlDB := db.InitPsqlDB()
	defer psqlDB.Close()

	// Создаем репозиторий
	orderRepo := repository.NewOrderRepo(psqlDB)

	// Тестовый заказ
	testOrder := &models.Order{
		OrderUID:    "b563feb7b2b84b6test",
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: models.Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: models.Payment{
			Transaction:  "b563feb7b2b84b6test",
			RequestID:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDT:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Items: []models.Item{
			{
				ChrtID:      9934930,
				TrackNumber: "WBILMTESTTRACK",
				Price:       453,
				Rid:         "ab4219087a764ae0btest",
				Name:        "Mascaras",
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NMID:        2389212,
				Brand:       "Vivienne Sabo",
				Status:      202,
			},
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test",
		DeliveryService:   "meest",
		ShardKey:          9,
		SmID:              99,
		DateCreated:       time.Date(2021, 11, 26, 6, 22, 19, 0, time.UTC),
		OofShard:          1,
	}

	// Сохраняем заказ
	fmt.Println("Сохраняем заказ...")
	err := orderRepo.SaveOrder(testOrder)
	if err != nil {
		log.Fatalf("Ошибка сохранения заказа: %v", err)
	}
	fmt.Println("Заказ успешно сохранен!")

	// Получаем заказ обратно
	fmt.Println("\nПолучаем заказ обратно...")
	retrievedOrder, err := orderRepo.GetOrderById(testOrder.OrderUID)
	if err != nil {
		log.Fatalf("Ошибка получения заказа: %v", err)
	}

	// Выводим результат
	orderJSON, _ := json.MarshalIndent(retrievedOrder, "", "  ")
	fmt.Printf("Полученный заказ:\n%s\n", string(orderJSON))

	// Проверяем, что данные совпадают
	if retrievedOrder.OrderUID == testOrder.OrderUID &&
		retrievedOrder.TrackNumber == testOrder.TrackNumber &&
		retrievedOrder.Delivery.Name == testOrder.Delivery.Name {
		fmt.Println("\n✅ Тест прошел успешно! Данные совпадают.")
	} else {
		fmt.Println("\n❌ Тест не прошел. Данные не совпадают.")
	}
}
