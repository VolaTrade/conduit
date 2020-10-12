CREATE TABLE transactions(
    trade_id UUID NOT NULL DEFAULT uuid_generate_v4 (),
    time_stamp TIMESTAMP,
    pair VARCHAR(8),
    price NUMERIC,
    quantity NUMERIC,
    is_maker boolean
);
