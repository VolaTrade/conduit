CREATE TABLE transactions(
    trade_id UID DEFAULT NOT NULL,
    time_stamp TIMESTAMP, 
    pair VARCHAR(8),
    price NUMERIC,
    quantity NUMERIC,
    is_maker boolean
);
