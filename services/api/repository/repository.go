package repository

import "context"

type Repository interface {
	CreateOrder(ctx context.Context, params CreateOrderParams) (Order, error)
	GetOrderByID(ctx context.Context, id int32) (Order, error)
	ListOrders(ctx context.Context, params ListOrdersParams) ([]Order, error)
	UpdateOrderStatus(ctx context.Context, params UpdateOrderStatusParams) (Order, error)
	Close()
	Ping(ctx context.Context) error
}
