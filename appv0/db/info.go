package db

import (
	"fmt"
	"time"
)

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

// GetInfo 获取用户基本信息
func (i *Info) GetInfo(id int) (*Info, error) {
	var ret Info
	if err := DB.Table(i.TableName()).Where("id = ?", id).First(&ret).Error; err != nil {
		fmt.Println("用户信息：未找到数据", err)
		return nil, err // 返回错误给调用方
	}
	return &ret, nil
}
