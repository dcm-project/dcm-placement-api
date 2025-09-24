package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Application struct {
	gorm.Model
	ID      uuid.UUID      `gorm:"primaryKey;"`
	Name    string         `gorm:"name;not null"`
	Service string         `gorm:"service;not null"`
	Zones   pq.StringArray `gorm:"type:text[]"`
	Tier    string         `gorm:"tier;not null"`
}

type ApplicationList []Application

type RequestedVm struct {
	gorm.Model
	ID       uuid.UUID `gorm:"primaryKey;"`
	Name     string    `gorm:"name;not null"`
	Env      string    `gorm:"not null"`
	Ram      int       `gorm:"not null"`
	Os       string    `gorm:"not null"`
	Cpu      int       `gorm:"not null"`
	Region   string    `gorm:"not null"`
	Role     string    `gorm:"not null"`
	TenantId string    `gorm:"not null"`
}

type RequestedVmList []RequestedVm

type DeclaredVm struct {
	gorm.Model
	ID            uuid.UUID   `gorm:"primaryKey;"`
	RequestedVmID uuid.UUID   `gorm:"type:uuid;not null"`       // Foreign key to RequestedVm
	RequestedVm   RequestedVm `gorm:"foreignKey:RequestedVmID"` // Relationship to RequestedVm
	IPAddress     string      `gorm:"not null"`
	Gateway       string      `gorm:"not null"`
	Netmask       string      `gorm:"not null"`
	DnsName       string      `gorm:"not null"`
	CreatedAt     time.Time   `gorm:"not null"`
}

type DeclaredVmList []DeclaredVm
