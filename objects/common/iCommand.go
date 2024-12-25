package common

type Command interface {
	GetCommand() string
	GetAlias() []string
	Execute(args ...string) error
}
