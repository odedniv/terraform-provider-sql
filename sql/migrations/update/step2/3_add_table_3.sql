-- +migrate Up
CREATE TABLE table_3 (a integer, b integer);

-- +migrate Down
DROP TABLE table_3;
