package onepasswordcli

type onePasswordCliCmd struct {
	args []string
}

func (o *onePasswordCliCmd) GetArgs() []string {
	return o.args
}

func (o *onePasswordCliCmd) WithAdd() *onePasswordCliCmd {
	o.args = append(o.args, "add")
	return o
}

func (o *onePasswordCliCmd) WithCreate() *onePasswordCliCmd {
	o.args = append(o.args, "create")
	return o
}

func (o *onePasswordCliCmd) WithEdit() *onePasswordCliCmd {
	o.args = append(o.args, "edit")
	return o
}

func (o *onePasswordCliCmd) WithGet() *onePasswordCliCmd {
	o.args = append(o.args, "get")
	return o
}

func (o *onePasswordCliCmd) WithDelete() *onePasswordCliCmd {
	o.args = append(o.args, "delete")
	return o
}

func (o *onePasswordCliCmd) WithUserEmail(email string) *onePasswordCliCmd {
	o.args = append(o.args, "user", email)
	return o
}

func (o *onePasswordCliCmd) WithUserID(id string) *onePasswordCliCmd {
	o.args = append(o.args, "user", id)
	return o
}

func (o *onePasswordCliCmd) WithName(name string) *onePasswordCliCmd {
	o.args = append(o.args, name)
	return o
}

func (o *onePasswordCliCmd) WithEmail(email string) *onePasswordCliCmd {
	o.args = append(o.args, email)
	return o
}

func (o *onePasswordCliCmd) WithNewName(name string) *onePasswordCliCmd {
	o.args = append(o.args, "--name", name)
	return o
}

func (o *onePasswordCliCmd) WithRole(role string) *onePasswordCliCmd {
	o.args = append(o.args, "--role", role)
	return o
}
