package events

type Logger interface {
	Log(scope Scope, op Operation, err error)
	LogFatal(scope Scope, op Operation, err error)
}
