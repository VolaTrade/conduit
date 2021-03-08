CREATE TABLE order_books(
	id NUMERIC NOT NULL,
	pair VARCHAR NOT NULL,
	time_stamp TIMESTAMP NOT NULL,
	bids jsonb NOT NULL,
	asks jsonb NOT NULL,
	PRIMARY KEY(id, pair, time_stamp)
);