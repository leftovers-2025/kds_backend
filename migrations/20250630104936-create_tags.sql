
-- +migrate Up
CREATE TABLE `tags`(
    `id` BINARY(16) PRIMARY KEY COMMENT "タグID",
    `name` VARCHAR(255) NOT NULL COMMENT "タグ名"
);

-- +migrate Down
DROP TABLE `tags`;
