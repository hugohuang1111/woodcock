
-- change password
-- SET PASSWORD FOR 'root'@'localhost' = PASSWORD('newpass');

-- delete user
-- DROP USER 'poker'@'localhost';

-- create user
CREATE USER 'poker'@'localhost' IDENTIFIED BY 'poker1111';

-- create database
CREATE DATABASE poker COLLATE 'utf8_general_ci';
GRANT ALL ON poker.* TO 'poker'@'localhost';

-- select database
USE poker;

ALTER TABLE user AUTO_INCREMENT=10001;
-- create user table
CREATE TABLE IF NOT EXISTS user (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name varchar(128) NOT NULL,
    passwd varchar(32) NOT NULL,
    PRIMARY KEY(id),
    UNIQUE (name)
) ENGINE = innoDB DEFAULT CHARACTER SET = utf8;
