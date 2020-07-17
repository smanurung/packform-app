create user pack_admin;
create database packform-db;
grant all privileges on database packform-db to pack_admin;
alter user pack_admin with encrypted password '<password>';

create table delivery (
    id serial primary key,
    order_item_id integer not null,
    delivered_quantity integer not null
);

create table order_items (
    id serial primary key,
    order_id integer not null,
    price_per_unit real not null,
    quantity integer not null,
    product varchar(128) not null
);

create table orders (
    id serial primary key,
    created_at timestamp not null,
    order_name varchar(128) not null,
    customer_id varchar(128) not null
);