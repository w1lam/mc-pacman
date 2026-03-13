package events

type Logger interface {
	Log(Event)
	Close() error
}
