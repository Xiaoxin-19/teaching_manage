package dao

import (
	"context"

	"gorm.io/gorm"
)

type UnitWorkDao interface {
	Execute(context.Context, func(context.Context) error) error
}

type UnitWorkGormDao struct {
	db DBGetter
}

type txKey struct{}

func NewUnitWorkDao(db DBGetter) *UnitWorkGormDao {
	return &UnitWorkGormDao{db: db}
}

func (u *UnitWorkGormDao) Execute(ctx context.Context, fn func(context.Context) error) error {
	db := u.db()
	return db.Transaction(func(tx *gorm.DB) error {
		return fn(context.WithValue(ctx, txKey{}, tx))
	})
}

// GetDBFromCtx 是一个公共辅助函数，供所有 DAO 使用，用于从 context 中提取事务
func GetDBFromCtx(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}
	return defaultDB
}
