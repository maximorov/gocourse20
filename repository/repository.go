package repository

import (
	"context"
	"gorm.io/gorm"
	"solid/model"
)

// Liskov Substitution Principle
type UserRepository interface {
	Create(ctx context.Context, user model.User) (*model.User, error)
	Update(ctx context.Context, user model.User) error
}

type DBUserRepository struct {
	db *gorm.DB
}

func (r *DBUserRepository) Create(ctx context.Context, user model.User) (*model.User, error) {
	err := r.db.WithContext(ctx).Create(&user).Error
	return &user, err
}

func (r *DBUserRepository) Update(ctx context.Context, user model.User) error {
	return r.db.WithContext(ctx).Save(&user).Error
}

type MemoryUserRepository struct {
	users map[int]model.User
}

func (r *MemoryUserRepository) Create(_ context.Context, user model.User) (*model.User, error) {
	if r.users == nil {
		r.users = map[int]model.User{}
	}

	r.users[user.Id] = user

	return &user, nil
}

func (r *MemoryUserRepository) Update(_ context.Context, user model.User) error {
	if r.users == nil {
		r.users = map[int]model.User{}
	}
	r.users[user.Id] = user

	return nil
}
