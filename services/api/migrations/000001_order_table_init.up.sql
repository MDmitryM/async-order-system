create table orders
(
    id serial primary key,
    user_id int not null,
    total int not null,
    status varchar(50) not null default 'Pending',
    created_at timestamp with time zone default current_timestamp
);