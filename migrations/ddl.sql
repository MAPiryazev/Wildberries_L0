create table city(
    id serial primary key,
    name varchar(255) not null 
);

create table brand(
    id serial primary key,
    name varchar(255) not null
);

create table region(
    id serial primary key,
    name varchar(255) not null
);

create table bank(
    id serial primary key,
    name varchar(255) not null
);

create table currency(
    id serial primary key,
    name varchar(255) not null
);

create table delivery_service(
    id serial primary key,
    name varchar(255) not null
);

create table delivery(
    id serial primary key,
    name varchar(255) not null,
    phone varchar(20),
    zip varchar(20),
    city_id int not null, 
    address varchar(255) not null,
    region_id int not null, 
    email varchar(255),
    foreign key (city_id) references city(id) on delete restrict,
    foreign key (region_id) references region(id) on delete restrict
);

create table payment(
    id serial primary key, 
    transaction varchar(50) not null,
    request_id varchar(50),
    currency_id int,
    provider varchar(50),
    amount int,
    payment_dt bigint,
    bank_id int, 
    delivery_cost int,
    goods_total int,
    custom_fee int,
    foreign key (currency_id) references currency(id) on delete restrict,
    foreign key (bank_id) references bank(id) on delete restrict
);

create table FactOrders(
    id serial primary key,
    order_uid varchar(50) not null unique, 
    track_number varchar(50) not null,
    entry varchar(50),
    delivery_id int, 
    payment_id int, 
    locale varchar(10),
    internal_signature varchar(255), 
    customer_id varchar(50),
    delivery_service_id int, 
    shardkey int,
    sm_id bigint,
    date_created timestamp,
    oof_shard int,
    foreign key (delivery_id) references delivery(id) on delete cascade,
    foreign key (payment_id) references payment(id) on delete cascade,
    foreign key (delivery_service_id) references delivery_service(id) on delete restrict
);

create table FactItems(
    id serial primary key,
    order_uid varchar(50) not null,
    chrt_id bigint,
    track_number varchar(50) not null,
    price int,
    rid varchar(50),
    name varchar(255),
    sale int,
    size varchar(10),
    total_price int,
    nm_id bigint,
    brand_id int,
    status int,
    foreign key (brand_id) references brand(id) on delete restrict,
    foreign key (order_uid) references FactOrders(order_uid) on delete cascade
);
