-- +goose Up
CREATE TABLE accounts (
	id TEXT NOT NULL DEFAULT,
	name TEXT NOT NULL
);
