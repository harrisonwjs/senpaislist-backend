// Code generated by sqlc. DO NOT EDIT.
// source: bookmarks.sql

package db

import (
	"context"
	"database/sql"
)

const createBookmark = `-- name: CreateBookmark :one
INSERT INTO bookmarks (
    owner,
    anime_id,
    bookmark_type
) VALUES (
    $1,$2,$3
) RETURNING owner, anime_id, bookmark_type, id, created_at
`

type CreateBookmarkParams struct {
	Owner        string        `json:"owner"`
	AnimeID      sql.NullInt64 `json:"anime_id"`
	BookmarkType string        `json:"bookmark_type"`
}

func (q *Queries) CreateBookmark(ctx context.Context, arg CreateBookmarkParams) (Bookmark, error) {
	row := q.db.QueryRowContext(ctx, createBookmark, arg.Owner, arg.AnimeID, arg.BookmarkType)
	var i Bookmark
	err := row.Scan(
		&i.Owner,
		&i.AnimeID,
		&i.BookmarkType,
		&i.ID,
		&i.CreatedAt,
	)
	return i, err
}
