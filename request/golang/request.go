package main

type ResponseBase struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type LoginRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UserAccount  string `json:"account"`
	UserName     string `json:"name"`
	AccessToken  string `json:"accessToken"`
	AccessExpire int64  `json:"accessExpire"`
	RefreshAfter int64  `json:"refreshAfter"`
}

type LoginResponseBody struct {
	ResponseBase
	Data LoginResponse `json:"data,omitempty"`
}

type CheckResponse struct {
	Account   string `json:"account"`
	Remaining int    `json:"remaining"`
	Balance   int    `json:"balance"`
}

type CheckResponseBody struct {
	ResponseBase
	Data CheckResponse `json:"data,omitempty"`
}

type OpenResponse struct {
	Account   string `json:"account"`
	Amount    int    `json:"amount"`
	Remaining int    `json:"remaining"`
	Balance   int    `json:"balance"`
}

type OpenResponseBody struct {
	ResponseBase
	Data OpenResponse `json:"data,omitempty"`
}
