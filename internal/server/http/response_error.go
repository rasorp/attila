// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package http

type ResponseError struct {
	ErrorBody `json:"error"`
}

type ErrorBody struct {
	Msg  string `json:"message"`
	Code int    `json:"code"`
}

func NewResponseError(e error, c int) *ResponseError {
	return &ResponseError{
		ErrorBody: ErrorBody{
			Msg:  e.Error(),
			Code: c,
		},
	}
}

func (e *ResponseError) StatusCode() int { return e.Code }

func (e *ResponseError) Error() string { return e.Msg }

func (e *ResponseError) String() string { return e.Msg }
