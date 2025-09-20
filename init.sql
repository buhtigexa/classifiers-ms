-- Create a new UTF-8 `classifiersdb` database.
-- CREATE DATABASE classifiersdb CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Switch to using the `classifiersdb` database.
USE classifiersdb;


create table if not exists classifiers (
	id integer not null primary key auto_increment,
	name varchar(100) not null,
	created datetime not null
);




CREATE USER 'appuser'@'localhost';
GRANT SELECT, INSERT, UPDATE , DELETE ON classifiersdb.* TO 'appuser'@'localhost';
ALTER USER 'appuser'@'localhost' IDENTIFIED by 'websecret';
