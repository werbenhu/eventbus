// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 werbenhu
// SPDX-FileContributor: werbenhu

package eventbus

type err struct {
	Msg  string
	Code int
}

func (e err) String() string {
	return e.Msg
}

func (e err) Error() string {
	return e.Msg
}

var (
	ErrHandlerIsNotFunc  = err{Code: 10000, Msg: "handler is not a function"}
	ErrHandlerParamNum   = err{Code: 10001, Msg: "the number of parameters of the handler must be two"}
	ErrHandlerFirstParam = err{Code: 10002, Msg: "the first of parameters of the handler must be a string"}
	ErrNoSubscriber      = err{Code: 10003, Msg: "no subscriber on topic"}
	ErrChannelClosed     = err{Code: 10004, Msg: "channel is closed"}

	ErrGroupExisted    = err{Code: 10001, Msg: "group already existed"}
	ErrNoResultMatched = err{Code: 10002, Msg: "no result matched"}
	ErrKeyExisted      = err{Code: 10003, Msg: "key already existed"}
)
