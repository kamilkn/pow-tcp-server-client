package cache

type Logger interface {
	Debug(msg string, args ...any)
}
