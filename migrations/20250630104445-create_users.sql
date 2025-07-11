
-- +migrate Up
CREATE TABLE `users`(
    `id` BINARY(16) PRIMARY KEY COMMENT 'ユーザーID',
    `email` VARCHAR(255) NOT NULL COMMENT 'メールアドレス',
    `name` VARCHAR(255) NOT NULL COMMENT 'フルネーム',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新日時'
);

CREATE TABLE `google_ids`(
    `user_id` BINARY(16) PRIMARY KEY COMMENT 'ユーザーID',
    `google_id` VARCHAR(255) NOT NULL UNIQUE COMMENT 'GoogleユーザーID',
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE `roles`(
    `user_id` BINARY(16) PRIMARY KEY COMMENT 'ユーザーID',
    `role` VARCHAR(255) NOT NULL COMMENT 'ロール',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新日時',
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- +migrate Down
DROP TABLE `roles`;
DROP TABLE `google_ids`;
DROP TABLE `users`;
