-- +migrate Up
CREATE TABLE table_1 (a integer, b integer);

-- +migrate Down
DROP TABLE table_1;
