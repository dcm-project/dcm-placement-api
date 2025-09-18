package store

import (
	"context"

	"github.com/dcm-project/dcm-placement-api/internal/store/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DeclaredVm interface {
	List(ctx context.Context) (model.DeclaredVmList, error)
	Create(ctx context.Context, vm model.DeclaredVm) (*model.DeclaredVm, error)
}

type DeclaredVmStore struct {
	db *gorm.DB
}

var _ DeclaredVm = (*DeclaredVmStore)(nil)

func NewDeclaredVm(db *gorm.DB) DeclaredVm {
	return &DeclaredVmStore{db: db}
}

func (s *DeclaredVmStore) List(ctx context.Context) (model.DeclaredVmList, error) {
	var vms model.DeclaredVmList
	tx := s.db.Model(&vms)
	// Preload the related RequestedVm data
	result := tx.Preload("RequestedVm").Find(&vms)
	if result.Error != nil {
		return nil, result.Error
	}
	return vms, nil
}

func (s *DeclaredVmStore) Create(ctx context.Context, vm model.DeclaredVm) (*model.DeclaredVm, error) {
	result := s.db.Clauses(clause.Returning{}).Create(&vm)
	if result.Error != nil {
		return nil, result.Error
	}

	return &vm, nil
}
