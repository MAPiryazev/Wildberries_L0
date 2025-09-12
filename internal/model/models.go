package models

import (
	"time"
)

type Delivery struct {
	Name    string `json:"name" validate:"required"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip" validate:"required,numeric"`
	City    string `json:"city" validate:"required"`
	Address string `json:"address" validate:"required"`
	Region  string `json:"region" validate:"required"`
	Email   string `json:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction" validate:"required"`
	RequestID    string `json:"request" validate:"required"`
	Currency     string `json:"currency" validate:"required,oneof=RUB USD EUR"`
	Provider     string `json:"provider"`
	Amount       int32  `json:"amount" validate:"required"`
	PaymentDT    int64  `json:"payment_dt"`
	Bank         string `json:"bank" validate:"required"`
	DeliveryCost int32  `json:"delivery_cost" validate:"gte=0"`
	GoodsTotal   int32  `json:"goods_total" validate:"gte=0"`
	CustomFee    int32  `json:"custom_fee" validate:"gte=0"`
}

type Item struct {
	OrderUID    string `json:"order_uid" validate:"required"`
	ChrtID      int64  `json:"chrt_id" validate:"required"`
	TrackNumber string `json:"track_number" validate:"required"`
	Price       int32  `json:"price" validate:"required,gt=0"`
	Rid         string `json:"rid"`
	Name        string `json:"name" validate:"required"`
	Sale        int32  `json:"sale" validate:"gte=0,lte=100"`
	Size        string `json:"size" validate:"lte=10"`
	TotalPrice  int32  `json:"total_price" validate:"gte=0"`
	NMID        int64  `json:"nm_id"`
	Brand       string `json:"brand" validate:"required"`
	Status      int32  `json:"status"`
}

type Order struct {
	OrderUID          string    `json:"order_uid" validate:"required"`
	TrackNumber       string    `json:"track_number" validate:"required"`
	Entry             string    `json:"entry"`
	Delivery          Delivery  `json:"delivery" validate:"required"`
	Payment           Payment   `json:"payment" validate:"required"`
	Items             []Item    `json:"items" validate:"required,min=1,dive"`
	Locale            string    `json:"locale" validate:"lte=10"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id" validate:"required"`
	DeliveryService   string    `json:"delivery_service" validate:"required"`
	ShardKey          int32     `json:"shardkey"`
	SmID              int64     `json:"sm_id" validate:"gte=0"`
	DateCreated       time.Time `json:"date_created" validate:"required"`
	OofShard          int32     `json:"oof_shard"`
}

type DLQMessage struct {
	OriginalValue []byte    `json:"original_value"`
	Error         string    `json:"error"`
	Timestamp     time.Time `json:"timestamp"`
}
