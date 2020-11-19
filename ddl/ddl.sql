CREATE TABLE transactions(
    trade_id NUMERIC,
    time_stamp TIMESTAMP,
    pair VARCHAR,
    price NUMERIC,
    quantity NUMERIC,
    is_maker boolean,
    PRIMARY KEY(trade_id, time_stamp, pair, quantity, is_maker)
);

CREATE TABLE order_books(
	id NUMERIC NOT NULL,
	pair VARCHAR NOT NULL,
	time_stamp TIMESTAMP NOT NULL,
	bids jsonb NOT NULL,
	asks jsonb NOT NULL,
	PRIMARY KEY(id, pair, time_stamp)
);

