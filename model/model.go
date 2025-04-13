package model

import (
	"gorm.io/gorm"

	"qwbot/model/task"
)

func InitTable(db *gorm.DB) error {
	// 自动迁移数据库表结构
	err := db.AutoMigrate(
		&task.Task{},
	)
	if err != nil {
		return err
	}
	return nil
}