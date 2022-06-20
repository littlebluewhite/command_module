-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS  `time_template`
(
    `id`               int PRIMARY KEY AUTO_INCREMENT,
    `name`             varchar(255) UNIQUE NOT NULL,
    `repeat_type`      ENUM ('daily', 'weekly', 'monthly'),
    `start_date`       date                NOT NULL,
    `end_date`         date,
    `start_time`       time                NOT NULL,
    `end_time`         time                NOT NULL,
    `interval_seconds` int,
    `updated_at`       datetime,
    `created_at`       datetime DEFAULT (now())
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS  `weekly_repeat`
(
    `id`               int PRIMARY KEY AUTO_INCREMENT,
    `time_template_id` int,
    `weekly_condition` json
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS  `monthly_repeat`
(
    `id`                    int PRIMARY KEY AUTO_INCREMENT,
    `time_template_id`      int,
    `first_week_condition`  json,
    `second_week_condition` json,
    `third_week_condition`  json,
    `fourth_week_condition` json,
    `monthly_condition`     json
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS  `schedule`
(
    `id`               int PRIMARY KEY AUTO_INCREMENT,
    `name`             varchar(255) UNIQUE NOT NULL,
    `description`      varchar(255),
    `time_template_id` int,
    `enable`           boolean DEFAULT false
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS  `schedule_command_ref`
(
    `schedule_id` int,
    `command_id`  int
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS  `command`
(
    `id`          int PRIMARY KEY AUTO_INCREMENT,
    `name`        varchar(255) UNIQUE NOT NULL,
    `protocol`    ENUM ('http', 'socket', 'websocket', 'snmp'),
    `description` varchar(255),
    `updated_at`  datetime,
    `created_at`  datetime DEFAULT (now())
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS  `https_command`
(
    `id`                 int PRIMARY KEY AUTO_INCREMENT,
    `command_id`         int UNIQUE,
    `method`             ENUM ('get', 'post', 'patch', 'put', 'delete'),
    `url`                varchar(255) NOT NULL,
    `authorization_type` ENUM ('basic', 'token'),
    `header`             json,
    `body_type`          ENUM ('plain', 'formData', 'xWwwFormUrlencoded', 'json', 'binary'),
    `body`               varchar(255),
    `parser_id`          int UNIQUE
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS  `header_template`
(
    `id`   int PRIMARY KEY AUTO_INCREMENT,
    `data` json
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS  `snmp_command`
(
    `id`         int PRIMARY KEY AUTO_INCREMENT,
    `command_id` int UNIQUE
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS  `socket_command`
(
    `id`         int PRIMARY KEY AUTO_INCREMENT,
    `command_id` int UNIQUE
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS  `websocket_command`
(
    `id`         int PRIMARY KEY AUTO_INCREMENT,
    `command_id` int UNIQUE
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS  `parser`
(
    `id`             int PRIMARY KEY AUTO_INCREMENT,
    `search_command` varchar(255) COMMENT 'ex: person.item.[]array.name'
);
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE `weekly_repeat`
    ADD FOREIGN KEY (`time_template_id`) REFERENCES `time_template` (`id`);
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE `monthly_repeat`
    ADD FOREIGN KEY (`time_template_id`) REFERENCES `time_template` (`id`);
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE `schedule`
    ADD FOREIGN KEY (`time_template_id`) REFERENCES `time_template` (`id`);
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE `schedule_command_ref`
    ADD FOREIGN KEY (`schedule_id`) REFERENCES `schedule` (`id`);
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE `schedule_command_ref`
    ADD FOREIGN KEY (`command_id`) REFERENCES `command` (`id`);
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE `https_command`
    ADD FOREIGN KEY (`command_id`) REFERENCES `command` (`id`);
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE `https_command`
    ADD FOREIGN KEY (`parser_id`) REFERENCES `parser` (`id`);
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE `snmp_command`
    ADD FOREIGN KEY (`command_id`) REFERENCES `command` (`id`);
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE `socket_command`
    ADD FOREIGN KEY (`command_id`) REFERENCES `command` (`id`);
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE `websocket_command`
    ADD FOREIGN KEY (`command_id`) REFERENCES `command` (`id`);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS weekly_repeat;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS monthly_repeat;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS schedule_command_ref;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS https_command;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS parser;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS header_template;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS socket_command;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS snmp_command;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS websocket_command;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS command;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS schedule;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS time_template;
-- +goose StatementEnd
