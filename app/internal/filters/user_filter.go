package filters

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"strconv"
	"task-manager/internal/auth"
	"task-manager/internal/domain/models"
)

type UserFilter struct {
	Role     *auth.UserRole
	Roles    []*auth.UserRole
	IsActive *bool
}

func ApplyUserFilter(ctx *gin.Context) (models.UserFilter, error) {
	filter := &UserFilter{}

	roleParam := ctx.Query("role")
	if roleParam != "" {
		role := auth.UserRole(roleParam)
		if !role.IsValid() {
			return filter, errors.New("invalid role filter")
		}
		filter.Role = &role
	}

	rolesParams := ctx.QueryArray("roles[]")
	log.Println("rolesParams: ", rolesParams)
	if len(rolesParams) > 0 {
		for _, r := range rolesParams {
			role := auth.UserRole(r)
			if !role.IsValid() {
				return filter, errors.New("invalid roles filter")
			}
			filter.Roles = append(filter.Roles, &role)
		}
	}

	isActiveParam := ctx.Query("is_active")
	if isActiveParam != "" {
		isActive, err := strconv.ParseBool(isActiveParam)
		if err != nil {
			return filter, errors.New("invalid is_active filter")
		}
		filter.IsActive = &isActive
	}

	return filter, nil
}

func (filter *UserFilter) Apply(db *gorm.DB) (*gorm.DB, error) {
	if filter.Role != nil {
		db = db.Where("role = ?", *filter.Role)
	}
	log.Println("filter.Roles: ", filter.Roles)
	if len(filter.Roles) > 0 {
		db = db.Where("role IN ?", filter.Roles)
	}

	if filter.IsActive != nil {
		db = db.Where("is_active = ?", *filter.IsActive)
	}
	log.Println(db.Statement.SQL.String())
	return db, nil
}
