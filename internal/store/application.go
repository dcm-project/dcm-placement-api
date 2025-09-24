package store

import (
	"context"

	"github.com/dcm-project/dcm-placement-api/internal/store/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Application interface {
	List(ctx context.Context) (model.ApplicationList, error)
	Create(ctx context.Context, app model.Application) (*model.Application, error)
}

type ApplicationStore struct {
	db *gorm.DB
}

var _ Application = (*ApplicationStore)(nil)

func NewApplication(db *gorm.DB) Application {
	return &ApplicationStore{db: db}
}

func (s *ApplicationStore) List(ctx context.Context) (model.ApplicationList, error) {
	var apps model.ApplicationList
	tx := s.db.Model(&apps)
	result := tx.Find(&apps)
	if result.Error != nil {
		return nil, result.Error
	}
	return apps, nil
}

func (s *ApplicationStore) Create(ctx context.Context, app model.Application) (*model.Application, error) {
	result := s.db.Clauses(clause.Returning{}).Create(&app)
	if result.Error != nil {
		return nil, result.Error
	}

	return &app, nil
}
