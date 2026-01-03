package repository

import (
	"context"
	"teaching_manage/backend/dao"
	"teaching_manage/backend/entity"
	"teaching_manage/backend/model"

	"gorm.io/gorm"
)

type SubjectRepository interface {
	GetSubjectList(ctx context.Context, key string, offset int, limit int) ([]entity.Subject, int64, error)
	CreateSubject(ctx context.Context, subject entity.Subject) error
	UpdateSubject(ctx context.Context, subject entity.Subject) error
	DeleteSubject(ctx context.Context, id uint) error
	GetAllSubjects(ctx context.Context) ([]entity.Subject, error)
}

func (s SubjectRepositoryImpl) GetAllSubjects(ctx context.Context) ([]entity.Subject, error) {
	subjects, _, err := s.dao.GetSubjectList(ctx, "", 0, -1)
	if err != nil {
		return nil, err
	}
	var entities []entity.Subject
	for _, s := range subjects {
		entities = append(entities, entity.Subject{
			ID:            s.ID,
			SubjectNumber: s.SubjectNumber,
			Name:          s.Name,
			CreatedAt:     s.CreatedAt,
			UpdatedAt:     s.UpdatedAt,
		})
	}
	return entities, nil
}

type SubjectRepositoryImpl struct {
	dao dao.SubjectDao
}

func NewSubjectRepository(dao dao.SubjectDao) SubjectRepository {
	return &SubjectRepositoryImpl{dao: dao}
}

func (sr SubjectRepositoryImpl) GetSubjectList(ctx context.Context, key string, offset int, limit int) ([]entity.Subject, int64, error) {
	subjects, total, err := sr.dao.GetSubjectList(ctx, key, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	var entities []entity.Subject
	for _, s := range subjects {
		entities = append(entities, entity.Subject{
			ID:            s.ID,
			SubjectNumber: s.SubjectNumber,
			Name:          s.Name,
			CreatedAt:     s.CreatedAt,
			UpdatedAt:     s.UpdatedAt,
			StudentSubjects: func() []entity.StudentSubject {
				var studentSubjects []entity.StudentSubject
				for _, ss := range s.StudentSubjects {
					studentSubjects = append(studentSubjects, entity.StudentSubject{
						ID:      ss.ID,
						Teacher: entity.Teacher{ID: ss.ID},
						Student: entity.Student{ID: ss.StudentID},
						Balance: ss.Balance,
					})
				}
				return studentSubjects
			}(),
		})
	}
	return entities, total, nil
}

func (sr SubjectRepositoryImpl) CreateSubject(ctx context.Context, subject entity.Subject) error {
	err := sr.dao.CreateSubject(ctx, &model.Subject{
		Name: subject.Name,
	})
	return err
}

func (sr SubjectRepositoryImpl) UpdateSubject(ctx context.Context, subject entity.Subject) error {
	err := sr.dao.UpdateSubject(ctx, &model.Subject{
		Model: gorm.Model{ID: subject.ID},
		Name:  subject.Name,
	})
	return err
}
func (sr SubjectRepositoryImpl) DeleteSubject(ctx context.Context, id uint) error {
	err := sr.dao.DeleteSubject(ctx, id)
	return err
}
