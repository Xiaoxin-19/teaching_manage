package repository

import (
	"context"
	"teaching_manage/dao"
	"teaching_manage/entity"
)

type RecordRepository interface {
	CreateRecord(ctx context.Context, record *entity.Record) error
	GetRecordList(ctx context.Context, stuKey string, teachKey string,
		startDate string, endDate string, offset int, limit int) ([]entity.Record, int64, int64, error)
	GetAllPendingRecordList(ctx context.Context) ([]entity.Record, error)
	ActivateRecord(ctx context.Context, recordID uint) error
	GetRecordByID(ctx context.Context, d uint) (entity.Record, error)
	DeleteRecordByID(ctx context.Context, id uint) error
}

type RecordRepositoryImpl struct {
	recordDao dao.RecordDAO
}

func NewRecordRepository(dao dao.RecordDAO) RecordRepository {
	return &RecordRepositoryImpl{recordDao: dao}
}

func (r *RecordRepositoryImpl) CreateRecord(ctx context.Context, record *entity.Record) error {
	recordModel := dao.Record{
		StudentID:    record.Student.ID,
		TeacherID:    record.Teacher.ID,
		TeachingDate: record.TeachingDate,
		StartTime:    record.StartTime,
		EndTime:      record.EndTime,
		Remark:       record.Remark,
	}
	return r.recordDao.CreateRecord(ctx, recordModel)
}

func (r *RecordRepositoryImpl) GetRecordList(ctx context.Context, stuKey string, teachKey string,
	startDate string, endDate string, offset int, limit int) ([]entity.Record, int64, int64, error) {
	records, total, pendingTotal, err := r.recordDao.GetRecordList(ctx, stuKey, teachKey, startDate, endDate, offset, limit)
	if err != nil {
		return nil, 0, 0, err
	}
	var result []entity.Record
	for _, rec := range records {
		result = append(result, entity.Record{
			ID:           rec.ID,
			CreatedAt:    rec.CreatedAt,
			UpdatedAt:    rec.UpdatedAt,
			Student:      entity.Student{ID: rec.StudentID, Name: rec.Student.Name, DeletedAt: rec.Student.DeletedAt.Time},
			Teacher:      entity.Teacher{ID: rec.TeacherID, Name: rec.Teacher.Name, DeletedAt: rec.Teacher.DeletedAt.Time},
			TeachingDate: rec.TeachingDate,
			StartTime:    rec.StartTime,
			EndTime:      rec.EndTime,
			Active:       rec.Active,
			Remark:       rec.Remark,
		})
	}
	return result, total, pendingTotal, nil
}

func (r *RecordRepositoryImpl) GetRecordByID(ctx context.Context, d uint) (entity.Record, error) {
	dbRecord, err := r.recordDao.GetRecordByID(ctx, d)
	if err != nil {
		return entity.Record{}, err
	}

	return entity.Record{
		ID:           dbRecord.ID,
		CreatedAt:    dbRecord.CreatedAt,
		UpdatedAt:    dbRecord.UpdatedAt,
		Student:      entity.Student{ID: dbRecord.StudentID, Name: dbRecord.Student.Name, DeletedAt: dbRecord.Student.DeletedAt.Time},
		Teacher:      entity.Teacher{ID: dbRecord.TeacherID, Name: dbRecord.Teacher.Name, DeletedAt: dbRecord.Teacher.DeletedAt.Time},
		TeachingDate: dbRecord.TeachingDate,
		StartTime:    dbRecord.StartTime,
		EndTime:      dbRecord.EndTime,
		Active:       dbRecord.Active,
		Remark:       dbRecord.Remark,
	}, nil
}

func (r *RecordRepositoryImpl) ActivateRecord(ctx context.Context, recordID uint) error {
	return r.recordDao.ActivateRecord(ctx, recordID)
}

func (r *RecordRepositoryImpl) DeleteRecordByID(ctx context.Context, id uint) error {
	return r.recordDao.DeleteRecordByID(ctx, id)
}

func (r *RecordRepositoryImpl) GetAllPendingRecordList(ctx context.Context) ([]entity.Record, error) {
	dbRecords, err := r.recordDao.GetAllPendingRecordList(ctx)
	if err != nil {
		return nil, err
	}
	var result []entity.Record
	for _, rec := range dbRecords {
		result = append(result, entity.Record{
			ID:           rec.ID,
			CreatedAt:    rec.CreatedAt,
			UpdatedAt:    rec.UpdatedAt,
			Student:      entity.Student{ID: rec.StudentID, Name: rec.Student.Name, DeletedAt: rec.Student.DeletedAt.Time},
			Teacher:      entity.Teacher{ID: rec.TeacherID, Name: rec.Teacher.Name, DeletedAt: rec.Teacher.DeletedAt.Time},
			TeachingDate: rec.TeachingDate,
			StartTime:    rec.StartTime,
			EndTime:      rec.EndTime,
			Active:       rec.Active,
			Remark:       rec.Remark,
		})
	}
	return result, nil
}
