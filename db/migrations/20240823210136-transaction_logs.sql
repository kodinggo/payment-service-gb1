-- +migrate Up
CREATE TABLE
    `transaction_logs` (
        id bigint not NULL AUTO_INCREMENT,
        transaction_id bigint not NULL,
        status varchar(50) not NULL,
        `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
        `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        `deleted_at` TIMESTAMP NULL DEFAULT NULL,
        PRIMARY KEY (`id`),
        FOREIGN KEY (`transaction_id`) REFERENCES transactions (`id`)
    );

-- +migrate Down
DROP TABLE IF EXISTS transaction_logs;