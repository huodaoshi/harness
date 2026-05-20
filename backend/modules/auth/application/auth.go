package application

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"

	authpkg "github.com/huodaoshi/harness/backend/modules/auth/pkg"

	"github.com/huodaoshi/harness/backend/modules/auth/domain"
	"github.com/huodaoshi/harness/backend/pkg/apierror"
)

// ScanStatus values.
const (
	ScanStatusPending  = 0
	ScanStatusScanned  = 1
	ScanStatusConfirmed = 2
)

// TTL constants (seconds).
const (
	ttlSMSCode      = 300  // 5 minutes
	ttlSMSSendLimit = 60   // 1 minute cool-down
	ttlSMSLock      = 600  // 10 minutes lock
	ttlScanSession  = 300 // 5 minutes
	ttlBindTicket   = 600 // 10 minutes

	maxSMSAttempts  = 5
	maxSMSDailyCount = 10
	maxIPCount      = 20
	ipTTL           = 3600 // 1 hour
)

// ScanStatusResult is the return type for GetScanStatus.
type ScanStatusResult struct {
	Status     int
	AuthCode   string
	BindTicket string
	WxOpenID   string
}

// LoginResult is the return type for WeChat Mini Program login.
type LoginResult struct {
	TokenPair  *domain.TokenPair
	BindTicket string
	WxOpenID   string
}

// AuthService defines all authentication operations.
type AuthService interface {
	SendSMSCode(ctx context.Context, phone, ip string) error
	VerifySMSCode(ctx context.Context, phone, code, anonID string) (*domain.TokenPair, error)
	InitWxScan(ctx context.Context, anonID string) (*domain.ScanSession, error)
	ConfirmWxScan(ctx context.Context, scanToken, wxCode, anonID string) error
	GetScanStatus(ctx context.Context, scanToken string) (*ScanStatusResult, error)
	MiniProgramLogin(ctx context.Context, wxCode, anonID string) (*LoginResult, error)
	BindPhone(ctx context.Context, bindTicket, phone, code, anonID string) (*domain.TokenPair, error)
	RefreshToken(ctx context.Context, refreshToken string) (*domain.TokenPair, error)
	Logout(ctx context.Context, userID string) error
	ExchangeAuthCode(ctx context.Context, authCode string) (*domain.TokenPair, error)
	// AdminStaticLogin issues a TokenPair when fixed admin credentials match (see conf admin_static_login).
	AdminStaticLogin(ctx context.Context, username, password string) (*domain.TokenPair, error)
}

// AdminStaticLoginOpts configures optional fixed admin username/password login.
type AdminStaticLoginOpts struct {
	Enabled  bool
	Username string
	Password string
	UserID   string
	UID      int64
}

func (o AdminStaticLoginOpts) active() bool {
	return o.Enabled && o.Username != "" && o.Password != "" && o.UserID != ""
}

// authService is the private implementation of AuthService.
type authService struct {
	userRepo     domain.UserRepo
	smsCodeRepo  domain.SMSCodeRepo
	tokenRepo    domain.TokenRepo
	scanRepo     domain.ScanRepo
	bindRepo     domain.BindTicketRepo
	wxClient     domain.WeChatClient
	signer       authpkg.JWTSigner
	smsService   domain.SMSService
	refreshTTL   int
	authCodeRepo domain.AuthCodeRepo
	adminStatic  AdminStaticLoginOpts
}

// NewAuthService creates an AuthService with all required dependencies.
func NewAuthService(
	userRepo domain.UserRepo,
	smsCodeRepo domain.SMSCodeRepo,
	tokenRepo domain.TokenRepo,
	scanRepo domain.ScanRepo,
	bindRepo domain.BindTicketRepo,
	wxClient domain.WeChatClient,
	signer authpkg.JWTSigner,
	smsService domain.SMSService,
	refreshTTL int,
	authCodeRepo domain.AuthCodeRepo,
	adminStatic AdminStaticLoginOpts,
) AuthService {
	return &authService{
		userRepo:     userRepo,
		smsCodeRepo:  smsCodeRepo,
		tokenRepo:    tokenRepo,
		scanRepo:     scanRepo,
		bindRepo:     bindRepo,
		wxClient:     wxClient,
		signer:       signer,
		smsService:   smsService,
		refreshTTL:   refreshTTL,
		authCodeRepo: authCodeRepo,
		adminStatic:  adminStatic,
	}
}

// SendSMSCode generates and sends a 6-digit SMS verification code.
func (s *authService) SendSMSCode(ctx context.Context, phone, ip string) error {
	// Check phone-level lock.
	locked, err := s.smsCodeRepo.CheckLock(ctx, phone)
	if err != nil {
		return fmt.Errorf("application: send sms: check lock: %w", err)
	}
	if locked {
		return apierror.ErrAccountLocked
	}

	// Check per-phone send frequency.
	limited, err := s.smsCodeRepo.CheckSendLimit(ctx, phone)
	if err != nil {
		return fmt.Errorf("application: send sms: check send limit: %w", err)
	}
	if limited {
		return apierror.ErrSMSTooFrequent
	}

	// Check daily quota.
	daily, err := s.smsCodeRepo.IncrDailyCount(ctx, phone)
	if err != nil {
		return fmt.Errorf("application: send sms: incr daily count: %w", err)
	}
	if daily > maxSMSDailyCount {
		return apierror.ErrSMSDailyLimit
	}

	// Check per-IP rate limit.
	ipCount, err := s.smsCodeRepo.IncrIPCount(ctx, ip, ipTTL)
	if err != nil {
		return fmt.Errorf("application: send sms: incr ip count: %w", err)
	}
	if ipCount > maxIPCount {
		return apierror.ErrIPRateLimit
	}

	// Generate 6-digit code.
	code, err := generate6DigitCode()
	if err != nil {
		return fmt.Errorf("application: send sms: generate code: %w", err)
	}

	// Store code.
	if err = s.smsCodeRepo.SetCode(ctx, phone, code, ttlSMSCode); err != nil {
		return fmt.Errorf("application: send sms: set code: %w", err)
	}

	// Set send frequency limit.
	if err = s.smsCodeRepo.SetSendLimit(ctx, phone, ttlSMSSendLimit); err != nil {
		return fmt.Errorf("application: send sms: set send limit: %w", err)
	}

	// Send SMS.
	if err = s.smsService.Send(ctx, phone, code); err != nil {
		return fmt.Errorf("application: send sms: send: %w", err)
	}

	return nil
}

// VerifySMSCode verifies an SMS code and returns a TokenPair on success.
func (s *authService) VerifySMSCode(ctx context.Context, phone, code, anonID string) (*domain.TokenPair, error) {
	// Check lock.
	locked, err := s.smsCodeRepo.CheckLock(ctx, phone)
	if err != nil {
		return nil, fmt.Errorf("application: verify sms: check lock: %w", err)
	}
	if locked {
		return nil, apierror.ErrAccountLocked
	}

	// Get stored code.
	stored, err := s.smsCodeRepo.GetCode(ctx, phone)
	if err != nil {
		return nil, fmt.Errorf("application: verify sms: get code: %w", err)
	}
	if stored == "" {
		return nil, apierror.ErrSMSCodeExpired
	}

	if stored != code {
		// Increment failure attempts.
		attempts, incrErr := s.smsCodeRepo.IncrAttempts(ctx, phone)
		if incrErr != nil {
			return nil, fmt.Errorf("application: verify sms: incr attempts: %w", incrErr)
		}
		if attempts >= maxSMSAttempts {
			if lockErr := s.smsCodeRepo.SetLock(ctx, phone, ttlSMSLock); lockErr != nil {
				return nil, fmt.Errorf("application: verify sms: set lock: %w", lockErr)
			}
			return nil, apierror.ErrAccountLocked
		}
		return nil, apierror.ErrSMSCodeWrong
	}

	// Code is correct — delete it.
	if err = s.smsCodeRepo.DeleteCode(ctx, phone); err != nil {
		return nil, fmt.Errorf("application: verify sms: delete code: %w", err)
	}

	// Find or create user.
	user, err := s.userRepo.FindByPhone(ctx, phone)
	if err != nil {
		return nil, fmt.Errorf("application: verify sms: find user: %w", err)
	}
	if user == nil {
		user = &domain.User{Phone: phone}
		if err = s.userRepo.Create(ctx, user); err != nil {
			return nil, fmt.Errorf("application: verify sms: create user: %w", err)
		}
	}

	// Issue tokens.
	pair, err := s.issueTokenPair(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("application: verify sms: issue tokens: %w", err)
	}

	// Migrate anon sessions asynchronously.
	s.migrateAnon(user.UserID, anonID)

	return pair, nil
}

// InitWxScan creates a new WeChat QR scan session.
func (s *authService) InitWxScan(ctx context.Context, anonID string) (*domain.ScanSession, error) {
	token, err := generateToken()
	if err != nil {
		return nil, fmt.Errorf("application: init wx scan: generate token: %w", err)
	}

	session := &domain.ScanSession{
		ScanToken: token,
		Status:    ScanStatusPending,
		AnonID:    anonID,
	}

	if err = s.scanRepo.Set(ctx, session, ttlScanSession); err != nil {
		return nil, fmt.Errorf("application: init wx scan: set session: %w", err)
	}

	return session, nil
}

// ConfirmWxScan is called by the WeChat client after scanning the QR code.
func (s *authService) ConfirmWxScan(ctx context.Context, scanToken, wxCode, anonID string) error {
	session, err := s.scanRepo.Get(ctx, scanToken)
	if err != nil {
		return fmt.Errorf("application: confirm wx scan: get session: %w", err)
	}
	if session == nil {
		return apierror.ErrScanNotFound
	}

	// Exchange wxCode for openID / unionID.
	openID, unionID, err := s.wxClient.Code2Session(ctx, wxCode)
	if err != nil {
		return apierror.ErrWxCodeBad
	}

	// Find or create user.
	user, err := s.userRepo.FindByWxOpenID(ctx, openID)
	if err != nil {
		return fmt.Errorf("application: confirm wx scan: find user: %w", err)
	}

	if user != nil {
		// Existing user — issue full tokens in session.
		pair, issueErr := s.issueTokenPair(ctx, user)
		if issueErr != nil {
			return fmt.Errorf("application: confirm wx scan: issue tokens: %w", issueErr)
		}
		session.Status = ScanStatusConfirmed
		session.UserID = user.UserID
		session.WxOpenID = openID
		session.WxUnionID = unionID
		// Store access token in bind_ticket field for retrieval convenience.
		session.BindTicket = pair.AccessToken
	} else {
		// New WeChat user — need phone binding.
		ticket, ticketErr := generateToken()
		if ticketErr != nil {
			return fmt.Errorf("application: confirm wx scan: generate bind ticket: %w", ticketErr)
		}
		if setErr := s.bindRepo.Set(ctx, ticket, openID, unionID, anonID, ttlBindTicket); setErr != nil {
			return fmt.Errorf("application: confirm wx scan: set bind ticket: %w", setErr)
		}
		session.Status = ScanStatusScanned
		session.WxOpenID = openID
		session.WxUnionID = unionID
		session.BindTicket = ticket
	}

	if err = s.scanRepo.Confirm(ctx, session); err != nil {
		return fmt.Errorf("application: confirm wx scan: confirm session: %w", err)
	}
	return nil
}

// GetScanStatus polls the QR scan session status.
func (s *authService) GetScanStatus(ctx context.Context, scanToken string) (*ScanStatusResult, error) {
	session, err := s.scanRepo.Get(ctx, scanToken)
	if err != nil {
		return nil, fmt.Errorf("application: get scan status: get session: %w", err)
	}
	if session == nil {
		return nil, apierror.ErrScanNotFound
	}

	result := &ScanStatusResult{Status: session.Status}

	switch session.Status {
	case ScanStatusConfirmed:
		// BindTicket field was used to store the access token (see ConfirmWxScan).
		// We need to issue a full TokenPair here — not practical without the user;
		// instead retrieve userID from session.
		if session.UserID != "" {
			user, findErr := s.userRepo.FindByID(ctx, session.UserID)
			if findErr != nil {
				return nil, fmt.Errorf("application: get scan status: find user: %w", findErr)
			}
			if user == nil {
				return nil, fmt.Errorf("application: get scan status: user not found: %w", apierror.ErrInternal)
			}
			pair, issueErr := s.issueTokenPair(ctx, user)
			if issueErr != nil {
				return nil, fmt.Errorf("application: get scan status: issue tokens: %w", issueErr)
			}
			authCode := uuid.New().String()
			if setErr := s.authCodeRepo.SetAuthCode(ctx, authCode, pair); setErr != nil {
				return nil, fmt.Errorf("application: get scan status: set auth code: %w", setErr)
			}
			result.AuthCode = authCode
		}
	case ScanStatusScanned:
		result.BindTicket = session.BindTicket
		result.WxOpenID = session.WxOpenID
	}

	return result, nil
}

// MiniProgramLogin authenticates a WeChat Mini Program user.
func (s *authService) MiniProgramLogin(ctx context.Context, wxCode, anonID string) (*LoginResult, error) {
	openID, unionID, err := s.wxClient.Code2Session(ctx, wxCode)
	if err != nil {
		return nil, apierror.ErrWxCodeBad
	}

	user, err := s.userRepo.FindByWxOpenID(ctx, openID)
	if err != nil {
		return nil, fmt.Errorf("application: mini program login: find user: %w", err)
	}

	if user != nil {
		// Existing user.
		pair, issueErr := s.issueTokenPair(ctx, user)
		if issueErr != nil {
			return nil, fmt.Errorf("application: mini program login: issue tokens: %w", issueErr)
		}
		s.migrateAnon(user.UserID, anonID)
		return &LoginResult{TokenPair: pair}, nil
	}

	// New WeChat user — issue bind ticket.
	ticket, ticketErr := generateToken()
	if ticketErr != nil {
		return nil, fmt.Errorf("application: mini program login: generate bind ticket: %w", ticketErr)
	}
	if setErr := s.bindRepo.Set(ctx, ticket, openID, unionID, anonID, ttlBindTicket); setErr != nil {
		return nil, fmt.Errorf("application: mini program login: set bind ticket: %w", setErr)
	}

	return &LoginResult{BindTicket: ticket, WxOpenID: openID}, nil
}

// BindPhone binds a phone to a WeChat account identified by a bind ticket.
func (s *authService) BindPhone(ctx context.Context, bindTicket, phone, code, anonID string) (*domain.TokenPair, error) {
	// Retrieve bind ticket.
	openID, unionID, ticketAnonID, err := s.bindRepo.Get(ctx, bindTicket)
	if err != nil {
		return nil, fmt.Errorf("application: bind phone: get bind ticket: %w", err)
	}
	if openID == "" {
		return nil, apierror.ErrBindTicketBad
	}

	// Verify SMS code.
	stored, err := s.smsCodeRepo.GetCode(ctx, phone)
	if err != nil {
		return nil, fmt.Errorf("application: bind phone: get code: %w", err)
	}
	if stored == "" {
		return nil, apierror.ErrSMSCodeExpired
	}
	if stored != code {
		return nil, apierror.ErrSMSCodeWrong
	}
	if err = s.smsCodeRepo.DeleteCode(ctx, phone); err != nil {
		return nil, fmt.Errorf("application: bind phone: delete code: %w", err)
	}

	// Delete bind ticket.
	if err = s.bindRepo.Delete(ctx, bindTicket); err != nil {
		return nil, fmt.Errorf("application: bind phone: delete ticket: %w", err)
	}

	// Find existing phone user or create new.
	user, err := s.userRepo.FindByPhone(ctx, phone)
	if err != nil {
		return nil, fmt.Errorf("application: bind phone: find by phone: %w", err)
	}
	if user == nil {
		user = &domain.User{Phone: phone, WxOpenID: openID, WxUnionID: unionID}
		if err = s.userRepo.Create(ctx, user); err != nil {
			return nil, fmt.Errorf("application: bind phone: create user: %w", err)
		}
	} else {
		// Attach WeChat info to existing phone user.
		if err = s.userRepo.UpdateWxInfo(ctx, user.UserID, openID, unionID); err != nil {
			return nil, fmt.Errorf("application: bind phone: update wx info: %w", err)
		}
	}

	pair, err := s.issueTokenPair(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("application: bind phone: issue tokens: %w", err)
	}

	// Use ticketAnonID if caller's anonID is empty.
	effectiveAnonID := anonID
	if effectiveAnonID == "" {
		effectiveAnonID = ticketAnonID
	}
	s.migrateAnon(user.UserID, effectiveAnonID)

	return pair, nil
}

// RefreshToken validates a refresh token and issues a new access token.
// The refresh token itself is not rotated.
func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*domain.TokenPair, error) {
	// Parse claims from the refresh token to extract userID, uid and role.
	parsedUserID, uid, role, parseErr := s.signer.Parse(refreshToken)
	if parseErr != nil {
		return nil, apierror.ErrTokenExpired
	}

	// Verify that the token is still stored in Redis (revocation check).
	stored, err := s.tokenRepo.GetRefresh(ctx, parsedUserID)
	if err != nil {
		return nil, fmt.Errorf("application: refresh token: get refresh: %w", err)
	}
	if stored == "" || stored != refreshToken {
		return nil, apierror.ErrTokenExpired
	}

	accessToken, signErr := s.signer.Sign(parsedUserID, uid, role)
	if signErr != nil {
		return nil, fmt.Errorf("application: refresh token: sign access token: %w", signErr)
	}

	return &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.refreshTTL,
	}, nil
}

// Logout removes the refresh token for the user.
func (s *authService) Logout(ctx context.Context, userID string) error {
	if err := s.tokenRepo.DeleteRefresh(ctx, userID); err != nil {
		return fmt.Errorf("application: logout: delete refresh: %w", err)
	}
	return nil
}

// ExchangeAuthCode exchanges a short-lived auth code for a TokenPair.
// The code is consumed (deleted) atomically on first use.
func (s *authService) ExchangeAuthCode(ctx context.Context, authCode string) (*domain.TokenPair, error) {
	pair, err := s.authCodeRepo.ConsumeAuthCode(ctx, authCode)
	if err != nil {
		return nil, fmt.Errorf("application: exchange auth code: consume: %w", err)
	}
	if pair == nil {
		return nil, apierror.ErrAuthCodeExpired
	}
	return pair, nil
}

// AdminStaticLogin validates configured static credentials and issues JWT for UserRoleAdmin.
func (s *authService) AdminStaticLogin(ctx context.Context, username, password string) (*domain.TokenPair, error) {
	if !s.adminStatic.active() {
		return nil, apierror.ErrForbidden
	}
	if username != s.adminStatic.Username || password != s.adminStatic.Password {
		return nil, apierror.ErrAdminLoginInvalid
	}
	if err := s.userRepo.UpsertStaticAdmin(ctx, s.adminStatic.UserID, s.adminStatic.UID, domain.UserRoleAdmin); err != nil {
		return nil, fmt.Errorf("application: admin static login: upsert admin user: %w", err)
	}
	user, err := s.userRepo.FindByID(ctx, s.adminStatic.UserID)
	if err != nil {
		return nil, fmt.Errorf("application: admin static login: find admin: %w", err)
	}
	if user == nil {
		return nil, apierror.ErrInternal
	}
	pair, err := s.issueTokenPair(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("application: admin static login: issue tokens: %w", err)
	}
	return pair, nil
}

// issueTokenPair signs an access token and stores a refresh token.
func (s *authService) issueTokenPair(ctx context.Context, user *domain.User) (*domain.TokenPair, error) {
	accessToken, err := s.signer.Sign(user.UserID, user.UID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("application: issue token pair: sign access: %w", err)
	}

	refreshTTL := time.Duration(s.refreshTTL) * time.Second
	refreshToken, err := s.signer.SignWithTTL(user.UserID, user.UID, user.Role, refreshTTL)
	if err != nil {
		return nil, fmt.Errorf("application: issue token pair: sign refresh: %w", err)
	}

	if err = s.tokenRepo.SetRefresh(ctx, user.UserID, refreshToken, s.refreshTTL); err != nil {
		return nil, fmt.Errorf("application: issue token pair: store refresh: %w", err)
	}

	return &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.refreshTTL,
	}, nil
}

// migrateAnon asynchronously migrates anonymous sessions to the registered user.
func (s *authService) migrateAnon(userID, anonID string) {
	if anonID == "" {
		return
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := s.userRepo.MigrateAnonSessions(ctx, anonID, userID); err != nil {
			_ = err // best-effort; log in production
		}
	}()
}

// generate6DigitCode generates a cryptographically random 6-digit numeric string.
func generate6DigitCode() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", fmt.Errorf("generate6DigitCode: %w", err)
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// generateToken generates a random 32-byte hex string for use as a scan token
// or bind ticket.
func generateToken() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generateToken: %w", err)
	}
	return fmt.Sprintf("%x", b), nil
}
