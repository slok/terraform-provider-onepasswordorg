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

type VaultGroupAccess struct {
	VaultID     string
	GroupID     string
	Permissions AccessPermissions
}

type VaultUserAccess struct {
	VaultID     string
	UserID      string
	Permissions AccessPermissions
}

// More information in https://developer.1password.com/docs/cli/vault-permissions.
type AccessPermissions struct {
	AllowViewing         bool
	AllowEditing         bool
	AllowManaging        bool
	ViewItems            bool
	CreateItems          bool
	EditItems            bool
	ArchiveItems         bool
	DeleteItems          bool
	ViewAndCopyPasswords bool
	ViewItemHistory      bool
	ImportItems          bool
	ExportItems          bool
	CopyAndShareItems    bool
	PrintItems           bool
	ManageVault          bool
}

// Item represents a 1password item.
type Item struct {
	ID       string
	Vault    Vault
	Title    string
	Fields   []Field
	Sections []Section
	URLs     []URL
	Tags     []string
	Category string
}

type Section struct {
	ID    string
	Label string
}

type URL struct {
	URL     string
	Primary bool
}

type Field struct {
	Section  *Section
	ID       string
	Type     string
	Purpose  string
	Label    string
	Value    string
	Generate bool
}
