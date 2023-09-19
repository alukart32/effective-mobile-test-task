CREATE TABLE IF NOT EXISTS "persons" (
    id uuid  PRIMARY KEY,
    name VARCHAR NOT NULL CHECK(LENGTH(name) > 0),
    surname VARCHAR NOT NULL CHECK(LENGTH(surname) > 0),
    patronymic VARCHAR NOT NULL,
    nation VARCHAR,
    gender VARCHAR,
    age INT CHECK(age > 0)
);