package service

import (
	"context"
	"log"
	"task-manager/internal/auth"
	"task-manager/internal/domain/models"
	"task-manager/internal/domain/repository"
	"task-manager/internal/myerrors"
	"task-manager/internal/utils"
	"time"
)

type UserService interface {
	Register(ctx context.Context, name string, email string, password string, role *auth.UserRole) (*models.User, error)
	GetAll(ctx context.Context) ([]models.User, error)
	GetById(ctx context.Context, id string) (*models.User, error)
	Update(ctx context.Context, ID string, name *string, email *string, password *string, role *auth.UserRole, isActive *bool) (*models.User, error)
	Delete(ctx context.Context, id string) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return userService{repo: repo}
}

func (s userService) Register(ctx context.Context, name string, email string, password string, role *auth.UserRole) (*models.User, error) {
	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	if role == nil {
		role = new(auth.UserRole)
	}
	if !role.IsValid() {
		*role = auth.UserRoleViewer
	}

	newUserId := utils.NewUUID()

	user := &models.User{
		ID:       newUserId,
		Name:     name,
		Email:    email,
		Password: passwordHash,
		Role:     *role,
		IsActive: true,
	}

	err = s.repo.CreateUser(ctx, user)
	if err != nil {
		return user, myerrors.EntityAlreadyExists("user")
	}

	user, err = s.repo.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s userService) GetAll(ctx context.Context) ([]models.User, error) {
	users, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s userService) GetById(ctx context.Context, id string) (*models.User, error) {
	user, err := s.repo.FindUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s userService) Update(ctx context.Context, ID string, name *string, email *string, password *string, role *auth.UserRole, isActive *bool) (*models.User, error) {
	user, err := s.repo.FindUserByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	var changed = false

	if name != nil && user.Name != *name {
		user.Name = *name
		changed = true
	}

	if email != nil && user.Email != *email {
		user.Email = *email
		changed = true
	}

	if role != nil && role.IsValid() && user.Role != *role {
		user.Role = *role
		changed = true
	}

	if isActive != nil && user.IsActive != *isActive {
		user.IsActive = *isActive
		changed = true
	}

	if password != nil {
		err := utils.CheckPasswordHash(*password, user.Password)
		if err == nil {
			newPassword, err := utils.HashPassword(*password)
			if err != nil {
				return nil, err
			}
			user.Password = newPassword
			changed = true
		}
	}

	if changed {
		user.UpdatedAt = time.Now()
		user, err = s.repo.Update(ctx, user)
		if err != nil {
			return nil, err
		}
	}

	log.Println(user)

	return user, nil
}

func (s userService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
