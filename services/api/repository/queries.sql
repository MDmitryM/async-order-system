-- name: CreateOrder :one
insert into orders (user_id, total, status, payment_method) 
values ($1, $2, $3, $4)
returning *;

-- name: GetOrderByID :one
select * from orders where id = $1;

-- name: UpdateOrderStatus :one
update orders set status = $2 where id = $1 returning *;

-- name: ListOrders :many
select * from orders order by created_at desc limit $1 offset $2;