package errors

type Error interface {
	Error() string
	Code() string
}

func New(code string, err string) error {
	return &errorString{code: code, err: err}
}

type errorString struct {
	code string
	err  string
}

func (e errorString) Error() string {
	return e.err
}

func (e errorString) Code() string {
	return e.code
}

var (
	InternalError = New("INTERNAL_ERROR", "Something went wrong. Please, try again.")
	InvalidMethod = New("INVALID_METHOD", "The method passed is invalid.")
)
