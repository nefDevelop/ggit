package testx

// InitNoopLogger is a init function to initialize a noop logger.
// It is now a dummy function because the log package handles "noop" internally.
var InitNoopLogger = func(name string, vs ...any) (any, error) {
	return nil, nil
}
