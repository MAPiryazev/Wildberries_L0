-- 1. Вставка справочников (если нет)
INSERT INTO city (name)
SELECT 'Kiryat Mozkin'
WHERE NOT EXISTS (SELECT 1 FROM city WHERE name = 'Kiryat Mozkin');

INSERT INTO region (name)
SELECT 'Kraiot'
WHERE NOT EXISTS (SELECT 1 FROM region WHERE name = 'Kraiot');

INSERT INTO brand (name)
SELECT 'Vivienne Sabo'
WHERE NOT EXISTS (SELECT 1 FROM brand WHERE name = 'Vivienne Sabo');

INSERT INTO bank (name)
SELECT 'alpha'
WHERE NOT EXISTS (SELECT 1 FROM bank WHERE name = 'alpha');

INSERT INTO currency (name)
SELECT 'USD'
WHERE NOT EXISTS (SELECT 1 FROM currency WHERE name = 'USD');

INSERT INTO delivery_service (name)
SELECT 'meest'
WHERE NOT EXISTS (SELECT 1 FROM delivery_service WHERE name = 'meest');


-- 2. Вставка delivery
INSERT INTO delivery (name, phone, zip, city_id, address, region_id, email)
SELECT 
    'Test Testov',
    '+9720000000',
    '2639809',
    (SELECT id FROM city WHERE name = 'Kiryat Mozkin'),
    'Ploshad Mira 15',
    (SELECT id FROM region WHERE name = 'Kraiot'),
    'test@gmail.com'
WHERE NOT EXISTS (
    SELECT 1 FROM delivery 
    WHERE name = 'Test Testov' AND phone = '+9720000000' AND zip = '2639809'
);


-- 3. Вставка payment
INSERT INTO payment (transaction, request_id, currency_id, provider, amount, payment_dt, bank_id, delivery_cost, goods_total, custom_fee)
SELECT 
    'b563feb7b2b84b6test',
    '',
    (SELECT id FROM currency WHERE name = 'USD'),
    'wbpay',
    1817,
    1637907727,
    (SELECT id FROM bank WHERE name = 'alpha'),
    1500,
    317,
    0
WHERE NOT EXISTS (
    SELECT 1 FROM payment WHERE transaction = 'b563feb7b2b84b6test'
);


-- 4. Вставка FactOrders
INSERT INTO FactOrders (
    order_uid, track_number, entry, delivery_id, payment_id, locale, internal_signature, customer_id, delivery_service_id, shardkey, sm_id, date_created, oof_shard
)
SELECT 
    'b563feb7b2b84b6test',
    'WBILMTESTTRACK',
    'WBIL',
    d.id,
    p.id,
    'en',
    '',
    'test',
    ds.id,
    9,
    99,
    '2021-11-26T06:22:19Z'::timestamp,
    1
FROM delivery d
JOIN payment p ON p.transaction = 'b563feb7b2b84b6test'
JOIN delivery_service ds ON ds.name = 'meest'
WHERE NOT EXISTS (SELECT 1 FROM FactOrders WHERE order_uid = 'b563feb7b2b84b6test');


-- 5. Вставка FactItems
INSERT INTO FactItems (
    order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand_id, status
)
SELECT
    'b563feb7b2b84b6test',
    9934930,
    'WBILMTESTTRACK',
    453,
    'ab4219087a764ae0btest',
    'Mascaras',
    30,
    '0',
    317,
    2389212,
    (SELECT id FROM brand WHERE name = 'Vivienne Sabo'),
    202
WHERE NOT EXISTS (
    SELECT 1 FROM FactItems WHERE order_uid = 'b563feb7b2b84b6test' AND chrt_id = 9934930
);
