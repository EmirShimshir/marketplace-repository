package entity

import (
	"github.com/EmirShimshir/marketplace-domain/domain"
	"github.com/google/uuid"
	"github.com/guregu/null"
)

const (
	PgUserCustomer  = "Customer"
	PgUserSeller    = "Seller"
	PgUserModerator = "Moderator"
)

type PgUser struct {
	ID       uuid.UUID   `db:"id"`
	CartID   uuid.UUID   `db:"cart_id"`
	Name     string      `db:"name"`
	Surname  string      `db:"surname"`
	Phone    null.String `db:"phone"`
	Email    string      `db:"email"`
	Password string      `db:"password"`
	Role     string      `db:"role"`
}

func (u *PgUser) ToDomain() domain.User {
	var userRole domain.UserRole
	switch u.Role {
	case PgUserCustomer:
		userRole = domain.UserCustomer
	case PgUserSeller:
		userRole = domain.UserSeller
	case PgUserModerator:
		userRole = domain.UserModerator
	}
	return domain.User{
		ID:       domain.ID(u.ID.String()),
		CartID:   domain.ID(u.CartID.String()),
		Name:     u.Name,
		Surname:  u.Surname,
		Phone:    u.Phone,
		Email:    u.Email,
		Password: u.Password,
		Role:     userRole,
	}
}

func NewPgUser(user domain.User) PgUser {
	id, _ := uuid.Parse(user.ID.String())
	cartID, _ := uuid.Parse(user.CartID.String())
	var userRole string
	switch user.Role {
	case domain.UserCustomer:
		userRole = PgUserCustomer
	case domain.UserSeller:
		userRole = PgUserSeller
	case domain.UserModerator:
		userRole = PgUserModerator
	}
	return PgUser{
		ID:       id,
		CartID:   cartID,
		Name:     user.Name,
		Surname:  user.Surname,
		Phone:    user.Phone,
		Email:    user.Email,
		Password: user.Password,
		Role:     userRole,
	}
}
