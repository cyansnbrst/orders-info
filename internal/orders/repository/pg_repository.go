package repository

import (
	"database/sql"
	"errors"

	"cyansnbrst.com/order-info/internal/models"
	"cyansnbrst.com/order-info/internal/orders"
	"cyansnbrst.com/order-info/pkg/db"
)

// Orders repository
type ordersRepo struct {
	db *sql.DB
}

// Orders repository constructor
func NewOrdersRepository(db *sql.DB) orders.Repository {
	return &ordersRepo{db: db}
}

// Get order by UID
func (r *ordersRepo) Get(uid string) (*models.Order, error) {
	queryOrder := `
	SELECT 
		o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature, 
		o.customer_id, o.delivery_service, o.shardkey, o.sm_id, o.date_created, 
		o.oof_shard, 
		d.name AS delivery_name, d.phone AS delivery_phone, d.zip AS delivery_zip, 
		d.city AS delivery_city, d.address AS delivery_address, d.region AS delivery_region, 
		d.email AS delivery_email, 
		p.transaction AS payment_transaction, p.request_id AS payment_request_id, 
		p.currency AS payment_currency, p.provider AS payment_provider, p.amount AS payment_amount, 
		p.payment_dt AS payment_date, p.bank AS payment_bank, p.delivery_cost AS payment_delivery_cost, 
		p.goods_total AS payment_goods_total, p.custom_fee AS payment_custom_fee 
	FROM 
		orders o
	LEFT JOIN 
		delivery d ON o.order_uid = d.order_uid
	LEFT JOIN 
		payment p ON o.order_uid = p.order_uid
	WHERE 
		o.order_uid = $1;`

	var order models.Order
	err := r.db.QueryRow(queryOrder, uid).Scan(
		&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale,
		&order.InternalSignature, &order.CustomerID, &order.DeliveryService,
		&order.ShardKey, &order.SmID, &order.DateCreated, &order.OofShard,
		&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip,
		&order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region,
		&order.Delivery.Email, &order.Payment.Transaction, &order.Payment.RequestID,
		&order.Payment.Currency, &order.Payment.Provider, &order.Payment.Amount,
		&order.Payment.PaymentDt, &order.Payment.Bank, &order.Payment.DeliveryCost,
		&order.Payment.GoodsTotal, &order.Payment.CustomFee,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, db.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	queryItems := `
	SELECT 
		i.chrt_id, i.track_number, i.price, i.rid, i.name, i.sale, i.size, 
		i.total_price, i.nm_id, i.brand, i.status 
	FROM 
		items i 
	WHERE 
		i.order_uid = $1;`

	rows, err := r.db.Query(queryItems, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.Item
		if err := rows.Scan(
			&item.ChrtID, &item.TrackNumber, &item.Price, &item.RID,
			&item.Name, &item.Sale, &item.Size, &item.TotalPrice,
			&item.NmID, &item.Brand, &item.Status,
		); err != nil {
			return nil, err
		}
		order.Items = append(order.Items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &order, nil
}

// Save order to database
func (r *ordersRepo) Save(order *models.Order) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	queryOrder := `
	INSERT INTO orders (
		order_uid, track_number, entry, locale, internal_signature, 
		customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	ON CONFLICT (order_uid) DO UPDATE SET
		track_number = EXCLUDED.track_number,
		entry = EXCLUDED.entry,
		locale = EXCLUDED.locale,
		internal_signature = EXCLUDED.internal_signature,
		customer_id = EXCLUDED.customer_id,
		delivery_service = EXCLUDED.delivery_service,
		shardkey = EXCLUDED.shardkey,
		sm_id = EXCLUDED.sm_id,
		date_created = EXCLUDED.date_created,
		oof_shard = EXCLUDED.oof_shard;`

	_, err = tx.Exec(queryOrder,
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale,
		order.InternalSignature, order.CustomerID, order.DeliveryService,
		order.ShardKey, order.SmID, order.DateCreated, order.OofShard,
	)
	if err != nil {
		return err
	}

	queryDelivery := `
	INSERT INTO delivery (
		order_uid, name, phone, zip, city, address, region, email
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	ON CONFLICT (order_uid) DO UPDATE SET
		name = EXCLUDED.name,
		phone = EXCLUDED.phone,
		zip = EXCLUDED.zip,
		city = EXCLUDED.city,
		address = EXCLUDED.address,
		region = EXCLUDED.region,
		email = EXCLUDED.email;`

	_, err = tx.Exec(queryDelivery,
		order.OrderUID, order.Delivery.Name, order.Delivery.Phone,
		order.Delivery.Zip, order.Delivery.City, order.Delivery.Address,
		order.Delivery.Region, order.Delivery.Email,
	)
	if err != nil {
		return err
	}

	queryPayment := `
	INSERT INTO payment (
		order_uid, transaction, request_id, currency, provider, amount, 
		payment_dt, bank, delivery_cost, goods_total, custom_fee
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	ON CONFLICT (order_uid) DO UPDATE SET
		transaction = EXCLUDED.transaction,
		request_id = EXCLUDED.request_id,
		currency = EXCLUDED.currency,
		provider = EXCLUDED.provider,
		amount = EXCLUDED.amount,
		payment_dt = EXCLUDED.payment_dt,
		bank = EXCLUDED.bank,
		delivery_cost = EXCLUDED.delivery_cost,
		goods_total = EXCLUDED.goods_total,
		custom_fee = EXCLUDED.custom_fee;`

	_, err = tx.Exec(queryPayment,
		order.OrderUID, order.Payment.Transaction, order.Payment.RequestID,
		order.Payment.Currency, order.Payment.Provider, order.Payment.Amount,
		order.Payment.PaymentDt, order.Payment.Bank, order.Payment.DeliveryCost,
		order.Payment.GoodsTotal, order.Payment.CustomFee,
	)

	if err != nil {
		return err
	}

	queryDeleteItems := `DELETE FROM items WHERE order_uid = $1;`
	_, err = tx.Exec(queryDeleteItems, order.OrderUID)
	if err != nil {
		return err
	}

	queryItem := `
	INSERT INTO items (
		order_uid, chrt_id, track_number, price, rid, name, sale, size, 
		total_price, nm_id, brand, status
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);`

	for _, item := range order.Items {
		_, err = tx.Exec(queryItem,
			order.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.RID,
			item.Name, item.Sale, item.Size, item.TotalPrice,
			item.NmID, item.Brand, item.Status,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// Get all orders from database
func (r *ordersRepo) GetAll() ([]models.Order, error) {
	query := `
		SELECT 
			o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature, 
			o.customer_id, o.delivery_service, o.shardkey, o.sm_id, o.date_created, 
			o.oof_shard, 
			d.name AS delivery_name, d.phone AS delivery_phone, d.zip AS delivery_zip, 
			d.city AS delivery_city, d.address AS delivery_address, d.region AS delivery_region, d.email AS delivery_email, 
			p.transaction AS payment_transaction, p.request_id AS payment_request_id, 
			p.currency AS payment_currency, p.provider AS payment_provider, p.amount AS payment_amount, 
			p.payment_dt AS payment_date, p.bank AS payment_bank, p.delivery_cost AS payment_delivery_cost, 
			p.goods_total AS payment_goods_total, p.custom_fee AS payment_custom_fee,
			i.chrt_id, i.track_number AS item_track_number, i.price, i.rid, i.name AS item_name, 
			i.sale, i.size, i.total_price, i.nm_id, i.brand, i.status
		FROM 
			orders o
		LEFT JOIN 
			delivery d ON o.order_uid = d.order_uid
		LEFT JOIN 
			payment p ON o.order_uid = p.order_uid
		LEFT JOIN 
			items i ON o.order_uid = i.order_uid;
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	var currentOrder *models.Order

	for rows.Next() {
		var order models.Order
		var item models.Item
		if err := rows.Scan(
			&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale,
			&order.InternalSignature, &order.CustomerID, &order.DeliveryService,
			&order.ShardKey, &order.SmID, &order.DateCreated, &order.OofShard,
			&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip,
			&order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region,
			&order.Delivery.Email, &order.Payment.Transaction, &order.Payment.RequestID,
			&order.Payment.Currency, &order.Payment.Provider, &order.Payment.Amount,
			&order.Payment.PaymentDt, &order.Payment.Bank, &order.Payment.DeliveryCost,
			&order.Payment.GoodsTotal, &order.Payment.CustomFee,
			&item.ChrtID, &item.TrackNumber, &item.Price, &item.RID, &item.Name,
			&item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status,
		); err != nil {
			return nil, err
		}

		if currentOrder == nil || currentOrder.OrderUID != order.OrderUID {
			if currentOrder != nil {
				orders = append(orders, *currentOrder)
			}
			currentOrder = &order
			currentOrder.Items = []models.Item{item}
		} else {
			currentOrder.Items = append(currentOrder.Items, item)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if currentOrder != nil {
		orders = append(orders, *currentOrder)
	}

	return orders, nil
}
