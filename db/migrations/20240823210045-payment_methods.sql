-- +migrate Up
CREATE TABLE
    `payment_methods` (
        `id` bigint NOT NULL AUTO_INCREMENT,
        `name` varchar(100) NOT NULL,
        `bank_code` varchar(50) NOT NULL,
        `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
        `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        `deleted_at` TIMESTAMP NULL DEFAULT NULL,
        PRIMARY KEY (`id`)
    );

-- +migrate Down
DROP TABLE IF EXISTS payment_methods;