package model

// 400, 408(timeout), 422(already) // 500(Internal server)
const INTERNAL_SERVER_ERR = 500
const ALREADY_USED = 422
const AUTH_ERR = 401
const SUCCESS = 200
const BAD_REQUEST = 400

type ServerError struct {
	Status	int
	Msg		string
}

type Option func(*ServerError)

func Status(status int) Option {
	return func(args *ServerError) {
		args.Status = status
	}
}

func Msg(msg string) Option {
	return func(args *ServerError) {
		args.Msg= msg
	}
}

func NewServerError(options ...Option) *ServerError {
	pServerErr := &ServerError{
		Status: 400,
		Msg: 	"Parameter is incorrect",
	}
	for _, option := range options {
		option(pServerErr)
	}
	return pServerErr
}
