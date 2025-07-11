
-- +migrate Up
CREATE TABLE `posts`(
    `id` BINARY(16) PRIMARY KEY COMMENT '投稿ID',
    `user_id` BINARY(16) NOT NULL COMMENT 'ユーザーID',
    `location_id` BINARY(16) NOT NULL COMMENT 'ロケーションID',
    `description` VARCHAR(255) NOT NULL COMMENT '概要',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新日時',
    FOREIGN KEY (`user_id`) REFERENCES `users`(`id`),
    FOREIGN KEY (`location_id`) REFERENCES `locations`(`id`)
);

CREATE TABLE `post_tags`(
    `post_id` BINARY(16) NOT NULL COMMENT '投稿ID',
    `tag_id` BINARY(16) NOT NULL COMMENT 'タグID',
    FOREIGN KEY (`post_id`) REFERENCES `posts`(`id`),
    FOREIGN KEY (`tag_id`) REFERENCES `tags`(`id`),
    UNIQUE (`post_id`, `tag_id`)
);

CREATE TABLE `post_images`(
    `post_id` BINARY(16) NOT NULL COMMENT '投稿ID',
    `image_url` VARCHAR(255) NOT NULL COMMENT 'タグID',
    FOREIGN KEY (`post_id`) REFERENCES `posts`(`id`)
);

-- +migrate Down
DROP TABLE `post_images`;
DROP TABLE `post_tags`;
DROP TABLE `posts`;
