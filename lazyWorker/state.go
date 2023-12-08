package lazyWorker

//go:generate stringer -type State -linecomment
type State int32

const (
	_             State = iota // unknown
	StateCreated               // created
	StateRunning               // running
	StateStopping              // stopping
	StateStopped               // stopped
)
