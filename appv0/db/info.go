package db

import "time"

// Info undefined
type Info struct {
	ID        int64     `json:"id" gorm:"id"`
	Name      string    `json:"name" gorm:"name"`
	CreatedAt time.Time `json:"created_at" gorm:"created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"updated_at"`
}

// TableName 表名称
func (*Info) TableName() string {
	return "info"
}
