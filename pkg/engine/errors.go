package engine

// PlayRunErrorMessage represents a message that PlayRunError contains
type PlayRunErrorMessage string

const (
	// PlayFinished is a message provided when there's no more computing to do for the Play
	PlayFinished PlayRunErrorMessage = "play finished"

	// UnknownMessage is a message provided when unknown error happened
	UnknownMessage PlayRunErrorMessage = ""
)

// PlayRunError represents an error during the run
type PlayRunError struct {
	Message PlayRunErrorMessage
}

var _ error = &PlayRunError{}

func (e PlayRunError) Error() string {
	return string(e.Message)
}

// NewError creates a new engine error with provided message
func NewError(m PlayRunErrorMessage) PlayRunError {
	return PlayRunError{m}
}

// IsPlayEndedErorr checks if error indicates end of the play
func IsPlayEndedErorr(err error) bool {
	switch MessageForError(err) {
	case PlayFinished:
		return true
	}
	return false
}

// MessageForError returns the message for an error
func MessageForError(err error) PlayRunErrorMessage {
	switch t := err.(type) {
	case PlayRunError:
		return t.Message
	}
	return UnknownMessage
}
