// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type TrackUser struct {
	ID      int32       `json:"id"`
	Name    string      `json:"name"`
	Counter pgtype.Int4 `json:"counter"`
}
