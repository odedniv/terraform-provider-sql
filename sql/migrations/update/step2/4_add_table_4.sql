-- +migrate Up
CREATE TABLE table_4 (a integer, b integer);

-- +migrate Down
DROP TABLE table_4;
