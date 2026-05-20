package apierror

import (
	"context"
	"errors"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// CodedError carries a business error code and HTTP status.
type CodedError struct {
	Code       int
	HTTPStatus int
}

func (e *CodedError) Error() string {
	return fmt.Sprintf("error code %d (http %d)", e.Code, e.HTTPStatus)
}

// Predefined business errors.
var (
	ErrBadRequest      = &CodedError{Code: 4000, HTTPStatus: consts.StatusBadRequest}
	ErrSMSTooFrequent  = &CodedError{Code: 4001, HTTPStatus: consts.StatusTooManyRequests}
	ErrSMSDailyLimit   = &CodedError{Code: 4002, HTTPStatus: consts.StatusTooManyRequests}
	ErrSMSCodeWrong    = &CodedError{Code: 4003, HTTPStatus: consts.StatusBadRequest}
	ErrAccountLocked   = &CodedError{Code: 4004, HTTPStatus: consts.StatusForbidden}
	ErrSMSCodeExpired  = &CodedError{Code: 4005, HTTPStatus: consts.StatusBadRequest}
	ErrScanNotFound    = &CodedError{Code: 4006, HTTPStatus: consts.StatusNotFound}
	ErrBindTicketBad   = &CodedError{Code: 4007, HTTPStatus: consts.StatusBadRequest}
	ErrTokenExpired    = &CodedError{Code: 4008, HTTPStatus: consts.StatusUnauthorized}
	ErrWxCodeBad       = &CodedError{Code: 4009, HTTPStatus: consts.StatusBadRequest}
	ErrIPRateLimit     = &CodedError{Code: 4010, HTTPStatus: consts.StatusTooManyRequests}
	ErrForbidden       = &CodedError{Code: 4011, HTTPStatus: consts.StatusForbidden}
	ErrNotFound        = &CodedError{Code: 4012, HTTPStatus: consts.StatusNotFound}
	ErrRateLimit       = &CodedError{Code: 4013, HTTPStatus: consts.StatusTooManyRequests}
	ErrAnonIDInvalid   = &CodedError{Code: 4014, HTTPStatus: consts.StatusBadRequest}
	ErrAuthCodeExpired = &CodedError{Code: 4015, HTTPStatus: consts.StatusBadRequest}
	ErrAdminLoginInvalid = &CodedError{Code: 4016, HTTPStatus: consts.StatusUnauthorized}
	ErrAgentError      = &CodedError{Code: 5001, HTTPStatus: consts.StatusInternalServerError}
	ErrInternal        = &CodedError{Code: 5000, HTTPStatus: consts.StatusInternalServerError}
)

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var defaultMessages = map[int]string{
	4000: "bad request",
	4001: "sms send too frequent",
	4002: "sms daily limit exceeded",
	4003: "sms code wrong",
	4004: "account locked",
	4005: "sms code expired",
	4006: "scan session not found",
	4007: "bind ticket invalid",
	4008: "token expired or invalid",
	4009: "wechat code invalid",
	4010: "ip rate limit",
	4011: "forbidden",
	4012: "not found",
	4013: "请求过于频繁，请稍后再试",
	4014: "invalid anonymous id",
	4015: "auth code expired",
	4016: "admin login invalid",
	5000: "internal error",
	5001: "agent error",
}

// Render writes a JSON error response.
func Render(_ context.Context, c *app.RequestContext, err error) {
	var coded *CodedError
	if !errors.As(err, &coded) {
		coded = ErrInternal
	}
	msg := defaultMessages[coded.Code]
	if msg == "" {
		msg = "error"
	}
	c.JSON(coded.HTTPStatus, errorResponse{Code: coded.Code, Message: msg})
}
