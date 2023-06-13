package errors

var (
	ErrorConflict = NewErrConflict("object")
	ErrorNotFound = NewErrNotFound("object")
)

// error not found
type ErrNotFound struct {
	text string
}

func NewErrNotFound(text string) *ErrNotFound {
	return &ErrNotFound{
		text: text,
	}
}

func (e *ErrNotFound) Error() string {
	return e.text + " not found"
}

// error conflict
type ErrConflict struct {
	text string
}

func NewErrConflict(text string) *ErrConflict {
	return &ErrConflict{
		text: text,
	}
}

func (e *ErrConflict) Error() string {
	return e.text + " already exist"
}
