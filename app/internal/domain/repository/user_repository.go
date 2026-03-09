package repository

import (
	"context"
	"gorm.io/gorm"
	"task-manager/internal/domain/models"
)

type UserRepository interface {
	CreateUser(context.Context, *models.User) error
	FindUserByID(context.Context, string) (*models.User, error)
	FindUserByEmail(context.Context, string) (*models.User, error)
	FindAll(context.Context) ([]models.User, error)
	Update(context.Context, *models.User) (*models.User, error)
	Delete(context.Context, string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) FindUserByID(ctx context.Context, userID string) (*models.User, error) {
	user := &models.User{}
	err := r.db.WithContext(ctx).First(user, "id = ?", userID).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	err := r.db.WithContext(ctx).First(user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) FindAll(ctx context.Context) ([]models.User, error) {
	var users []models.User
	err := r.db.WithContext(ctx).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) (*models.User, error) {
	err := r.db.WithContext(ctx).Model(user).Updates(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) Delete(ctx context.Context, userID string) error {
	user := &models.User{}
	err := r.db.WithContext(ctx).First(user, "id = ?", userID).Error
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Delete(user).Error
}
