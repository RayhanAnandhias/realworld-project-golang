// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package models

const TableNameArticleTag = "article_tag"

// ArticleTag mapped from table <article_tag>
type ArticleTag struct {
	IDArticle int32 `gorm:"column:id_article;type:integer;primaryKey" json:"id_article"`
	IDTag     int32 `gorm:"column:id_tag;type:integer;primaryKey" json:"id_tag"`
}

// TableName ArticleTag's table name
func (*ArticleTag) TableName() string {
	return TableNameArticleTag
}