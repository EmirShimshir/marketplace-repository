package entity

import (
	"github.com/EmirShimshir/marketplace-core/domain"
	"github.com/google/uuid"
	"github.com/guregu/null"
)

const (
	MgUserCustomer  = "Customer"
	MgUserSeller    = "Seller"
	MgUserModerator = "Moderator"
)

type MgUser struct {
	ID       string   `bson:"_id"`
	CartID   string   `bson:"cart_id"`
	Name     string      `bson:"name"`
	Surname  string      `bson:"surname"`
	Phone    null.String `bson:"phone,omitempty"`
	Email    string      `bson:"email"`
	Password string      `bson:"password"`
	Role     string      `bson:"role"`
}

func (u *MgUser) ToDomain() domain.User {
	var userRole domain.UserRole
	switch u.Role {
	case MgUserCustomer:
		userRole = domain.UserCustomer
	case MgUserSeller:
		userRole = domain.UserSeller
	case MgUserModerator:
		userRole = domain.UserModerator
	}
	return domain.User{
		ID:       domain.ID(u.ID),
		CartID:   domain.ID(u.CartID),
		Name:     u.Name,
		Surname:  u.Surname,
		Phone:    u.Phone,
		Email:    u.Email,
		Password: u.Password,
		Role:     userRole,
	}
}

func NewMgUser(user domain.User) MgUser {
	id, _ := uuid.Parse(user.ID.String())
	cartID, _ := uuid.Parse(user.CartID.String())
	var userRole string
	switch user.Role {
	case domain.UserCustomer:
		userRole = MgUserCustomer
	case domain.UserSeller:
		userRole = MgUserSeller
	case domain.UserModerator:
		userRole = MgUserModerator
	}
	return MgUser{
		ID:       id.String(),
		CartID:   cartID.String(),
		Name:     user.Name,
		Surname:  user.Surname,
		Phone:    user.Phone,
		Email:    user.Email,
		Password: user.Password,
		Role:     userRole,
	}
}
