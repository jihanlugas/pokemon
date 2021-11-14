package response

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
)

var json jsoniter.API

func init() {
	json = jsoniter.ConfigFastest
}

// ErrorResponse for error type when request validation failed
type ErrorResponse struct {
	IsError bool        `json:"error"`
	Message string      `json:"message"`
	Payload interface{} `json:"payload"`
}

// SuccessResponse type for Success Response
type SuccessResponse struct {
	IsSuccess bool        `json:"success"`
	Message   string      `json:"message"`
	Payload   interface{} `json:"payload"`
}

// Payload used to pass to Error or Success method
type Payload map[string]interface{}

func (e *ErrorResponse) Error() string {
	return e.Message
}

// Error generate ErrorResponse error
func Error(msg string, payload interface{}) *ErrorResponse {
	return &ErrorResponse{true, msg, payload}
}

// ErrorForce generate ErrorResponse error type with additional forceLogout:true in payload
func ErrorForce(msg string, payload Payload) *ErrorResponse {
	payload["forceLogout"] = true
	return &ErrorResponse{true, msg, payload}
}

// Success generate SuccessResponse with given payload
func Success(msg string, payload interface{}) *SuccessResponse {
	return &SuccessResponse{true, msg, payload}
}

// SendJSON response to the client browser with Plain context.String
func (s *SuccessResponse) SendJSON(c echo.Context) error {
	return sendJSON(c, s, http.StatusOK)
}

// SendJSON response to the client browser with Plain context.String
func (e *ErrorResponse) SendJSON(c echo.Context) error {
	return sendJSON(c, e, http.StatusBadRequest)
}

func sendJSON(c echo.Context, i interface{}, httpStat int) error {

	if js, err := json.Marshal(i); err != nil {
		panic(err)
	} else {
		return c.Blob(httpStat, echo.MIMEApplicationJSONCharsetUTF8, js)
	}
}

func ValidationError(err error) *Payload {
	return &Payload{
		"listError": getListError(err),
	}
}

type ListErrorStack struct {
	listError Payload
}

func (e *ListErrorStack) StackError(field, msg string) *ListErrorStack {
	e.listError[field] = FieldError{
		Field: field,
		Msg:   msg,
	}
	return e
}

func (e *ListErrorStack) Build() *Payload {
	return &Payload{
		"listError": e.listError,
	}
}

func ListErrorComposer() *ListErrorStack {
	return &ListErrorStack{listError: make(Payload)}
}
