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

// Vault represents a 1password vault.
type Vault struct {
	ID          string
	Name        string
	Description string
}

// MembershipRole represents a 1password user membership role.
type MembershipRole int

const (
	MembershipRoleMember MembershipRole = iota
	MembershipRoleManager
)

// Role represents a 1password user membership into a group.
type Membership struct {
	UserID  string
	GroupID string
	Role    MembershipRole
}
