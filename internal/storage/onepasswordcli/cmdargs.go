package onepasswordcli

type onePasswordCliCmd struct {
	args []string
}

func (o *onePasswordCliCmd) GetArgs() []string {
	return o.args
}

func (o *onePasswordCliCmd) AddArg() *onePasswordCliCmd {
	o.args = append(o.args, "add")
	return o
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

func (o *onePasswordCliCmd) DeleteArg() *onePasswordCliCmd {
	o.args = append(o.args, "delete")
	return o
}

func (o *onePasswordCliCmd) RemoveArg() *onePasswordCliCmd {
	o.args = append(o.args, "remove")
	return o
}

func (o *onePasswordCliCmd) UsersArg() *onePasswordCliCmd {
	o.args = append(o.args, "users")
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
