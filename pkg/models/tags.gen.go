package models

import (
	"time"
)

const TableNameTag = "tags"

// Tag mapped from table <tags>
type Tag struct {
	ID        int32     `gorm:"column:id;type:integer;primaryKey;autoIncrement:true" json:"id"`
	Name      string    `gorm:"column:name;type:text;not null" json:"name"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type TagResponse struct {
	Tags []string `json:"tags"`
}

// TableName Tag's table name
func (*Tag) TableName() string {
	return TableNameTag
}
