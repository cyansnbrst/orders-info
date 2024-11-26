CREATE TABLE orders (
    order_uid TEXT PRIMARY KEY,
    track_number TEXT UNIQUE NOT NULL,
    entry TEXT,
    locale TEXT,
    internal_signature TEXT,
    customer_id TEXT NOT NULL,
    delivery_service TEXT,
    shardkey TEXT,
    sm_id INT,
    date_created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    oof_shard TEXT
);

CREATE TABLE delivery (
    id SERIAL PRIMARY KEY,
    order_uid TEXT REFERENCES orders(order_uid) ON DELETE CASCADE UNIQUE,
    name TEXT NOT NULL,
    phone TEXT NOT NULL,
    zip TEXT,
    city TEXT,
    address TEXT,
    region TEXT,
    email TEXT
);

CREATE TABLE payment (
    id SERIAL PRIMARY KEY,
    order_uid TEXT REFERENCES orders(order_uid) ON DELETE CASCADE UNIQUE,
    transaction TEXT NOT NULL,
    request_id TEXT,
    currency TEXT NOT NULL,
    provider TEXT,
    amount INT NOT NULL,
    payment_dt BIGINT NOT NULL,
    bank TEXT,
    delivery_cost INT,
    goods_total INT,
    custom_fee INT
);

CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    order_uid TEXT REFERENCES orders(order_uid) ON DELETE CASCADE,
    chrt_id INT NOT NULL,
    track_number TEXT NOT NULL,
    price INT NOT NULL,
    rid TEXT,
    name TEXT NOT NULL,
    sale INT,
    size TEXT,
    total_price INT NOT NULL,
    nm_id INT,
    brand TEXT,
    status INT NOT NULL
);
