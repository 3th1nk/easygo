package ldap

type Error struct {
	Type ErrType
	Err  error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Type.String()
}

func isError(err error, typ ErrType) bool {
	if v, _ := err.(*Error); v != nil {
		return v.Type == typ
	}
	return false
}

//go:generate stringer -type ErrType -trimprefix Err -output error_string.go
type ErrType int

const (
	ErrRequestFail ErrType = iota
	ErrAuthFail
	ErrObjectNotExist
	ErrAccountNotExist
)

func IsRequestFail(err error) bool { return isError(err, ErrRequestFail) }

func IsAuthFail(err error) bool { return isError(err, ErrAuthFail) }

func IsObjectNotExist(err error) bool { return isError(err, ErrObjectNotExist) }

func IsAccountNotExist(err error) bool { return isError(err, ErrAccountNotExist) }
