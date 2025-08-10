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
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Инициализируем БД
	psqlDB := db.InitPsqlDB()
	defer psqlDB.Close()

	// Создаем репозиторий
	orderRepo := repository.NewOrderRepo(psqlDB)

	// Тестовый заказ
	testOrder := models.Order{
		OrderUID:    "a1b2c3d4-e5f6-7890-1234-56789abcdef0", // новый UUID
		TrackNumber: "TRACK1234567890",
		Entry:       "WEB",
		Delivery: models.Delivery{
			Name:    "Доставка Плюс",
			Phone:   "+7 999 888 77 66",
			Zip:     "123456",
			City:    "Москва",
			Address: "ул. Новый Арбат, д. 10",
			Region:  "Москва",
			Email:   "delivery@example.com",
		},
		Payment: models.Payment{
			Transaction:  "tx9876543210",
			RequestID:    "req0987654321",
			Currency:     "RUB",
			Provider:     "Sberbank",
			Amount:       14999,
			PaymentDT:    1688700000,
			Bank:         "Сбербанк",
			DeliveryCost: 350,
			GoodsTotal:   14649,
			CustomFee:    0,
		},
		DeliveryService:   "DHL",
		Locale:            "ru",
		InternalSignature: "sig0987654321",
		CustomerID:        "cust123456",
		ShardKey:          5,
		SmID:              1234567890,
		DateCreated:       time.Date(2025, 8, 10, 22, 30, 0, 0, time.UTC),
		OofShard:          1,
		Items: []models.Item{
			{
				ChrtID:      555555,
				TrackNumber: "TRACK1234567890",
				Price:       4999,
				Rid:         "RID555",
				Name:        "Новые Кроссовки",
				Sale:        15,
				Size:        "42",
				TotalPrice:  4249,
				NMID:        111111,
				Brand:       "Nike",
				Status:      1,
			},
			{
				ChrtID:      666666,
				TrackNumber: "TRACK1234567890",
				Price:       9999,
				Rid:         "RID666",
				Name:        "Куртка зимняя",
				Sale:        5,
				Size:        "L",
				TotalPrice:  9499,
				NMID:        222222,
				Brand:       "Adidas",
				Status:      1,
			},
		},
	}

	// Сохраняем заказ
	fmt.Println("Сохраняем заказ...")
	err := orderRepo.SaveOrder(&testOrder)
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
