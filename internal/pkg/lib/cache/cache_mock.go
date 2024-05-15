package cache

type mockLogger struct { //nolint:unused // mock
	cancelSignalHandled bool
}

func (l *mockLogger) Debug(msg string, _ ...any) { //nolint:unused // mock
	l.cancelSignalHandled = msg == "context canceled"
}
