package repository

import (
	"context"
	"teaching_manage/backend/dao"
)

type UnitWork interface {
	Execute(context.Context, func(context.Context) error) error
}

type UnitWorkImpl struct {
	dao dao.UnitWorkDao
}

func (u *UnitWorkImpl) Execute(ctx context.Context, fn func(context.Context) error) error {
	return u.dao.Execute(ctx, fn)
}

func NewUnitWork(dao dao.UnitWorkDao) UnitWork {
	return &UnitWorkImpl{dao: dao}
}
