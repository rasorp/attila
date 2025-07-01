// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

const (
	Version = "v1alpha1"

	defaultUserAgent = "go-attila" + "/" + Version

	defaultAddress = "http://127.0.0.1:8080"
)

type ResponseError struct {
	ErrorBody `json:"error"`
}

type ErrorBody struct {
	Msg  string `json:"message"`
	Code int    `json:"code"`
}

func (e *ResponseError) Error() string {
	return e.Msg
}

func (e *ResponseError) StatusCode() int {
	return e.Code
}

type Response struct {
	*http.Response
}

type Config struct {
	Address    string
	HTTPClient *http.Client
	UserAgent  string
}

func DefaultConfig() *Config {
	return &Config{
		Address:   defaultAddress,
		UserAgent: defaultUserAgent,
	}
}

type Client struct {
	client    *http.Client
	address   string
	userAgent string
}

func NewClient(cfg *Config) *Client {

	address := cfg.Address
	if address == "" {
		address = defaultAddress
	}

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	userAgent := cfg.UserAgent
	if userAgent == "" {
		userAgent = defaultUserAgent
	}

	return &Client{
		address:   address,
		client:    httpClient,
		userAgent: userAgent,
	}
}

type RequestOption func(req *http.Request)

func (c *Client) NewRequest(method, path string, body interface{}, opts ...RequestOption) (*http.Request, error) {

	if !strings.HasPrefix(path, "/") {
		return nil, errors.New("path missing slash prefix")
	}

	var buf io.ReadWriter

	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	fullURL := c.address + path

	req, err := http.NewRequest(method, fullURL, buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("User-Agent", c.userAgent)

	for _, opt := range opts {
		opt(req)
	}

	return req, nil
}

func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.bareDo(ctx, req)
	if err != nil {
		return resp, err
	}

	defer func(Body io.ReadCloser) { _ = Body.Close() }(resp.Body)

	switch v := v.(type) {
	case nil:
	case io.Writer:
		_, err = io.Copy(v, resp.Body)
	default:
		decErr := json.NewDecoder(resp.Body).Decode(v)
		if decErr == io.EOF {
			decErr = nil // ignore EOF errors caused by empty response body
		}
		if decErr != nil {
			err = decErr
		}
	}
	return resp, err
}

func (c *Client) bareDo(ctx context.Context, req *http.Request) (*Response, error) {

	if ctx == nil {
		return nil, errors.New("context must be non-nil")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			return nil, err
		}
	}

	if err = checkResponse(resp); err != nil {
		return nil, err
	}
	return &Response{resp}, nil
}

func checkResponse(r *http.Response) error {

	if c := r.StatusCode; http.StatusOK <= c && c < http.StatusMultipleChoices {
		return nil
	}

	errorResponse := ResponseError{
		ErrorBody: ErrorBody{
			Code: r.StatusCode,
		},
	}

	data, err := io.ReadAll(r.Body)
	if err == nil && data != nil {
		if err = json.Unmarshal(data, &errorResponse); err != nil {
			errorResponse.Msg = err.Error()
		}
	}

	return &errorResponse
}
