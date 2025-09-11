package store

import (
	"gorm.io/gorm"
)

type Store interface {
	Close() error
	DeclaredVm() DeclaredVm
	RequestedVm() RequestedVm
}

type DataStore struct {
	db          *gorm.DB
	declaredVm  DeclaredVm
	requestedVm RequestedVm
}

func NewStore(db *gorm.DB) Store {
	return &DataStore{
		db:          db,
		declaredVm:  NewDeclaredVm(db),
		requestedVm: NewRequestedVm(db),
	}
}

func (s *DataStore) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (s *DataStore) DeclaredVm() DeclaredVm {
	return s.declaredVm
}

func (s *DataStore) RequestedVm() RequestedVm {
	return s.requestedVm
}
