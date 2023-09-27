// Package main provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.15.0 DO NOT EDIT.
package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
)

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// LoginBody defines model for LoginBody.
type LoginBody struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// LoginResponse defines model for LoginResponse.
type LoginResponse struct {
	Token string `json:"token"`
}

// UserResponse defines model for UserResponse.
type UserResponse struct {
	Username string `json:"username"`
}

// ValueResponse defines model for ValueResponse.
type ValueResponse struct {
	Value *string `json:"value,omitempty"`
}

// PostLoginFormdataRequestBody defines body for PostLogin for application/x-www-form-urlencoded ContentType.
type PostLoginFormdataRequestBody = LoginBody

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// PostLoginWithBody request with any body
	PostLoginWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	PostLoginWithFormdataBody(ctx context.Context, body PostLoginFormdataRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetPing request
	GetPing(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetUsersMe request
	GetUsersMe(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) PostLoginWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostLoginRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PostLoginWithFormdataBody(ctx context.Context, body PostLoginFormdataRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostLoginRequestWithFormdataBody(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetPing(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetPingRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetUsersMe(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetUsersMeRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewPostLoginRequestWithFormdataBody calls the generic PostLogin builder with application/x-www-form-urlencoded body
func NewPostLoginRequestWithFormdataBody(server string, body PostLoginFormdataRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	bodyStr, err := runtime.MarshalForm(body, nil)
	if err != nil {
		return nil, err
	}
	bodyReader = strings.NewReader(bodyStr.Encode())
	return NewPostLoginRequestWithBody(server, "application/x-www-form-urlencoded", bodyReader)
}

// NewPostLoginRequestWithBody generates requests for PostLogin with any type of body
func NewPostLoginRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/login")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewGetPingRequest generates requests for GetPing
func NewGetPingRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/ping")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetUsersMeRequest generates requests for GetUsersMe
func NewGetUsersMeRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/users/me")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// PostLoginWithBodyWithResponse request with any body
	PostLoginWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostLoginResponse, error)

	PostLoginWithFormdataBodyWithResponse(ctx context.Context, body PostLoginFormdataRequestBody, reqEditors ...RequestEditorFn) (*PostLoginResponse, error)

	// GetPingWithResponse request
	GetPingWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetPingResponse, error)

	// GetUsersMeWithResponse request
	GetUsersMeWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetUsersMeResponse, error)
}

type PostLoginResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *LoginResponse
	JSON401      *ErrorResponse
}

// Status returns HTTPResponse.Status
func (r PostLoginResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r PostLoginResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetPingResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *ValueResponse
}

// Status returns HTTPResponse.Status
func (r GetPingResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetPingResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetUsersMeResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *UserResponse
}

// Status returns HTTPResponse.Status
func (r GetUsersMeResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetUsersMeResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// PostLoginWithBodyWithResponse request with arbitrary body returning *PostLoginResponse
func (c *ClientWithResponses) PostLoginWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostLoginResponse, error) {
	rsp, err := c.PostLoginWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostLoginResponse(rsp)
}

func (c *ClientWithResponses) PostLoginWithFormdataBodyWithResponse(ctx context.Context, body PostLoginFormdataRequestBody, reqEditors ...RequestEditorFn) (*PostLoginResponse, error) {
	rsp, err := c.PostLoginWithFormdataBody(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostLoginResponse(rsp)
}

// GetPingWithResponse request returning *GetPingResponse
func (c *ClientWithResponses) GetPingWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetPingResponse, error) {
	rsp, err := c.GetPing(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetPingResponse(rsp)
}

// GetUsersMeWithResponse request returning *GetUsersMeResponse
func (c *ClientWithResponses) GetUsersMeWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetUsersMeResponse, error) {
	rsp, err := c.GetUsersMe(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetUsersMeResponse(rsp)
}

// ParsePostLoginResponse parses an HTTP response from a PostLoginWithResponse call
func ParsePostLoginResponse(rsp *http.Response) (*PostLoginResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &PostLoginResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest LoginResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	}

	return response, nil
}

// ParseGetPingResponse parses an HTTP response from a GetPingWithResponse call
func ParseGetPingResponse(rsp *http.Response) (*GetPingResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetPingResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest ValueResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseGetUsersMeResponse parses an HTTP response from a GetUsersMeWithResponse call
func ParseGetUsersMeResponse(rsp *http.Response) (*GetUsersMeResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetUsersMeResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest UserResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (POST /login)
	PostLogin(ctx echo.Context) error

	// (GET /ping)
	GetPing(ctx echo.Context) error

	// (GET /users/me)
	GetUsersMe(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// PostLogin converts echo context to params.
func (w *ServerInterfaceWrapper) PostLogin(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostLogin(ctx)
	return err
}

// GetPing converts echo context to params.
func (w *ServerInterfaceWrapper) GetPing(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetPing(ctx)
	return err
}

// GetUsersMe converts echo context to params.
func (w *ServerInterfaceWrapper) GetUsersMe(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetUsersMe(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.POST(baseURL+"/login", wrapper.PostLogin)
	router.GET(baseURL+"/ping", wrapper.GetPing)
	router.GET(baseURL+"/users/me", wrapper.GetUsersMe)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/7RUTY/TMBD9K9bAMd20wCk3kFAFLGKFgAviYJJp6iX1mJlJS1XlvyM7/djQdBFoNycr",
	"M37v+fmNd1DSKpBHrwLFDqRc4sqm5Wtm4o8ogbxg/BGYArI6TOUVitg6FXQbEAoQZedr6DIQtdrKnZLz",
	"ijUydF0GjD9bx1hB8fXQmB3RvmWHLfT9FkuNaNdUO/+Kqu25iGBFNsTVqIpWkL1djUn8Q8axMzshXlRy",
	"2RKlH+j/zta3jeF/FrzH8f84zxjJF9u0eJllHctjFF1kcX5BsVihlOyCOvJQwNtW1CyIjaJo7M5AnTZx",
	"+ycUNYK8RoYM1sjSb5lFKRTQ2+CggOdX06tZsl+XSUbeRK+TPBI9p7ym2jhvlIwu0dgQGlfaVEuwnNZv",
	"KijghkTTxUFvEooewlSSV/QJ/Q5C/muy2WwmC+LVpOUGfUkVVqfpiKunjAso4El+Gp98Pzv5Ka/d8GKU",
	"W0w/eu/TOZ9Np/couRXy/0h8vNlEPjTtw7vo+ovp7MEoh4/ECOXLVpfodY9uFtY1WEEfpjzErBQ7qDHJ",
	"GN7bHPWmz9KjGTYchRH1B6FxoCTvR28vdtg4RzVly4xem61pqK6xivmMG88SOUeNky7v8TEPN3hMLoQh",
	"fb8DAAD//1Vuc4YEBgAA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
