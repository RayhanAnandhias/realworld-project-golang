package models

import (
	"time"
)

const TableNameComment = "comments"

// Comment mapped from table <comments>
type Comment struct {
	ID        int32     `gorm:"column:id;type:integer;primaryKey;autoIncrement:true" json:"id"`
	IDAuthor  int32     `gorm:"column:id_author;type:integer;not null" json:"id_author"`
	IDArticle int32     `gorm:"column:id_article;type:integer;not null" json:"id_article"`
	Body      string    `gorm:"column:body;type:text;not null" json:"body"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type CommentRequest struct {
	Body string `json:"body" binding:"required"`
}

type CommentCreateRequest struct {
	Comment CommentRequest `json:"comment" binding:"required"`
}

type CommentQueryResult struct {
	ID        int32     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	IDAuthor  int32     `json:"id_author"`
	IDArticle int32     `json:"id_article"`
	Username  string    `json:"username"`
	Bio       *string   `json:"bio"`
	Image     *string   `json:"image"`
	Following bool      `json:"following"`
}

type CommentResponse struct {
	ID        int32        `json:"id"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
	Body      string       `json:"body"`
	Author    *UserProfile `json:"author"`
}

type CommentResponseData struct {
	Comment *CommentResponse `json:"comment"`
}

// TableName Comment's table name
func (*Comment) TableName() string {
	return TableNameComment
}
