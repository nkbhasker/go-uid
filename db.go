package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type DBStore interface {
	DB() *gorm.DB
	CloseDB() error
}

type dbStore struct {
	db *gorm.DB
}

func Init(postgresUrl string) (DBStore, error) {
	schema.RegisterSerializer("id", NewIdSerializer())

	db, err := gorm.Open(postgres.Open(postgresUrl), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto migrate model
	if err = db.AutoMigrate(&User{}); err != nil {
		return nil, err
	}

	return &dbStore{db: db}, nil
}

func (s dbStore) DB() *gorm.DB {
	return s.db
}

func (s *dbStore) CloseDB() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
