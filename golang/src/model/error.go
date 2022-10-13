package model

// 400, 408(timeout), 422(already) // 500(Internal server)
type ServerError struct {
	Status	string
	Msg		string
}

type Option func(*ServerError)

func Status(status string) Option {
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
		Status: "400",
		Msg: 	"Parameter is incorrect",
	}
	for _, option := range options {
		option(pServerErr)
	}
	return pServerErr
}
