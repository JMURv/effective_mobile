package model

import (
	"time"
)

type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type Song struct {
	ID          uint64    `json:"id"`
	Group       string    `json:"group"`
	Song        string    `json:"song"`
	ReleaseDate time.Time `json:"release_date"`
	Lyrics      []string  `json:"lyrics"`
	Link        string    `json:"link"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PaginatedSongs struct {
	Data        any   `json:"data"`
	Count       int64 `json:"count"`
	TotalPages  int   `json:"total_pages"`
	CurrentPage int   `json:"current_page"`
	HasNextPage bool  `json:"has_next_page"`
}
