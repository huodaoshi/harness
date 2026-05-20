package domain

import (
	"context"
	"time"
)

// User is the core user entity.
type User struct {
	UserID    string    `bson:"user_id"`
	UID       int64     `bson:"uid"`
	Phone     string    `bson:"phone"`
	WxOpenID  string    `bson:"wx_openid"`
	WxUnionID string    `bson:"wx_unionid"`
	AnonID    string    `bson:"anon_id"`
	Role      int       `bson:"role"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

// ScanSession holds WeChat QR scan session state.
type ScanSession struct {
	ScanToken  string
	Status     int
	AnonID     string
	UserID     string
	WxOpenID   string
	WxUnionID  string
	BindTicket string
}

// TokenPair holds a JWT access/refresh token pair.
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

// RefreshPayload carries the claims embedded in a refresh token.
type RefreshPayload struct {
	UserID string
	UID    int64
	Role   int
}

// UserRepo is the repository interface for user persistence.
type UserRepo interface {
	Create(ctx context.Context, u *User) error
	FindByID(ctx context.Context, userID string) (*User, error)
	FindByPhone(ctx context.Context, phone string) (*User, error)
	FindByWxOpenID(ctx context.Context, openID string) (*User, error)
	UpdateWxInfo(ctx context.Context, userID, openID, unionID string) error
	MigrateAnonSessions(ctx context.Context, anonID, userID string) error
	UpsertStaticAdmin(ctx context.Context, userID string, uid int64, role UserRole) error
}

// ScanRepo is the repository interface for WeChat QR scan sessions.
type ScanRepo interface {
	Set(ctx context.Context, session *ScanSession, ttlSeconds int) error
	Get(ctx context.Context, scanToken string) (*ScanSession, error)
	Confirm(ctx context.Context, session *ScanSession) error
}

// SMSCodeRepo is the repository interface for SMS verification codes and rate limits.
type SMSCodeRepo interface {
	SetCode(ctx context.Context, phone, code string, ttlSeconds int) error
	GetCode(ctx context.Context, phone string) (string, error)
	IncrAttempts(ctx context.Context, phone string) (int64, error)
	SetSendLimit(ctx context.Context, phone string, ttlSeconds int) error
	CheckSendLimit(ctx context.Context, phone string) (bool, error)
	IncrDailyCount(ctx context.Context, phone string) (int64, error)
	IncrIPCount(ctx context.Context, ip string, ttlSeconds int) (int64, error)
	SetLock(ctx context.Context, phone string, ttlSeconds int) error
	CheckLock(ctx context.Context, phone string) (bool, error)
	DeleteCode(ctx context.Context, phone string) error
}

// TokenRepo is the repository interface for refresh token storage.
type TokenRepo interface {
	SetRefresh(ctx context.Context, userID, token string, ttlSeconds int) error
	GetRefresh(ctx context.Context, userID string) (string, error)
	DeleteRefresh(ctx context.Context, userID string) error
}

// BindTicketRepo is the repository interface for WeChat bind ticket storage.
type BindTicketRepo interface {
	Set(ctx context.Context, ticket, openID, unionID, anonID string, ttlSeconds int) error
	Get(ctx context.Context, ticket string) (openID, unionID, anonID string, err error)
	Delete(ctx context.Context, ticket string) error
}

// WeChatClient is the interface for WeChat platform API calls.
type WeChatClient interface {
	Code2Session(ctx context.Context, code string) (openID, unionID string, err error)
}

// SMSService is the interface for sending SMS messages.
type SMSService interface {
	Send(ctx context.Context, phone, code string) error
}

// AuthCodeRepo is the repository interface for short-lived auth codes.
type AuthCodeRepo interface {
	SetAuthCode(ctx context.Context, authCode string, pair *TokenPair) error
	ConsumeAuthCode(ctx context.Context, authCode string) (*TokenPair, error)
}
