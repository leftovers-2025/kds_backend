
-- +migrate Up
CREATE TABLE `notifications`(
    `user_id` BINARY(16) PRIMARY KEY COMMENT 'ユーザーID',
    `enabled` BOOLEAN NOT NULL DEFAULT FALSE COMMENT '有効化',
    FOREIGN KEY (`user_id`) REFERENCES `users`(`id`)
);

CREATE TABLE `tag_notifications`(
    `user_id` BINARY(16) NOT NULL COMMENT 'ユーザーID',
    `tag_id` BINARY(16) NOT NULL COMMENT 'タグID',
    FOREIGN KEY (`user_id`) REFERENCES `users`(`id`),
    FOREIGN KEY (`tag_id`) REFERENCES `tags`(`id`)
);

CREATE TABLE `location_notifications`(
    `user_id` BINARY(16) NOT NULL COMMENT 'ユーザーID',
    `location_id` BINARY(16) NOT NULL COMMENT 'ロケーションID',
    FOREIGN KEY (`user_id`) REFERENCES `users`(`id`),
    FOREIGN KEY (`location_id`) REFERENCES `locations`(`id`)
);

-- +migrate Down
DROP TABLE `location_notifications`;
DROP TABLE `tag_notifications`;
DROP TABLE `notifications`;
