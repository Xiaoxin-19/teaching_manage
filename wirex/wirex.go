package wirex

import (
	"teaching_manage/dao"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewGormDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("teaching_manage.db"), &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, err
	}
	// Auto migrate the Student and Teacher models
	err = db.AutoMigrate(&dao.Student{}, &dao.Teacher{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
