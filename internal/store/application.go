package store

import (
	"context"
	"encoding/base64"
	"strconv"

	"github.com/dcm-project/dcm-placement-api/internal/store/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	// defaultPageSize defines the maximum number of items that can be returned in a single page
	defaultPageSize = 100
)

type Application interface {
	List(ctx context.Context, pageSize *int, pageToken *string) (model.ApplicationList, *string, error)
	Create(ctx context.Context, app model.Application) (*model.Application, error)
	Update(ctx context.Context, app model.Application) (*model.Application, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Get(ctx context.Context, id uuid.UUID) (*model.Application, error)
}

type ApplicationStore struct {
	db *gorm.DB
}

var _ Application = (*ApplicationStore)(nil)

func NewApplication(db *gorm.DB) Application {
	return &ApplicationStore{db: db}
}

func (s *ApplicationStore) List(ctx context.Context, pageSize *int, pageToken *string) (model.ApplicationList, *string, error) {
	var apps model.ApplicationList

	// Default page size
	limit := defaultPageSize
	if pageSize != nil {
		limit = *pageSize
	}

	// Parse page token to get offset
	offset := 0
	if pageToken != nil && *pageToken != "" {
		decoded, err := base64.StdEncoding.DecodeString(*pageToken)
		if err == nil {
			if parsedOffset, err := strconv.Atoi(string(decoded)); err == nil {
				offset = parsedOffset
			}
		}
	}

	// Query with limit and offset
	tx := s.db.Model(&apps).Limit(limit + 1).Offset(offset)
	result := tx.Find(&apps)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	// Check if there are more results
	var nextPageToken *string
	if len(apps) > limit {
		apps = apps[:limit]
		nextOffset := offset + limit
		token := base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(nextOffset)))
		nextPageToken = &token
	}

	return apps, nextPageToken, nil
}

func (s *ApplicationStore) Delete(ctx context.Context, id uuid.UUID) error {
	result := s.db.Delete(&model.Application{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *ApplicationStore) Create(ctx context.Context, app model.Application) (*model.Application, error) {
	result := s.db.Clauses(clause.Returning{}).Create(&app)
	if result.Error != nil {
		return nil, result.Error
	}

	return &app, nil
}

func (s *ApplicationStore) Update(ctx context.Context, app model.Application) (*model.Application, error) {
	result := s.db.Save(&app)
	if result.Error != nil {
		return nil, result.Error
	}

	return &app, nil
}

func (s *ApplicationStore) Get(ctx context.Context, id uuid.UUID) (*model.Application, error) {
	var app model.Application
	result := s.db.First(&app, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &app, nil
}
