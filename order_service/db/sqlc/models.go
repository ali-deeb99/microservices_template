// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Order struct {
	ID     int32       `json:"id"`
	Name   string      `json:"name"`
	Note   pgtype.Text `json:"note"`
	Status int64       `json:"status"`
}
