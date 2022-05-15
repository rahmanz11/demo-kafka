alter table match_order_info add column pay_with varchar(255);
alter table match_order_info add column put_proceeds varchar(255);
alter table match_order_info add column status varchar(255);
alter table match_order_info add column paid boolean default false;