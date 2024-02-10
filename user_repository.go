package main

import "gorm.io/gorm"

type UserRepository interface {
	NewUser(name string, email string) (*User, error)
	CreateUser(user *User) error
	GetUser(id Identifier) (*User, error)
}

type userRepository struct {
	IdGenerator IdGenerator
	DB          *gorm.DB
}

func NewUserRepository(db *gorm.DB, idGenerator IdGenerator) UserRepository {
	return &userRepository{
		DB:          db,
		IdGenerator: idGenerator,
	}
}

func (r *userRepository) NewUser(name string, email string) (*User, error) {
	id, err := r.IdGenerator.Next("user")
	if err != nil {
		return nil, err
	}

	return &User{
		ID:    id,
		Name:  name,
		Email: email,
	}, nil
}

func (r *userRepository) CreateUser(user *User) error {
	return r.DB.Create(user).Error
}

func (r *userRepository) GetUser(id Identifier) (*User, error) {
	var user User
	tx := r.DB.Find(&user, id)
	if tx.RowsAffected > 0 {
		return &user, nil
	}
	if tx.Error != nil {
		return nil, tx.Error
	}

	return nil, nil
}
