package e

func WithoutUserError() error {
	err := New("Database can't find this user")
	return err
}

func UserBeDeleted() error {
	err := New("User be deleted")
	return err
}

func UserBeStoped() error {
	err := New("User be stoped")
	return err
}

// New returns an error that formats as the given text.
func New(text string) error {
	return &errorString{text}
}

// errorString is a trivial implementation of error.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}
