// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: posts.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createPost = `-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, published_at, feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, created_at, updated_at, title, url, published_at, feed_id
`

type CreatePostParams struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	Url         string
	PublishedAt time.Time
	FeedID      uuid.UUID
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, createPost,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Title,
		arg.Url,
		arg.PublishedAt,
		arg.FeedID,
	)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Url,
		&i.PublishedAt,
		&i.FeedID,
	)
	return i, err
}

const deletePost = `-- name: DeletePost :exec
DELETE FROM posts WHERE id = $1
`

func (q *Queries) DeletePost(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deletePost, id)
	return err
}

const getPosts = `-- name: GetPosts :many
SELECT posts.id, posts.created_at, posts.updated_at, posts.title, posts.url, posts.published_at, posts.feed_id
FROM feed_follows
JOIN feeds ON feed_follows.feed_id = feeds.id
JOIN posts ON feeds.id = posts.feed_id
and feed_follows.user_id = $1
`

func (q *Queries) GetPosts(ctx context.Context, userID uuid.UUID) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, getPosts, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Title,
			&i.Url,
			&i.PublishedAt,
			&i.FeedID,
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