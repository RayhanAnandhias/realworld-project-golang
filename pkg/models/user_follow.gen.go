// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package models

const TableNameUserFollow = "user_follow"

// UserFollow mapped from table <user_follow>
type UserFollow struct {
	IDUserA int32 `gorm:"column:id_user_a;type:integer;primaryKey" json:"id_user_a"`
	IDUserB int32 `gorm:"column:id_user_b;type:integer;primaryKey" json:"id_user_b"`
}

// TableName UserFollow's table name
func (*UserFollow) TableName() string {
	return TableNameUserFollow
}
