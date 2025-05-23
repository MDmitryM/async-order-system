// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package repository

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Order struct {
	ID            int32              `json:"id"`
	UserID        int32              `json:"user_id"`
	Total         int32              `json:"total"`
	Status        string             `json:"status"`
	CreatedAt     pgtype.Timestamptz `json:"created_at"`
	PaymentMethod string             `json:"payment_method"`
	ProductID     int32              `json:"product_id"`
}
