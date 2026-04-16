package converter

import (
	"errors"

	"github.com/mmms914/test-avito-internship/internal/api/models"
	"github.com/mmms914/test-avito-internship/internal/domain"
)

func UserRoleFromString(strRole string) (domain.Role, error) {
	switch strRole {
	case domain.AdminRole.String():
		return domain.AdminRole, nil
	case domain.UserRole.String():
		return domain.UserRole, nil
	default:
		return "", errors.New("invalid role")
	}
}

func UserToResponse(u *domain.User) *models.UserResponse {
	return &models.UserResponse{
		User: UserToObject(u),
	}
}

func UserToObject(u *domain.User) *models.UserObject {
	return &models.UserObject{
		ID:        u.ID(),
		Email:     u.Email(),
		Role:      u.Role().String(),
		CreatedAt: u.CreatedAt(),
	}
}
