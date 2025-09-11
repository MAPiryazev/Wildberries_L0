package main

import (
	"bufio"
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	models "github.com/MAPiryazev/Wildberries_L0/internal/model"
	"github.com/brianvoe/gofakeit/v7"
)

func main() {
	const totalOrders = 3000000
	const maxItemsPerOrder = 5

	file, err := os.Create("orders.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	rand.Seed(time.Now().UnixNano())
	gofakeit.Seed(time.Now().UnixNano())

	writer.Write([]byte("["))

	for i := 0; i < totalOrders; i++ {
		orderUID := "uid_" + strconv.Itoa(i)
		trackNumber := "WBILMTESTTRACK" + strconv.Itoa(i%1000)

		numItems := 1 + rand.Intn(maxItemsPerOrder)
		items := make([]models.Item, numItems)
		for j := 0; j < numItems; j++ {
			price := int32(gofakeit.Price(100, 1000))
			sale := int32(rand.Intn(51))
			totalPrice := price - (price*sale)/100

			items[j] = models.Item{
				OrderUID:    orderUID,
				ChrtID:      int64(gofakeit.Int32()),
				TrackNumber: trackNumber,
				Price:       price,
				Rid:         "rid_" + strconv.Itoa(i) + "_" + strconv.Itoa(j),
				Name:        gofakeit.ProductName(),
				Sale:        sale,
				Size:        strconv.Itoa(38 + rand.Intn(10)),
				TotalPrice:  totalPrice,
				NMID:        int64(1000000 + rand.Intn(9000000)),
				Brand:       gofakeit.Company(),
				Status:      202,
			}
		}

		amount := int32(0)
		for _, item := range items {
			amount += item.TotalPrice
		}
		if amount == 0 {
			amount = 1
		}

		var phone, email string
		switch rand.Intn(3) {
		case 0:
			phone = gofakeit.Phone()
		case 1:
			email = gofakeit.Email()
		default:
			phone = gofakeit.Phone()
			email = gofakeit.Email()
		}

		// 5% заказов будут битые
		isBroken := rand.Intn(100) < 5
		if isBroken {
			switch rand.Intn(5) {
			case 0:
				orderUID = "" // битый OrderUID
			case 1:
				phone = "" // битый контакт
			case 2:
				email = "" // битый контакт
			case 3:
				amount = -100 // битая сумма
			case 4:
				items = nil // пустой слайс товаров
			}
		}

		order := models.Order{
			OrderUID:    orderUID,
			TrackNumber: trackNumber,
			Entry:       "WBIL",
			Delivery: models.Delivery{
				Name:    gofakeit.Name(),
				Phone:   phone,
				Zip:     gofakeit.Zip(),
				City:    gofakeit.City(),
				Address: gofakeit.Street(),
				Region:  gofakeit.State(),
				Email:   email,
			},
			Payment: models.Payment{
				Transaction:  "txn_" + strconv.Itoa(i),
				RequestID:    "req_" + strconv.Itoa(i),
				Currency:     "USD",
				Provider:     "wbpay",
				Amount:       amount,
				PaymentDT:    time.Now().Unix(),
				Bank:         "discount",
				DeliveryCost: 500,
				GoodsTotal:   amount,
				CustomFee:    0,
			},
			Items:             items,
			Locale:            getShortLocale(),
			InternalSignature: "",
			CustomerID:        "cust_" + strconv.Itoa(i),
			DeliveryService:   "dhl",
			ShardKey:          int32(i % 10),
			SmID:              int64(100 + i),
			DateCreated:       time.Now(),
			OofShard:          int32(i % 10),
		}

		data, err := json.Marshal(order)
		if err != nil {
			log.Fatal(err)
		}

		if i != 0 {
			writer.Write([]byte(","))
		}
		writer.Write(data)

		if i%10000 == 0 {
			writer.Flush()
			log.Printf("Сгенерировано заказов: %d\n", i)
		}
	}

	writer.Write([]byte("]"))
	writer.Flush()

	log.Println("Генерация 3 млн заказов завершена!")
}

func getShortLocale() string {
	locales := []string{"en", "ru", "es", "fr", "de", "it", "pt", "ja", "ko", "zh"}
	return locales[rand.Intn(len(locales))]
}
