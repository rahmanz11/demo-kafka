alter table match_order_info drop column paid;
alter table match_order_info drop column pmt_method;
alter table match_order_info rename column order_id to receipt_id;
alter table match_order_info rename column amt to fund_amt;
alter table match_order_info add column fund_receipt_id varchar(255);
alter table match_order_info add column transaction_type varchar(255);
alter table match_order_info add column currency_code varchar(255);
alter table match_order_info add column curreyncy_name varchar(255);
alter table match_order_info add column paid_at TIMESTAMP;