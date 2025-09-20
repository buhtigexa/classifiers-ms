-- Create a new UTF-8 `classifiersdb` database.
-- CREATE DATABASE classifiersdb CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Switch to using the `classifiersdb` database.
USE classifiersdb;


create table if not exists classifiers (
	id integer not null primary key auto_increment,
	name varchar(100) not null,
	description text null,
	is_active boolean null default true,
	created_at datetime not null default current_timestamp
);




CREATE USER 'appuser'@'localhost';
GRANT SELECT, INSERT, UPDATE , DELETE ON classifiersdb.* TO 'appuser'@'localhost';
ALTER USER 'appuser'@'localhost' IDENTIFIED by 'websecret';
