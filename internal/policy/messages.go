package policy

const (
	ErrMissingTarget      = "target is required"
	ErrMissingInstrument  = "instrument is required"
	ErrMissingDescription = "description is required"
	ErrNotFound           = "observation not found"
	ErrImmutable          = "released observation cannot be changed"
	ErrInvalidDecision    = "decision must be approve or return"
	ErrInvalidTransition  = "observation is not ready for this transition"
)
