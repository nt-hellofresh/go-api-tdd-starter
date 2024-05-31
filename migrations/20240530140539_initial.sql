-- +goose Up
-- +goose StatementBegin
CREATE TABLE product
(
    id    VARCHAR CONSTRAINT product_pk PRIMARY KEY,
    name  TEXT  NOT NULL,
    price FLOAT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE product;
-- +goose StatementEnd
