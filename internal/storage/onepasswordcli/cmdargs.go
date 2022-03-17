package onepasswordcli

import "strings"

type onePasswordCliCmd struct {
	args []string
}

func (o *onePasswordCliCmd) GetArgs() []string {
	return o.args
}

func (o *onePasswordCliCmd) CreateArg() *onePasswordCliCmd {
	o.args = append(o.args, "create")
	return o
}

func (o *onePasswordCliCmd) EditArg() *onePasswordCliCmd {
	o.args = append(o.args, "edit")
	return o
}

func (o *onePasswordCliCmd) GetArg() *onePasswordCliCmd {
	o.args = append(o.args, "get")
	return o
}

func (o *onePasswordCliCmd) ListArg() *onePasswordCliCmd {
	o.args = append(o.args, "list")
	return o
}

func (o *onePasswordCliCmd) ProvisionArg() *onePasswordCliCmd {
	o.args = append(o.args, "provision")
	return o
}

func (o *onePasswordCliCmd) DeleteArg() *onePasswordCliCmd {
	o.args = append(o.args, "delete")
	return o
}

func (o *onePasswordCliCmd) GrantArg() *onePasswordCliCmd {
	o.args = append(o.args, "grant")
	return o
}

func (o *onePasswordCliCmd) RevokeArg() *onePasswordCliCmd {
	o.args = append(o.args, "revoke")
	return o
}

func (o *onePasswordCliCmd) UserArg() *onePasswordCliCmd {
	o.args = append(o.args, "user")
	return o
}

func (o *onePasswordCliCmd) GroupArg() *onePasswordCliCmd {
	o.args = append(o.args, "group")
	return o
}

func (o *onePasswordCliCmd) VaultArg() *onePasswordCliCmd {
	o.args = append(o.args, "vault")
	return o
}

func (o *onePasswordCliCmd) RawStrArg(s string) *onePasswordCliCmd {
	o.args = append(o.args, s)
	return o
}

func (o *onePasswordCliCmd) NameFlag(name string) *onePasswordCliCmd {
	o.args = append(o.args, "--name", name)
	return o
}

func (o *onePasswordCliCmd) DescriptionFlag(description string) *onePasswordCliCmd {
	if description == "" {
		return o
	}

	o.args = append(o.args, "--description", description)
	return o
}

func (o *onePasswordCliCmd) RoleFlag(role string) *onePasswordCliCmd {
	if role == "" {
		return o
	}
	o.args = append(o.args, "--role", role)
	return o
}

func (o *onePasswordCliCmd) GroupFlag(id string) *onePasswordCliCmd {
	o.args = append(o.args, "--group", id)
	return o
}

func (o *onePasswordCliCmd) UserFlag(id string) *onePasswordCliCmd {
	o.args = append(o.args, "--user", id)
	return o
}

func (o *onePasswordCliCmd) FormatJSONFlag() *onePasswordCliCmd {
	o.args = append(o.args, "--format", "json")
	return o
}

func (o *onePasswordCliCmd) EmailFlag(email string) *onePasswordCliCmd {
	o.args = append(o.args, "--email", email)
	return o
}

func (o *onePasswordCliCmd) VaultFlag(id string) *onePasswordCliCmd {
	o.args = append(o.args, "--vault", id)
	return o
}

func (o *onePasswordCliCmd) NoInputFlag() *onePasswordCliCmd {
	o.args = append(o.args, "--no-input")
	return o
}

func (o *onePasswordCliCmd) PermissionsFlag(permissions []string) *onePasswordCliCmd {
	if len(permissions) == 0 {
		return o
	}

	p := strings.Join(permissions, ",")
	o.args = append(o.args, "--permissions", p)

	return o
}
