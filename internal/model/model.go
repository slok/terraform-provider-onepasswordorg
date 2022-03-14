package model

// User represents a 1password user.
type User struct {
	ID    string
	Email string
	Name  string
}

// Group represents a 1password group.
type Group struct {
	ID          string
	Name        string
	Description string
}

// GroupRole represents a 1password user membership role.
type GroupRole int

const (
	GroupRoleUnknown GroupRole = iota
	GroupRoleMember
	GroupRoleAdmin
)

// GroupRole represents a 1password user membership into a group.
type GroupMembership struct {
	UserID  string
	GroupID string
	Role    GroupRole
}
