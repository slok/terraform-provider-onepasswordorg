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

// ServiceAccountVaultPermission is one of the limited permissions that 1Password
// allows when granting a service account access to a vault at creation time.
// Valid values: "read_items", "write_items", "share_items". "write_items" and
// "share_items" both require "read_items".
type ServiceAccountVaultPermission string

const (
ServiceAccountVaultPermissionReadItems  ServiceAccountVaultPermission = "read_items"
ServiceAccountVaultPermissionWriteItems ServiceAccountVaultPermission = "write_items"
ServiceAccountVaultPermissionShareItems ServiceAccountVaultPermission = "share_items"
)

// ServiceAccount represents a 1password service account.
//
// VaultAccess maps a vault identifier (name or UUID, as accepted by
// `op service-account create --vault`) to the set of permissions the service
// account should have on that vault. The 1Password CLI only allows configuring
// vault access at service-account creation time; therefore any change to
// VaultAccess (including adding/removing a vault or changing its permissions)
// requires destroying and recreating the service account.
type ServiceAccount struct {
ID          string
Name        string
Token       string // Only available at creation time.
VaultAccess map[string][]ServiceAccountVaultPermission
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
