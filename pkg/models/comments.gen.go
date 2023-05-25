// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

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

// TableName Comment's table name
func (*Comment) TableName() string {
	return TableNameComment
}
