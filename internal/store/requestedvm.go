package store

import (
	"context"

	"github.com/dcm-project/dcm-placement-api/internal/store/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RequestedVm interface {
	List(ctx context.Context) (model.RequestedVmList, error)
	Create(ctx context.Context, vm model.RequestedVm) (*model.RequestedVm, error)
}

type RequestedVmStore struct {
	db *gorm.DB
}

var _ RequestedVm = (*RequestedVmStore)(nil)

func NewRequestedVm(db *gorm.DB) RequestedVm {
	return &RequestedVmStore{db: db}
}

func (s *RequestedVmStore) List(ctx context.Context) (model.RequestedVmList, error) {
	var vms model.RequestedVmList
	tx := s.db.Model(&vms)
	result := tx.Find(&vms)
	if result.Error != nil {
		return nil, result.Error
	}
	return vms, nil
}

func (s *RequestedVmStore) Create(ctx context.Context, vm model.RequestedVm) (*model.RequestedVm, error) {
	result := s.db.Clauses(clause.Returning{}).Create(&vm)
	if result.Error != nil {
		return nil, result.Error
	}

	return &vm, nil
}
