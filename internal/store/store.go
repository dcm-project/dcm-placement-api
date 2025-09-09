package store

import (
	"gorm.io/gorm"
)

type Store interface {
	Close() error
}

type DataStore struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) Store {
	return &DataStore{
		db: db,
	}
}

func (s *DataStore) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
