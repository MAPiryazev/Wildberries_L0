package repository

import (
	"database/sql"

	models "github.com/MAPiryazev/Wildberries_L0/internal/model"
)

type OrderRepository interface {
	SaveOrder(order *models.Order) error
	GetOrderById(id string) (*models.Order, error)
}

type orderRepo struct {
	db *sql.DB
}

func NewOrderRepo(db *sql.DB) OrderRepository {
	return &orderRepo{db: db}
}

func (orderRepo *orderRepo) SaveOrder(order *models.Order) error {
	tx, err := orderRepo.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	cityID, err := orderRepo.cityIDFromName(tx, order.Delivery.City)
	if err != nil {
		return err
	}
	regionID, err := orderRepo.regionIDFromName(tx, order.Delivery.Region)
	if err != nil {
		return err
	}

	// Вставляем delivery
	var deliveryID int
	err = tx.QueryRow(`INSERT INTO delivery(name, phone, zip, city_id, address, region_id, email)
					   VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		cityID,
		order.Delivery.Address,
		regionID,
		order.Delivery.Email,
	).Scan(&deliveryID)
	if err != nil {
		return err
	}

	currencyID, err := orderRepo.currencyIDFromName(tx, order.Payment.Currency)
	if err != nil {
		return err
	}
	bankID, err := orderRepo.bankIDFromName(tx, order.Payment.Bank)
	if err != nil {
		return err
	}

	// Вставляем payment
	var paymentID int
	err = tx.QueryRow(`INSERT INTO payment(transaction, request_id, currency_id, provider, amount, payment_dt,
										   bank_id, delivery_cost, goods_total, custom_fee)
					   VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`,
		order.Payment.Transaction,
		order.Payment.RequestID,
		currencyID,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.PaymentDT,
		bankID,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFee,
	).Scan(&paymentID)
	if err != nil {
		return err
	}

	deliveryServiceID, err := orderRepo.deliveryServiceIDFromName(tx, order.DeliveryService)
	if err != nil {
		return err
	}

	// Вставляем основной заказ в factorders
	_, err = tx.Exec(`INSERT INTO factorders(order_uid, track_number, entry, delivery_id, payment_id, locale,
		internal_signature, customer_id, delivery_service_id, shardkey, sm_id, date_created, oof_shard)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		deliveryID,
		paymentID,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		deliveryServiceID,
		order.ShardKey,
		order.SmID,
		order.DateCreated,
		order.OofShard,
	)
	if err != nil {
		return err
	}

	// Вставляем items
	for _, item := range order.Items {
		brandID, err := orderRepo.brandIDFromName(tx, item.Brand)
		if err != nil {
			return err
		}
		_, err = tx.Exec(`INSERT INTO factitems(order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price,
		nm_id, brand_id, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
			order.OrderUID,
			item.ChrtID,
			item.TrackNumber,
			item.Price,
			item.Rid,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NMID,
			brandID,
			item.Status,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetOrderById implements OrderRepository.
func (o *orderRepo) GetOrderById(id string) (*models.Order, error) {
	panic("unimplemented")
}
