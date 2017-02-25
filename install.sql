-- Make sure your mysql process is running first!

DROP DATABASE IF EXISTS `SecurityNews`;
CREATE DATABASE IF NOT EXISTS `SecurityNews`  CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `SecurityNews`;

-- Enable client program to communicate with the server using utf8 character set
SET NAMES 'utf8';

DROP TABLE IF EXISTS `hackernews`;
create table IF NOT EXISTS `hackernews`(
    `id`      INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    `uuid`    varchar(36) not null,
    `title`    varchar(256) not null,
    `link`    varchar(256) not null,
    `post_date` timestamp  DEFAULT CURRENT_TIMESTAMP,
    `score`    varchar(64) not null,
    `user_name`    varchar(32) not null,
    `user_profile`    varchar(128) not null,
    `md5`     varchar(32) not null,
    UNIQUE(`uuid`),
    PRIMARY KEY (`id`)
)DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
alter table hackernews add unique index hackernews_md5_ux (md5);

EXPLAIN `hackernews`;

DROP TABLE IF EXISTS `threadpost`;
create table IF NOT EXISTS `threadpost`(
    `id`      INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    `title`    varchar(256) not null,
    `link`    varchar(256) not null,
    `post_date` timestamp  DEFAULT CURRENT_TIMESTAMP,
    `user_name`    varchar(32) not null,
    `user_profile`    varchar(128) not null,
    `short_content`   text,
    `cover_img`    varchar(256),
    `md5`     varchar(32) not null,
    PRIMARY KEY (`id`)
)DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
alter table threadpost add unique index threadpost_md5_ux (md5);

EXPLAIN `threadpost`;