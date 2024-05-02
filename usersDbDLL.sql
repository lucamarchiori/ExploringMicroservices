CREATE SCHEMA IF NOT EXISTS users;
CREATE TABLE IF NOT exists users.users (
	id serial4 NOT NULL,
	"name" varchar NULL,
	surname varchar NULL,
	email varchar NULL,
	"password" varchar NULL,
	created_at timestamp DEFAULT CURRENT_DATE NOT NULL
);