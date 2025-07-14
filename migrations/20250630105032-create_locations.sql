
-- +migrate Up
CREATE TABLE `locations`(
    `id` BINARY(16) PRIMARY KEY COMMENT "ロケーションID",
    `name` VARCHAR(255) NOT NULL UNIQUE COMMENT "ロケーション名"
);

-- +migrate Down
DROP TABLE `locations`;
