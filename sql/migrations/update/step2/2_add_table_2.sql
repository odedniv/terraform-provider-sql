-- +migrate Up
CREATE TABLE table_2 (c integer, d integer);

-- +migrate Down
DROP TABLE table_2;
