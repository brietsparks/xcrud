package cli

type Logger interface {
	Error(args ...interface{})
}
