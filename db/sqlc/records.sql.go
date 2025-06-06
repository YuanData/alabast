// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: records.sql

package db

import (
	"context"
	"encoding/json"
)

const createRecord = `-- name: CreateRecord :one
INSERT INTO records (
  username,
  content
) VALUES (
  $1, $2
) RETURNING id, username, content, created_at, updated_at
`

type CreateRecordParams struct {
	Username string          `json:"username"`
	Content  json.RawMessage `json:"content"`
}

func (q *Queries) CreateRecord(ctx context.Context, arg CreateRecordParams) (Record, error) {
	row := q.db.QueryRowContext(ctx, createRecord, arg.Username, arg.Content)
	var i Record
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Content,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getRecord = `-- name: GetRecord :one
SELECT id, username, content, created_at, updated_at FROM records
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetRecord(ctx context.Context, id int64) (Record, error) {
	row := q.db.QueryRowContext(ctx, getRecord, id)
	var i Record
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Content,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listRecords = `-- name: ListRecords :many
SELECT id, username, content, created_at, updated_at FROM records
WHERE username = $1
ORDER BY id
LIMIT $2
OFFSET $3
`

type ListRecordsParams struct {
	Username string `json:"username"`
	Limit    int32  `json:"limit"`
	Offset   int32  `json:"offset"`
}

func (q *Queries) ListRecords(ctx context.Context, arg ListRecordsParams) ([]Record, error) {
	rows, err := q.db.QueryContext(ctx, listRecords, arg.Username, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Record{}
	for rows.Next() {
		var i Record
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.Content,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
