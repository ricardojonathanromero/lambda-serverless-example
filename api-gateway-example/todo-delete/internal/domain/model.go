package domain

type ErrRes struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewErr(code, msg string) *ErrRes {
	return &ErrRes{Code: code, Message: msg}
}
