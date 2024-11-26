CREATE USER orders_user WITH PASSWORD 'orders_password';

CREATE DATABASE orders OWNER orders_user;

GRANT ALL PRIVILEGES ON DATABASE orders TO orders_user;