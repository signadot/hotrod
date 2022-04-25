CREATE DATABASE customer;
use customer;
CREATE TABLE customers
(
    id bigint unsigned NOT NULL,

    name varchar(255) NOT NULL,

    location varchar(255) NOT NULL,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE =utf8mb4_0900_ai_ci;
