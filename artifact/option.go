package artifact

// Option is a functional option for configuring a [Store] at construction time.
type Option func(*storeConfig)

type storeConfig struct {
	// maxPerSession, when > 0, limits how many artifacts a single session may
	// hold. Attempts to Set beyond this limit return ErrLimitExceeded.
	maxPerSession int
}

// ErrLimitExceeded is returned when a per-session artifact limit is breached.
var ErrLimitExceeded = newLimitError()

type limitError struct{}

func (limitError) Error() string { return "artifact: per-session limit exceeded" }
func newLimitError() error       { return limitError{} }

// WithMaxPerSession returns an Option that limits each session to at most n
// artifacts. A value of 0 (the default) means no limit.
func WithMaxPerSession(n int) Option {
	return func(c *storeConfig) {
		c.maxPerSession = n
	}
}

// NewStoreWithOptions creates a Store that respects the supplied options.
// For stores that need no options prefer the simpler [NewStore].
func NewStoreWithOptions(opts ...Option) *Store {
	cfg := &storeConfig{}
	for _, o := range opts {
		o(cfg)
	}
	s := NewStore()
	s.cfg = cfg
	return s
}
