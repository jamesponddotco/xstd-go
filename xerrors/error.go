package xerrors

// Error is an [imuutable error] type.
//
// [imuutable error]: https://dave.cheney.net/2016/04/07/constant-errors
type Error string

// Error implements the error interface for Error.
func (e Error) Error() string {
	return string(e)
}
