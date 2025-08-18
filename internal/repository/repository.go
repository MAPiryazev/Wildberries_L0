package repository

import (
	"database/sql"

	models "github.com/MAPiryazev/Wildberries_L0/internal/model"
)

type OrderRepository interface {
	SaveOrder(order *models.Order) error
	SaveOrdersBatch(orders []*models.Order) error
	GetOrderById(id string) (*models.Order, error)
	GetLastNOrders(limit int) ([]*models.Order, error)
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
func (orderRepo *orderRepo) GetOrderById(id string) (*models.Order, error) {
	var order models.Order

	// Основной заказ без factitems и brand
	err := orderRepo.db.QueryRow(`
        SELECT 
            fo.order_uid, fo.track_number, fo.entry, fo.locale, fo.internal_signature, fo.customer_id,
            fo.shardkey, fo.sm_id, fo.date_created, fo.oof_shard,
            d.name, d.phone, d.zip, d.address, d.email,
            c.name, r.name,
            p.transaction, p.request_id, cr.name, p.provider, p.amount, p.payment_dt,
            bn.name, p.delivery_cost, p.goods_total, p.custom_fee,
            ds.name
        FROM factorders fo
        LEFT JOIN delivery d ON fo.delivery_id = d.id
        LEFT JOIN city c ON d.city_id = c.id
        LEFT JOIN region r ON d.region_id = r.id
        LEFT JOIN payment p ON fo.payment_id = p.id
        LEFT JOIN currency cr ON p.currency_id = cr.id
        LEFT JOIN bank bn ON p.bank_id = bn.id
        LEFT JOIN delivery_service ds ON fo.delivery_service_id = ds.id
        WHERE fo.order_uid = $1
    `, id).Scan(
		&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerID,
		&order.ShardKey,
		&order.SmID,
		&order.DateCreated,
		&order.OofShard,
		&order.Delivery.Name,
		&order.Delivery.Phone,
		&order.Delivery.Zip,
		&order.Delivery.Address,
		&order.Delivery.Email,
		&order.Delivery.City,
		&order.Delivery.Region,
		&order.Payment.Transaction,
		&order.Payment.RequestID,
		&order.Payment.Currency,
		&order.Payment.Provider,
		&order.Payment.Amount,
		&order.Payment.PaymentDT,
		&order.Payment.Bank,
		&order.Payment.DeliveryCost,
		&order.Payment.GoodsTotal,
		&order.Payment.CustomFee,
		&order.DeliveryService,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Теперь подтягиваем items
	rows, err := orderRepo.db.Query(`
        SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand_id, status
        FROM factitems
        WHERE order_uid = $1
    `, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		var brandID int
		err := rows.Scan(
			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.Rid,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NMID,
			&brandID,
			&item.Status,
		)
		if err != nil {
			return nil, err
		}

		brandName, err := orderRepo.brandNameByID(brandID)
		if err != nil {
			return nil, err
		}
		item.Brand = brandName

		items = append(items, item)
	}
	order.Items = items

	return &order, nil
}

func (orderRepo *orderRepo) brandNameByID(id int) (string, error) {
	var name string
	err := orderRepo.db.QueryRow(`SELECT name FROM brand WHERE id = $1`, id).Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}

func (orderRepo *orderRepo) GetLastNOrders(limit int) ([]*models.Order, error) {
	rows, err := orderRepo.db.Query(`
        SELECT order_uid
        FROM factorders
        ORDER BY date_created DESC
        LIMIT $1
    `, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var uids []string
	for rows.Next() {
		var uid string
		if err := rows.Scan(&uid); err != nil {
			return nil, err
		}
		uids = append(uids, uid)
	}

	var orders []*models.Order
	for _, uid := range uids {
		order, err := orderRepo.GetOrderById(uid)
		if err != nil {
			return nil, err
		}
		if order != nil {
			orders = append(orders, order)
		}
	}

	return orders, nil
}

func (orderRepo *orderRepo) SaveOrdersBatch(orders []*models.Order) error {
	if len(orders) == 0 {
		return nil
	}

	tx, err := orderRepo.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, order := range orders {
		cityID, err := orderRepo.cityIDFromName(tx, order.Delivery.City)
		if err != nil {
			return err
		}

		regionID, err := orderRepo.regionIDFromName(tx, order.Delivery.Region)
		if err != nil {
			return err
		}

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
	}

	return tx.Commit()
}
