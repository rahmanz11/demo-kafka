create table order_info(
	id serial PRIMARY KEY,
	order_id VARCHAR(255),
	order_type VARCHAR(255),
	amt INTEGER,
	_from VARCHAR(255),
	_to VARCHAR(255),
	pmt_method VARCHAR(255),
	matched BOOLEAN DEFAULT false,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

create table match_order_info(
	id serial PRIMARY KEY,
	order_id VARCHAR(255),
	order_type VARCHAR(255),
	amt INTEGER,
	_from VARCHAR(255),
	_to VARCHAR(255),
	pmt_method VARCHAR(255),
	sell_order_id VARCHAR(255),
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);