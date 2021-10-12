package errors

type HTTPError interface {
	Code() int
}
