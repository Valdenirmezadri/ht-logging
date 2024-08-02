package logging

// Backend is the interface which a log backend need to implement to be able to
// be used as a logging backend.
type Backend interface {
	Log(Level, int, *Record) error
}
