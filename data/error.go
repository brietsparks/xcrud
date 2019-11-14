package data

type Error struct {
	Err error
	Msg string
}

func NewError(err error, msg string) *Error {
	return &Error{err, msg}
}

func (e Error) Error() string {
	return e.Msg
}

func (e Error) Unwrap() error {
	return e.Err
}
