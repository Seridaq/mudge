package command

type Command struct {
	name  string
	alias []string
}

func (com Command) GetCommand() string {
	return com.name
}

func (com Command) GetAlias() []string {
	return com.alias
}
