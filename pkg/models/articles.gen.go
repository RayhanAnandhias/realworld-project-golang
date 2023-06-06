package models

import (
	"time"
)

const TableNameArticle = "articles"

// Article mapped from table <articles>
type Article struct {
	ID          int32     `gorm:"column:id;type:integer;primaryKey;autoIncrement:true" json:"id"`
	IDAuthor    int32     `gorm:"column:id_author;type:integer;not null" json:"id_author"`
	Slug        string    `gorm:"column:slug;type:text;not null" json:"slug"`
	Title       string    `gorm:"column:title;type:text;not null" json:"title"`
	Description string    `gorm:"column:description;type:text" json:"description"`
	Body        string    `gorm:"column:body;type:text;not null" json:"body"`
	CreatedAt   time.Time `gorm:"column:created_at;type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type ArticleRequest struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Body        string   `json:"body" binding:"required"`
	TagList     []string `json:"tagList,omitempty"`
}

type ArticleUpdate struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Body        string `json:"body,omitempty"`
}

type ArticleUpdateRequest struct {
	Article ArticleUpdate `json:"article" binding:"required"`
}

type ArticleCreateRequest struct {
	Article ArticleRequest `json:"article" binding:"required"`
}

type ArticleCommon struct {
	Slug           string       `json:"slug"`
	Title          string       `json:"title"`
	Description    string       `json:"description"`
	Body           string       `json:"body"`
	TagList        []string     `json:"tagList"`
	CreatedAt      time.Time    `json:"createdAt"`
	UpdatedAt      time.Time    `json:"updatedAt"`
	Favorited      bool         `json:"favorited"`
	FavoritesCount int32        `json:"favoritesCount"`
	Author         *UserProfile `json:"author"`
}

type ArticleResponse struct {
	Article *ArticleCommon `json:"article"`
}

type ArticleQueryResult struct {
	ID              int32     `json:"id"`
	Slug            string    `json:"slug"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Body            string    `json:"body"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	IDAuthor        int32     `json:"id_author"`
	Username        string    `json:"username"`
	Bio             *string   `json:"bio"`
	Image           *string   `json:"image"`
	Following       bool      `json:"following"`
	TagID           int32     `json:"tag_id"`
	TagName         string    `json:"tag_name"`
	LikedBy         int32     `json:"liked_by"`
	LikedByUsername string    `json:"liked_by_username"`
	Favorited       bool      `json:"favorited"`
	FavoritesCount  int32     `json:"favorites_count"`
}

type FavoritesCountQueryResult struct {
	FavoritesCount int32 `json:"favorites_count"`
}

// TableName Article's table name
func (*Article) TableName() string {
	return TableNameArticle
}
