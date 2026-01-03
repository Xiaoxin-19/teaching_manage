package repository

import (
	"context"
	"teaching_manage/backend/dao"
	"teaching_manage/backend/entity"
	"teaching_manage/backend/model"
	"time"
)

type RecordRepository interface {
	CreateRecord(ctx context.Context, record *entity.Record) error
	GetRecordList(ctx context.Context, stuIds []uint, teachIDs []uint, subjectIDs []uint,
		startDate string, endDate string, offset int, limit int, active *bool) ([]entity.Record, int64, int64, error)
	GetAllPendingRecordList(ctx context.Context) ([]entity.Record, error)
	ActivateRecord(context.Context, uint) error
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
	recordModel := model.Record{
		StudentID:    record.Student.ID,
		TeacherID:    record.Teacher.ID,
		SubjectID:    record.Subject.ID,
		TeachingDate: record.TeachingDate,
		StartTime:    record.StartTime.Format("15:04"),
		EndTime:      record.EndTime.Format("15:04"),
		Remark:       record.Remark,
	}
	return r.recordDao.CreateRecord(ctx, recordModel)
}

func (r *RecordRepositoryImpl) GetRecordList(ctx context.Context, stuIDs []uint, teachIDs []uint, subjectIDs []uint,
	startDate string, endDate string, offset int, limit int, active *bool) ([]entity.Record, int64, int64, error) {
	records, total, pendingTotal, err := r.recordDao.GetRecordList(ctx, stuIDs, teachIDs, subjectIDs, startDate, endDate, offset, limit, active)
	if err != nil {
		return nil, 0, 0, err
	}
	var result []entity.Record
	for _, rec := range records {
		startTime, err := time.Parse("15:04", rec.StartTime)
		if err != nil {
			return nil, 0, 0, err
		}
		endTime, err := time.Parse("15:04", rec.EndTime)
		if err != nil {
			return nil, 0, 0, err
		}
		result = append(result, entity.Record{
			ID:           rec.ID,
			CreatedAt:    rec.CreatedAt,
			UpdatedAt:    rec.UpdatedAt,
			Student:      entity.Student{ID: rec.StudentID, Name: rec.Student.Name, StudentNumber: rec.Student.StudentNumber, Status: int(rec.Student.Status)},
			Teacher:      entity.Teacher{ID: rec.TeacherID, Name: rec.Teacher.Name, TeacherNumber: rec.Teacher.TeacherNumber},
			Subject:      entity.Subject{ID: rec.SubjectID, Name: rec.Subject.Name},
			TeachingDate: rec.TeachingDate,
			StartTime:    startTime,
			EndTime:      endTime,
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

	startTime, err := time.Parse("15:04", dbRecord.StartTime)
	if err != nil {
		return entity.Record{}, err
	}
	endTime, err := time.Parse("15:04", dbRecord.EndTime)
	if err != nil {
		return entity.Record{}, err
	}

	return entity.Record{
		ID:           dbRecord.ID,
		CreatedAt:    dbRecord.CreatedAt,
		UpdatedAt:    dbRecord.UpdatedAt,
		Student:      entity.Student{ID: dbRecord.StudentID, Name: dbRecord.Student.Name, StudentNumber: dbRecord.Student.StudentNumber, Status: int(dbRecord.Student.Status)},
		Teacher:      entity.Teacher{ID: dbRecord.TeacherID, Name: dbRecord.Teacher.Name, TeacherNumber: dbRecord.Teacher.TeacherNumber},
		Subject:      entity.Subject{ID: dbRecord.SubjectID, Name: dbRecord.Subject.Name},
		TeachingDate: dbRecord.TeachingDate,
		StartTime:    startTime,
		EndTime:      endTime,
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
		startTime, err := time.Parse("15:04", rec.StartTime)
		if err != nil {
			return nil, err
		}
		endTime, err := time.Parse("15:04", rec.EndTime)
		if err != nil {
			return nil, err
		}
		result = append(result, entity.Record{
			ID:           rec.ID,
			CreatedAt:    rec.CreatedAt,
			UpdatedAt:    rec.UpdatedAt,
			Student:      entity.Student{ID: rec.StudentID, Name: rec.Student.Name, StudentNumber: rec.Student.StudentNumber, Status: int(rec.Student.Status)},
			Teacher:      entity.Teacher{ID: rec.TeacherID, Name: rec.Teacher.Name, TeacherNumber: rec.Teacher.TeacherNumber},
			Subject:      entity.Subject{ID: rec.SubjectID, Name: rec.Subject.Name},
			TeachingDate: rec.TeachingDate,
			StartTime:    startTime,
			EndTime:      endTime,
			Active:       rec.Active,
			Remark:       rec.Remark,
		})
	}
	return result, nil
}
