// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 werbenhu
// SPDX-FileContributor: werbenhu

package eventbus

type err struct {
	Msg  string
	Code int
}

// String return the error's message
func (e err) String() string {
	return e.Msg
}

// Error return the error's message
func (e err) Error() string {
	return e.Msg
}

// Global variables that represent common errors that may be
// returned by the eventbus functions.
var (
	ErrHandlerIsNotFunc  = err{Code: 10000, Msg: "handler is not a function"}
	ErrHandlerParamNum   = err{Code: 10001, Msg: "the number of parameters of the handler must be two"}
	ErrHandlerFirstParam = err{Code: 10002, Msg: "the first of parameters of the handler must be a string"}
	ErrNoSubscriber      = err{Code: 10003, Msg: "no subscriber on topic"}
	ErrChannelClosed     = err{Code: 10004, Msg: "channel is closed"}
)
