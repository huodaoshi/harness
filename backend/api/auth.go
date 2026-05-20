package api

import (
	"context"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/huodaoshi/harness/backend/modules/auth/application"
	"github.com/huodaoshi/harness/backend/modules/auth/domain"
	"github.com/huodaoshi/harness/backend/pkg/apierror"
)

// ---------------------------------------------------------------------------
// Request structs
// ---------------------------------------------------------------------------

// SendSMSRequest is the request body for POST /v1/auth/sms/send.
type SendSMSRequest struct {
	Phone string `json:"phone" binding:"required"`
}

// VerifySMSRequest is the request body for POST /v1/auth/sms/verify.
type VerifySMSRequest struct {
	Phone  string `json:"phone"   binding:"required"`
	Code   string `json:"code"    binding:"required"`
	AnonID string `json:"anon_id"`
}

// InitWxScanRequest is the request body for POST /v1/auth/wx/scan/init.
type InitWxScanRequest struct {
	AnonID string `json:"anon_id"`
}

// ConfirmWxScanRequest is the request body for POST /v1/auth/wx/scan/confirm.
type ConfirmWxScanRequest struct {
	ScanToken string `json:"scan_token" binding:"required"`
	WxCode    string `json:"wx_code"    binding:"required"`
	AnonID    string `json:"anon_id"`
}

// MiniProgramLoginRequest is the request body for POST /v1/auth/wx/miniprogram/login.
type MiniProgramLoginRequest struct {
	WxCode string `json:"wx_code" binding:"required"`
	AnonID string `json:"anon_id"`
}

// BindPhoneRequest is the request body for POST /v1/auth/wx/bind.
type BindPhoneRequest struct {
	BindTicket string `json:"bind_ticket" binding:"required"`
	Phone      string `json:"phone"       binding:"required"`
	Code       string `json:"code"        binding:"required"`
	AnonID     string `json:"anon_id"`
}

// RefreshTokenRequest is the request body for POST /v1/auth/token/refresh.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// LogoutRequest is the request body for POST /v1/auth/logout.
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// AdminStaticLoginRequest is the request body for POST /v1/auth/admin/login.
type AdminStaticLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ExchangeAuthCodeRequest is the request body for POST /v1/auth/wx/scan/exchange.
type ExchangeAuthCodeRequest struct {
	AuthCode string `json:"auth_code" binding:"required"`
}

// ---------------------------------------------------------------------------
// Response structs
// ---------------------------------------------------------------------------

// TokenPairResponse is the JSON representation of a token pair.
type TokenPairResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

// InitScanResponse is returned by POST /v1/auth/wx/scan/init.
type InitScanResponse struct {
	ScanToken string `json:"scan_token"`
	QRContent string `json:"qr_content"`
}

// ScanStatusResponse is returned by GET /v1/auth/wx/scan/status/:scan_token.
type ScanStatusResponse struct {
	Status     int    `json:"status"`
	AuthCode   string `json:"auth_code,omitempty"`
	BindTicket string `json:"bind_ticket,omitempty"`
	WxOpenID   string `json:"wx_openid,omitempty"`
}

// MiniProgramLoginResponse is returned by POST /v1/auth/wx/miniprogram/login.
type MiniProgramLoginResponse struct {
	TokenPair  *TokenPairResponse `json:"token_pair,omitempty"`
	BindTicket string             `json:"bind_ticket,omitempty"`
	WxOpenID   string             `json:"wx_openid,omitempty"`
}

// ---------------------------------------------------------------------------
// AuthHandler
// ---------------------------------------------------------------------------

// AuthHandler handles all authentication HTTP requests.
type AuthHandler struct {
	svc application.AuthService
}

// NewAuthHandler creates an AuthHandler.
func NewAuthHandler(svc application.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// HandleSendSMS handles POST /v1/auth/sms/send.
func (h *AuthHandler) HandleSendSMS(ctx context.Context, c *app.RequestContext) {
	var req SendSMSRequest
	if err := c.BindAndValidate(&req); err != nil {
		apierror.Render(ctx, c, apierror.ErrBadRequest)
		return
	}

	ip := extractIP(c)
	if err := h.svc.SendSMSCode(ctx, req.Phone, ip); err != nil {
		apierror.Render(ctx, c, err)
		return
	}

	c.Status(consts.StatusNoContent)
}

// HandleVerifySMS handles POST /v1/auth/sms/verify.
func (h *AuthHandler) HandleVerifySMS(ctx context.Context, c *app.RequestContext) {
	var req VerifySMSRequest
	if err := c.BindAndValidate(&req); err != nil {
		apierror.Render(ctx, c, apierror.ErrBadRequest)
		return
	}

	pair, err := h.svc.VerifySMSCode(ctx, req.Phone, req.Code, req.AnonID)
	if err != nil {
		apierror.Render(ctx, c, err)
		return
	}

	c.JSON(consts.StatusOK, toTokenPairResponse(pair))
}

// HandleInitWxScan handles POST /v1/auth/wx/scan/init.
func (h *AuthHandler) HandleInitWxScan(ctx context.Context, c *app.RequestContext) {
	var req InitWxScanRequest
	if err := c.BindAndValidate(&req); err != nil {
		apierror.Render(ctx, c, apierror.ErrBadRequest)
		return
	}

		session, err := h.svc.InitWxScan(ctx, req.AnonID)
	if err != nil {
		apierror.Render(ctx, c, err)
		return
	}

	c.JSON(consts.StatusOK, InitScanResponse{
		ScanToken: session.ScanToken,
		QRContent: "weixin://wxpay/bizpayurl?scan_token=" + session.ScanToken,
	})
}

// HandleConfirmWxScan handles POST /v1/auth/wx/scan/confirm.
func (h *AuthHandler) HandleConfirmWxScan(ctx context.Context, c *app.RequestContext) {
	var req ConfirmWxScanRequest
	if err := c.BindAndValidate(&req); err != nil {
		apierror.Render(ctx, c, apierror.ErrBadRequest)
		return
	}

		if err := h.svc.ConfirmWxScan(ctx, req.ScanToken, req.WxCode, req.AnonID); err != nil {
		apierror.Render(ctx, c, err)
		return
	}

	c.Status(consts.StatusNoContent)
}

// HandleGetScanStatus handles GET /v1/auth/wx/scan/status/:scan_token.
func (h *AuthHandler) HandleGetScanStatus(ctx context.Context, c *app.RequestContext) {
	scanToken := c.Param("scan_token")
	if scanToken == "" {
		apierror.Render(ctx, c, apierror.ErrScanNotFound)
		return
	}

	result, err := h.svc.GetScanStatus(ctx, scanToken)
	if err != nil {
		apierror.Render(ctx, c, err)
		return
	}

	resp := ScanStatusResponse{
		Status:     result.Status,
		AuthCode:   result.AuthCode,
		BindTicket: result.BindTicket,
		WxOpenID:   result.WxOpenID,
	}

	c.JSON(consts.StatusOK, resp)
}

// HandleMiniProgramLogin handles POST /v1/auth/wx/miniprogram/login.
func (h *AuthHandler) HandleMiniProgramLogin(ctx context.Context, c *app.RequestContext) {
	var req MiniProgramLoginRequest
	if err := c.BindAndValidate(&req); err != nil {
		apierror.Render(ctx, c, apierror.ErrWxCodeBad)
		return
	}

	result, err := h.svc.MiniProgramLogin(ctx, req.WxCode, req.AnonID)
	if err != nil {
		apierror.Render(ctx, c, err)
		return
	}

	resp := MiniProgramLoginResponse{
		BindTicket: result.BindTicket,
		WxOpenID:   result.WxOpenID,
	}
	if result.TokenPair != nil {
		resp.TokenPair = toTokenPairResponse(result.TokenPair)
	}

	c.JSON(consts.StatusOK, resp)
}

// HandleBindPhone handles POST /v1/auth/wx/bind.
func (h *AuthHandler) HandleBindPhone(ctx context.Context, c *app.RequestContext) {
	var req BindPhoneRequest
	if err := c.BindAndValidate(&req); err != nil {
		apierror.Render(ctx, c, apierror.ErrBadRequest)
		return
	}

		pair, err := h.svc.BindPhone(ctx, req.BindTicket, req.Phone, req.Code, req.AnonID)
	if err != nil {
		apierror.Render(ctx, c, err)
		return
	}

	c.JSON(consts.StatusOK, toTokenPairResponse(pair))
}

// HandleRefreshToken handles POST /v1/auth/token/refresh.
func (h *AuthHandler) HandleRefreshToken(ctx context.Context, c *app.RequestContext) {
	var req RefreshTokenRequest
	if err := c.BindAndValidate(&req); err != nil {
		apierror.Render(ctx, c, apierror.ErrTokenExpired)
		return
	}

	pair, err := h.svc.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		apierror.Render(ctx, c, err)
		return
	}

	c.JSON(consts.StatusOK, toTokenPairResponse(pair))
}

// HandleLogout handles POST /v1/auth/logout.
func (h *AuthHandler) HandleLogout(ctx context.Context, c *app.RequestContext) {
	var req LogoutRequest
	if err := c.BindAndValidate(&req); err != nil {
		apierror.Render(ctx, c, apierror.ErrTokenExpired)
		return
	}

	userID, ok := UserIDFromContext(c)
	if !ok {
		apierror.Render(ctx, c, apierror.ErrTokenExpired)
		return
	}

	if err := h.svc.Logout(ctx, userID); err != nil {
		apierror.Render(ctx, c, err)
		return
	}

	c.Status(consts.StatusNoContent)
}

// HandleExchangeAuthCode handles POST /v1/auth/wx/scan/exchange.
// It exchanges a short-lived auth_code for a full TokenPair.
func (h *AuthHandler) HandleExchangeAuthCode(ctx context.Context, c *app.RequestContext) {
	var req ExchangeAuthCodeRequest
	if err := c.BindAndValidate(&req); err != nil {
		apierror.Render(ctx, c, apierror.ErrBadRequest)
		return
	}

	pair, err := h.svc.ExchangeAuthCode(ctx, req.AuthCode)
	if err != nil {
		apierror.Render(ctx, c, err)
		return
	}

	c.JSON(consts.StatusOK, toTokenPairResponse(pair))
}

// HandleAdminStaticLogin handles POST /v1/auth/admin/login (fixed admin user, conf admin_static_login).
func (h *AuthHandler) HandleAdminStaticLogin(ctx context.Context, c *app.RequestContext) {
	var req AdminStaticLoginRequest
	if err := c.BindAndValidate(&req); err != nil {
		apierror.Render(ctx, c, apierror.ErrBadRequest)
		return
	}

	pair, err := h.svc.AdminStaticLogin(ctx, req.Username, req.Password)
	if err != nil {
		apierror.Render(ctx, c, err)
		return
	}

	c.JSON(consts.StatusOK, toTokenPairResponse(pair))
}

// ---------------------------------------------------------------------------
// Route registration
// ---------------------------------------------------------------------------

// RegisterAuthRoutes registers all /v1/auth routes on the Hertz server.
func RegisterAuthRoutes(r *server.Hertz, h *AuthHandler, jwtAuth app.HandlerFunc) {
	v1 := r.Group("/v1/auth")
	v1.POST("/admin/login", h.HandleAdminStaticLogin)
	v1.POST("/sms/send", h.HandleSendSMS)
	v1.POST("/sms/verify", h.HandleVerifySMS)
	v1.POST("/wx/scan/init", h.HandleInitWxScan)
	v1.POST("/wx/scan/confirm", h.HandleConfirmWxScan)
	v1.GET("/wx/scan/status/:scan_token", h.HandleGetScanStatus)
	v1.POST("/wx/scan/exchange", h.HandleExchangeAuthCode)
	v1.POST("/wx/miniprogram/login", h.HandleMiniProgramLogin)
	v1.POST("/wx/bind", h.HandleBindPhone)
	v1.POST("/token/refresh", h.HandleRefreshToken)
	v1.POST("/logout", jwtAuth, h.HandleLogout)
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// toTokenPairResponse converts a domain.TokenPair to the API response type.
func toTokenPairResponse(pair *domain.TokenPair) *TokenPairResponse {
	if pair == nil {
		return nil
	}
	return &TokenPairResponse{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
		ExpiresIn:    pair.ExpiresIn,
	}
}

// extractIP returns the client IP address from well-known proxy headers,
// falling back to c.ClientIP().
func extractIP(c *app.RequestContext) string {
	if ip := string(c.GetHeader("X-Real-IP")); ip != "" {
		return ip
	}
	if forwarded := string(c.GetHeader("X-Forwarded-For")); forwarded != "" {
		parts := strings.SplitN(forwarded, ",", 2)
		if ip := strings.TrimSpace(parts[0]); ip != "" {
			return ip
		}
	}
	return c.ClientIP()
}
