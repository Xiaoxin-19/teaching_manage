package dao

/*
CREATE TABLE `teachers` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(10) NOT NULL COMMENT '教师姓名',
  `gender` char(3) DEFAULT NULL COMMENT '教师性别',
  `phone_code` varchar(11) DEFAULT NULL COMMENT '电话号码',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `delete_time` datetime DEFAULT NULL COMMENT '软删除标记时间戳',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `name` (`name`) USING BTREE COMMENT '教师姓名索引'
) ENGINE=InnoDB AUTO_INCREMENT=20 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;
*/

import (
	"context"

	"gorm.io/gorm"
)

type TeacherDao interface {
	CreateTeacher(ctx context.Context, t *Teacher) error
	UpdateTeacher(ctx context.Context, t *Teacher) error
	DeleteTeacher(ctx context.Context, id uint) error
	GetTeacherByID(ctx context.Context, id uint) (*Teacher, error)
	GetTeacherList(ctx context.Context, key string, offset int, limit int) ([]Teacher, int64, error)
}

type TeacherGormDao struct {
	db *gorm.DB
}

func NewTeacherDao(db *gorm.DB) TeacherDao {
	return &TeacherGormDao{db: db}
}

type Teacher struct {
	gorm.Model
	Name   string `gorm:"column:name;not null;comment:教师姓名;index;unique" json:"name"`
	Gender string `gorm:"column:gender;comment:教师性别" json:"gender"`
	Phone  string `gorm:"column:phone;comment:电话号码" json:"phone"`
	Remark string `gorm:"column:remark;comment:备注" json:"remark"`
}

func (s TeacherGormDao) CreateTeacher(ctx context.Context, t *Teacher) error {
	return gorm.G[Teacher](s.db).Create(ctx, &Teacher{
		Name:   t.Name,
		Gender: t.Gender,
		Phone:  t.Phone,
		Remark: t.Remark,
	})

}

func (s TeacherGormDao) UpdateTeacher(ctx context.Context, t *Teacher) error {
	_, err := gorm.G[Teacher](s.db).Where("id = ?", t.ID).Select("name", "gender", "phone", "remark").Updates(ctx, Teacher{
		Name:   t.Name,
		Gender: t.Gender,
		Phone:  t.Phone,
		Remark: t.Remark,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s TeacherGormDao) DeleteTeacher(ctx context.Context, id uint) error {
	_, err := gorm.G[Teacher](s.db).Where("id = ?", id).Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s TeacherGormDao) GetTeacherByID(ctx context.Context, id uint) (*Teacher, error) {
	t, err := gorm.G[Teacher](s.db).Where("id = ?", id).First(ctx)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// Get teacher list
func (s TeacherGormDao) GetTeacherList(ctx context.Context, key string, offset int, limit int) ([]Teacher, int64, error) {
	var teachers []Teacher
	query := gorm.G[Teacher](s.db).Where("")

	if key != "" {
		query = query.Where("name LIKE ?", "%"+key+"%")
	}
	total, err := query.Count(ctx, "*")
	if err != nil {
		return nil, 0, err
	}

	// 处理没有分页参数的情况
	if limit == 0 {
		offset = 0
		limit = int(total)
	}

	teachers, err = query.Offset(offset).Limit(limit).Find(ctx)
	return teachers, total, nil
}
