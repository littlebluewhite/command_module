USE schedule;

CREATE TABLE IF NOT EXISTS `time_templates`
(
    `id`           int PRIMARY KEY AUTO_INCREMENT,
    `name`         varchar(255) UNIQUE NOT NULL,
    `time_data_id` int UNIQUE          NOT NULL,
    `updated_at`   datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `created_at`   datetime DEFAULT (now())
    );



CREATE TABLE IF NOT EXISTS `time_data`
(
    `id`               int PRIMARY KEY AUTO_INCREMENT,
    `repeat_type`      ENUM ('daily', 'weekly', 'monthly'),
    `start_date`       date NOT NULL,
    `end_date`         date,
    `start_time`       time NOT NULL,
    `end_time`         time NOT NULL,
    `interval_seconds` int,
    `condition_type`   ENUM ('monthly_day', 'weekly_day', 'weekly_first', 'weekly_second', 'weekly_third', 'weekly_fourth', '') DEFAULT '',
    `condition`        json
    );



CREATE TABLE IF NOT EXISTS `schedules`
(
    `id`           int PRIMARY KEY AUTO_INCREMENT,
    `name`         varchar(255) UNIQUE NOT NULL,
    `description`  varchar(255),
    `time_data_id` int UNIQUE          NOT NULL,
    `command_id`   int,
    `enabled`      boolean  DEFAULT false,
    `updated_at`   datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `created_at`   datetime DEFAULT (now())
    );


CREATE TABLE IF NOT EXISTS `commands`
(
    `id`          int PRIMARY KEY AUTO_INCREMENT,
    `name`        varchar(255) UNIQUE NOT NULL,
    `protocol`    ENUM ('http', 'socket', 'websocket', 'snmp'),
    `description` varchar(255),
    `updated_at`  datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `created_at`  datetime DEFAULT (now())
    );


CREATE TABLE IF NOT EXISTS `https_commands`
(
    `id`                 int PRIMARY KEY AUTO_INCREMENT,
    `command_id`         int UNIQUE,
    `method`             ENUM ('GET', 'POST', 'PATCH', 'PUT', 'DELETE'),
    `url`                varchar(255) NOT NULL,
    `authorization_type` ENUM ('basic', 'token', '')                                                    DEFAULT '',
    `header`             json,
    `body_type`          ENUM ('text', 'html', 'xml', 'form_data', 'x-www_form_urlencoded', 'json', '') DEFAULT '',
    `body`               json,
    `parser_id`          int UNIQUE
    );


CREATE TABLE IF NOT EXISTS `header_templates`
(
    `id`   int PRIMARY KEY AUTO_INCREMENT,
    `name` varchar(255) UNIQUE NOT NULL,
    `data` json
    );


CREATE TABLE IF NOT EXISTS `snmp_commands`
(
    `id`         int PRIMARY KEY AUTO_INCREMENT,
    `command_id` int UNIQUE
);


CREATE TABLE IF NOT EXISTS `socket_commands`
(
    `id`         int PRIMARY KEY AUTO_INCREMENT,
    `command_id` int UNIQUE
);


CREATE TABLE IF NOT EXISTS `websocket_commands`
(
    `id`         int PRIMARY KEY AUTO_INCREMENT,
    `command_id` int UNIQUE
);


CREATE TABLE IF NOT EXISTS `parser`
(
    `id`             int PRIMARY KEY AUTO_INCREMENT,
    `search_command` varchar(255) COMMENT 'ex: person.item.[]array.name'
    );


ALTER TABLE `time_templates`
    ADD FOREIGN KEY (`time_data_id`) REFERENCES `time_data` (`id`);


ALTER TABLE `schedules`
    ADD FOREIGN KEY (`time_data_id`) REFERENCES `time_data` (`id`);


ALTER TABLE `schedules`
    ADD FOREIGN KEY (`command_id`) REFERENCES `commands` (`id`);


ALTER TABLE `https_commands`
    ADD FOREIGN KEY (`command_id`) REFERENCES `commands` (`id`);


ALTER TABLE `https_commands`
    ADD FOREIGN KEY (`parser_id`) REFERENCES `parser` (`id`);


ALTER TABLE `snmp_commands`
    ADD FOREIGN KEY (`command_id`) REFERENCES `commands` (`id`);


ALTER TABLE `socket_commands`
    ADD FOREIGN KEY (`command_id`) REFERENCES `commands` (`id`);


ALTER TABLE `websocket_commands`
    ADD FOREIGN KEY (`command_id`) REFERENCES `commands` (`id`);
