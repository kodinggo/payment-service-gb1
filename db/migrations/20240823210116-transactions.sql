-- +migrate Up
CREATE TABLE
    `transactions` (
        `id` bigint not NULL AUTO_INCREMENT,
        `user_id` bigint NOT NULL,
        `order_id` bigint NOT NULL,
        `payment_method_id` bigint not NULL,
        `status` enum ('success', 'failed') NOT NULL,
        `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
        `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        `deleted_at` TIMESTAMP NULL DEFAULT NULL,
        PRIMARY KEY (`id`),
        FOREIGN KEY (`payment_method_id`) REFERENCES payment_methods (`id`)
    );

-- +migrate Down
DROP TABLE IF EXISTS transactions;