// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package state

type ErrorResp struct {
	ErrorBody `json:"error"`
}

type ErrorBody struct {
	Msg  string `json:"message"`
	Code int    `json:"code"`
	err  error
}

func NewErrorResp(e error, c int) *ErrorResp {
	return &ErrorResp{
		ErrorBody: ErrorBody{
			err:  e,
			Code: c,
			Msg:  e.Error(),
		},
	}
}

func (e *ErrorResp) Error() string {
	return e.Msg
}

func (e *ErrorResp) Err() error {
	return e.err
}

func (e *ErrorResp) StatusCode() int {
	return e.Code
}

func (e *ErrorResp) String() string {
	return e.Msg
}
