CREATE TABLE transactions(
    trade_id NUMERIC,
    time_stamp TIMESTAMP,
    pair VARCHAR,
    price NUMERIC,
    quantity NUMERIC,
    is_maker boolean,
    PRIMARY KEY(trade_id, time_stamp, pair, quantity, is_maker)
);
