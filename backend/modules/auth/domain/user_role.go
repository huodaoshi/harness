package domain

// UserRole is the platform account role stored in Mongo users.role and JWT claims.
type UserRole int

const (
	UserRoleStandard UserRole = 0
	UserRoleAdmin    UserRole = 2
)

func (r UserRole) Int() int { return int(r) }
