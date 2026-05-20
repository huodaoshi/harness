package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/huodaoshi/harness/backend/modules/auth/domain"
	"github.com/huodaoshi/harness/backend/pkg/idgen"
)

// ---------------------------------------------------------------------------
// MongoUserRepo
// ---------------------------------------------------------------------------

// MongoUserRepo implements domain.UserRepo backed by MongoDB.
type MongoUserRepo struct {
	coll    *mongo.Collection
	counter *idgen.MongoCounter
	idg     *idgen.IDGenerator
}

// NewMongoUserRepo creates a MongoUserRepo.
func NewMongoUserRepo(db *mongo.Database, counter *idgen.MongoCounter, idg *idgen.IDGenerator) domain.UserRepo {
	return &MongoUserRepo{
		coll:    db.Collection("users"),
		counter: counter,
		idg:     idg,
	}
}

// Create inserts a new user. It generates a UserID via snowflake and a UID via
// the mongo counter before inserting.
func (r *MongoUserRepo) Create(ctx context.Context, u *domain.User) error {
	u.UserID = r.idg.Generate()

	uid, err := r.counter.Next(ctx, idgen.CounterUID)
	if err != nil {
		return fmt.Errorf("infra: user: create: get uid: %w", err)
	}
	u.UID = uid

	now := time.Now()
	u.CreatedAt = now
	u.UpdatedAt = now

	if _, err = r.coll.InsertOne(ctx, u); err != nil {
		return fmt.Errorf("infra: user: create: insert: %w", err)
	}
	return nil
}

// FindByPhone looks up a user by phone number.
func (r *MongoUserRepo) FindByPhone(ctx context.Context, phone string) (*domain.User, error) {
	var u domain.User
	err := r.coll.FindOne(ctx, bson.M{"phone": phone}).Decode(&u)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("infra: user: find by phone: %w", err)
	}
	return &u, nil
}

// FindByID looks up a user by UserID.
func (r *MongoUserRepo) FindByID(ctx context.Context, userID string) (*domain.User, error) {
	var u domain.User
	err := r.coll.FindOne(ctx, bson.M{"user_id": userID}).Decode(&u)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("infra: user: find by id: %w", err)
	}
	return &u, nil
}

// FindByWxOpenID looks up a user by WeChat OpenID.
func (r *MongoUserRepo) FindByWxOpenID(ctx context.Context, openID string) (*domain.User, error) {
	var u domain.User
	err := r.coll.FindOne(ctx, bson.M{"wx_openid": openID}).Decode(&u)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("infra: user: find by wx_openid: %w", err)
	}
	return &u, nil
}

// UpdateWxInfo sets wx_openid and wx_unionid on an existing user.
func (r *MongoUserRepo) UpdateWxInfo(ctx context.Context, userID, openID, unionID string) error {
	filter := bson.M{"user_id": userID}
	update := bson.M{
		"$set": bson.M{
			"wx_openid":  openID,
			"wx_unionid": unionID,
			"updated_at": time.Now(),
		},
	}
	if _, err := r.coll.UpdateOne(ctx, filter, update); err != nil {
		return fmt.Errorf("infra: user: update wx info: %w", err)
	}
	return nil
}

// MigrateAnonSessions reassigns all sessions belonging to anonID to userID.
func (r *MongoUserRepo) MigrateAnonSessions(ctx context.Context, anonID, userID string) error {
	db := r.coll.Database()
	sessions := db.Collection("sessions")
	filter := bson.M{"anon_id": anonID}
	update := bson.M{"$set": bson.M{"user_id": userID, "anon_id": ""}}
	if _, err := sessions.UpdateMany(ctx, filter, update); err != nil {
		return fmt.Errorf("infra: user: migrate anon sessions: %w", err)
	}
	return nil
}

// UpsertStaticAdmin inserts or updates the fixed admin user (does not use snowflake UserID).
func (r *MongoUserRepo) UpsertStaticAdmin(ctx context.Context, userID string, uid int64, role domain.UserRole) error {
	now := time.Now()
	filter := bson.M{"user_id": userID}
	update := bson.M{
		"$set": bson.M{
			"uid":        uid,
			"role":       role.Int(),
			"updated_at": now,
		},
		"$setOnInsert": bson.M{
			"user_id":    userID,
			"phone":      "",
			"wx_openid":  "",
			"wx_unionid": "",
			"anon_id":    "",
			"created_at": now,
		},
	}
	opts := options.Update().SetUpsert(true)
	if _, err := r.coll.UpdateOne(ctx, filter, update, opts); err != nil {
		return fmt.Errorf("infra: user: upsert static admin: %w", err)
	}
	return nil
}

// ---------------------------------------------------------------------------
// RedisAuthRepo — implements SMSCodeRepo, TokenRepo, ScanRepo, BindTicketRepo
// ---------------------------------------------------------------------------

// RedisAuthRepo implements all Redis-backed auth repositories.
type RedisAuthRepo struct {
	rdb redis.UniversalClient
}

// NewRedisAuthRepo creates a RedisAuthRepo.
func NewRedisAuthRepo(rdb redis.UniversalClient) *RedisAuthRepo {
	return &RedisAuthRepo{rdb: rdb}
}

// Redis key helpers.
func keySMSCode(phone string) string        { return fmt.Sprintf("auth:sms:code:%s", phone) }
func keySMSAttempts(phone string) string    { return fmt.Sprintf("auth:sms:attempts:%s", phone) }
func keySMSSendLimit(phone string) string   { return fmt.Sprintf("auth:sms:limit:%s", phone) }
func keySMSDailyCount(phone string) string  { return fmt.Sprintf("auth:sms:daily:%s", phone) }
func keySMSIPCount(ip string) string        { return fmt.Sprintf("auth:sms:ip:%s", ip) }
func keySMSLock(phone string) string        { return fmt.Sprintf("auth:sms:lock:%s", phone) }
func keyRefreshToken(userID string) string  { return fmt.Sprintf("auth:refresh:%s", userID) }
func keyScanSession(scanToken string) string { return fmt.Sprintf("auth:scan:%s", scanToken) }
func keyBindTicket(ticket string) string    { return fmt.Sprintf("auth:bind:%s", ticket) }

// ---------------------------------------------------------------------------
// SMSCodeRepo implementation
// ---------------------------------------------------------------------------

// SetCode stores a verification code for the given phone.
func (r *RedisAuthRepo) SetCode(ctx context.Context, phone, code string, ttlSeconds int) error {
	if err := r.rdb.Set(ctx, keySMSCode(phone), code, time.Duration(ttlSeconds)*time.Second).Err(); err != nil {
		return fmt.Errorf("infra: sms: set code: %w", err)
	}
	return nil
}

// GetCode retrieves the current verification code for a phone.
func (r *RedisAuthRepo) GetCode(ctx context.Context, phone string) (string, error) {
	val, err := r.rdb.Get(ctx, keySMSCode(phone)).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", fmt.Errorf("infra: sms: get code: %w", err)
	}
	return val, nil
}

// IncrAttempts increments and returns the verification attempt counter.
func (r *RedisAuthRepo) IncrAttempts(ctx context.Context, phone string) (int64, error) {
	val, err := r.rdb.Incr(ctx, keySMSAttempts(phone)).Result()
	if err != nil {
		return 0, fmt.Errorf("infra: sms: incr attempts: %w", err)
	}
	return val, nil
}

// SetSendLimit sets a cooldown key to rate-limit SMS sends per phone.
func (r *RedisAuthRepo) SetSendLimit(ctx context.Context, phone string, ttlSeconds int) error {
	if err := r.rdb.Set(ctx, keySMSSendLimit(phone), 1, time.Duration(ttlSeconds)*time.Second).Err(); err != nil {
		return fmt.Errorf("infra: sms: set send limit: %w", err)
	}
	return nil
}

// CheckSendLimit returns true if the phone is currently rate-limited.
func (r *RedisAuthRepo) CheckSendLimit(ctx context.Context, phone string) (bool, error) {
	exists, err := r.rdb.Exists(ctx, keySMSSendLimit(phone)).Result()
	if err != nil {
		return false, fmt.Errorf("infra: sms: check send limit: %w", err)
	}
	return exists > 0, nil
}

// IncrDailyCount increments the daily SMS send count. On first increment the
// key is set to expire at 23:59:59 of the current day.
func (r *RedisAuthRepo) IncrDailyCount(ctx context.Context, phone string) (int64, error) {
	key := keySMSDailyCount(phone)
	val, err := r.rdb.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("infra: sms: incr daily count: %w", err)
	}

	if val == 1 {
		// First increment — set TTL to end of today.
		now := time.Now()
		eod := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
		if expErr := r.rdb.ExpireAt(ctx, key, eod).Err(); expErr != nil {
			return val, fmt.Errorf("infra: sms: set daily expiry: %w", expErr)
		}
	}
	return val, nil
}

// IncrIPCount increments the per-IP SMS send counter with a rolling TTL.
func (r *RedisAuthRepo) IncrIPCount(ctx context.Context, ip string, ttlSeconds int) (int64, error) {
	key := keySMSIPCount(ip)
	pipe := r.rdb.Pipeline()
	incrCmd := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, time.Duration(ttlSeconds)*time.Second)
	if _, err := pipe.Exec(ctx); err != nil {
		return 0, fmt.Errorf("infra: sms: incr ip count: %w", err)
	}
	return incrCmd.Val(), nil
}

// SetLock locks a phone after too many failed attempts.
func (r *RedisAuthRepo) SetLock(ctx context.Context, phone string, ttlSeconds int) error {
	if err := r.rdb.Set(ctx, keySMSLock(phone), 1, time.Duration(ttlSeconds)*time.Second).Err(); err != nil {
		return fmt.Errorf("infra: sms: set lock: %w", err)
	}
	return nil
}

// CheckLock returns true if the phone is currently locked.
func (r *RedisAuthRepo) CheckLock(ctx context.Context, phone string) (bool, error) {
	exists, err := r.rdb.Exists(ctx, keySMSLock(phone)).Result()
	if err != nil {
		return false, fmt.Errorf("infra: sms: check lock: %w", err)
	}
	return exists > 0, nil
}

// DeleteCode removes the SMS verification code for a phone.
func (r *RedisAuthRepo) DeleteCode(ctx context.Context, phone string) error {
	if err := r.rdb.Del(ctx, keySMSCode(phone)).Err(); err != nil {
		return fmt.Errorf("infra: sms: delete code: %w", err)
	}
	return nil
}

// ---------------------------------------------------------------------------
// TokenRepo implementation
// ---------------------------------------------------------------------------

// SetRefresh stores a refresh token for the given userID.
func (r *RedisAuthRepo) SetRefresh(ctx context.Context, userID, token string, ttlSeconds int) error {
	if err := r.rdb.Set(ctx, keyRefreshToken(userID), token, time.Duration(ttlSeconds)*time.Second).Err(); err != nil {
		return fmt.Errorf("infra: token: set refresh: %w", err)
	}
	return nil
}

// GetRefresh retrieves the stored refresh token for a userID.
func (r *RedisAuthRepo) GetRefresh(ctx context.Context, userID string) (string, error) {
	val, err := r.rdb.Get(ctx, keyRefreshToken(userID)).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", fmt.Errorf("infra: token: get refresh: %w", err)
	}
	return val, nil
}

// DeleteRefresh removes the refresh token for a userID.
func (r *RedisAuthRepo) DeleteRefresh(ctx context.Context, userID string) error {
	if err := r.rdb.Del(ctx, keyRefreshToken(userID)).Err(); err != nil {
		return fmt.Errorf("infra: token: delete refresh: %w", err)
	}
	return nil
}

// ---------------------------------------------------------------------------
// ScanRepo implementation
// ---------------------------------------------------------------------------

// scanFields defines the Redis hash field names for a ScanSession.
const (
	scanFieldStatus     = "status"
	scanFieldAnonID     = "anon_id"
	scanFieldUserID     = "user_id"
	scanFieldWxOpenID   = "wx_openid"
	scanFieldWxUnionID  = "wx_unionid"
	scanFieldBindTicket = "bind_ticket"
	scanFieldScanToken  = "scan_token"
)

// Set stores a ScanSession as a Redis hash with a TTL.
func (r *RedisAuthRepo) Set(ctx context.Context, session *domain.ScanSession, ttlSeconds int) error {
	key := keyScanSession(session.ScanToken)
	fields := map[string]any{
		scanFieldScanToken:  session.ScanToken,
		scanFieldStatus:     session.Status,
		scanFieldAnonID:     session.AnonID,
		scanFieldUserID:     session.UserID,
		scanFieldWxOpenID:   session.WxOpenID,
		scanFieldWxUnionID:  session.WxUnionID,
		scanFieldBindTicket: session.BindTicket,
	}
	pipe := r.rdb.Pipeline()
	pipe.HSet(ctx, key, fields)
	pipe.Expire(ctx, key, time.Duration(ttlSeconds)*time.Second)
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("infra: scan: set: %w", err)
	}
	return nil
}

// Get retrieves a ScanSession by scan token.
func (r *RedisAuthRepo) Get(ctx context.Context, scanToken string) (*domain.ScanSession, error) {
	key := keyScanSession(scanToken)
	vals, err := r.rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("infra: scan: get: %w", err)
	}
	if len(vals) == 0 {
		return nil, nil
	}

	var status int
	if _, scanErr := fmt.Sscanf(vals[scanFieldStatus], "%d", &status); scanErr != nil {
		return nil, fmt.Errorf("infra: scan: get: parse status: %w", scanErr)
	}

	return &domain.ScanSession{
		ScanToken:  vals[scanFieldScanToken],
		Status:     status,
		AnonID:     vals[scanFieldAnonID],
		UserID:     vals[scanFieldUserID],
		WxOpenID:   vals[scanFieldWxOpenID],
		WxUnionID:  vals[scanFieldWxUnionID],
		BindTicket: vals[scanFieldBindTicket],
	}, nil
}

// Confirm overwrites a ScanSession, preserving the existing TTL (KEEPTTL).
func (r *RedisAuthRepo) Confirm(ctx context.Context, session *domain.ScanSession) error {
	key := keyScanSession(session.ScanToken)
	fields := map[string]any{
		scanFieldScanToken:  session.ScanToken,
		scanFieldStatus:     session.Status,
		scanFieldAnonID:     session.AnonID,
		scanFieldUserID:     session.UserID,
		scanFieldWxOpenID:   session.WxOpenID,
		scanFieldWxUnionID:  session.WxUnionID,
		scanFieldBindTicket: session.BindTicket,
	}
	// Use a pipeline: delete then re-set with KEEPTTL is not directly supported for
	// hash via a single command; instead we use HSet (which overwrites) and
	// rely on the existing TTL being preserved by not calling EXPIRE again.
	if err := r.rdb.HSet(ctx, key, fields).Err(); err != nil {
		return fmt.Errorf("infra: scan: confirm: %w", err)
	}
	return nil
}

// ---------------------------------------------------------------------------
// BindTicketRepo implementation
// ---------------------------------------------------------------------------

// Set stores bind ticket data.
func (r *RedisAuthRepo) SetBindTicket(ctx context.Context, ticket, openID, unionID, anonID string, ttlSeconds int) error {
	key := keyBindTicket(ticket)
	fields := map[string]any{
		"openid":  openID,
		"unionid": unionID,
		"anon_id": anonID,
	}
	pipe := r.rdb.Pipeline()
	pipe.HSet(ctx, key, fields)
	pipe.Expire(ctx, key, time.Duration(ttlSeconds)*time.Second)
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("infra: bind ticket: set: %w", err)
	}
	return nil
}

// GetBindTicket retrieves bind ticket data.
func (r *RedisAuthRepo) GetBindTicket(ctx context.Context, ticket string) (openID, unionID, anonID string, err error) {
	key := keyBindTicket(ticket)
	vals, getErr := r.rdb.HGetAll(ctx, key).Result()
	if getErr != nil {
		return "", "", "", fmt.Errorf("infra: bind ticket: get: %w", getErr)
	}
	if len(vals) == 0 {
		return "", "", "", nil
	}
	return vals["openid"], vals["unionid"], vals["anon_id"], nil
}

// DeleteBindTicket removes a bind ticket.
func (r *RedisAuthRepo) DeleteBindTicket(ctx context.Context, ticket string) error {
	if err := r.rdb.Del(ctx, keyBindTicket(ticket)).Err(); err != nil {
		return fmt.Errorf("infra: bind ticket: delete: %w", err)
	}
	return nil
}

// ---------------------------------------------------------------------------
// AuthCodeRepo implementation
// ---------------------------------------------------------------------------

// keyAuthCode returns the Redis key for a short-lived auth code.
func keyAuthCode(authCode string) string { return fmt.Sprintf("auth:scan:authcode:%s", authCode) }

// consumeAuthCodeScript atomically reads and deletes an auth code key.
var consumeAuthCodeScript = redis.NewScript(`
local v = redis.call('GET', KEYS[1])
if v then redis.call('DEL', KEYS[1]) end
return v
`)

// SetAuthCode stores a TokenPair as JSON under the auth code key with a 5-second TTL.
func (r *RedisAuthRepo) SetAuthCode(ctx context.Context, authCode string, pair *domain.TokenPair) error {
	jsonBytes, err := json.Marshal(pair)
	if err != nil {
		return fmt.Errorf("infra: auth code: marshal: %w", err)
	}
	if err = r.rdb.Set(ctx, keyAuthCode(authCode), jsonBytes, 5*time.Second).Err(); err != nil {
		return fmt.Errorf("infra: auth code: set: %w", err)
	}
	return nil
}

// ConsumeAuthCode atomically reads and deletes the auth code, returning the
// stored TokenPair. Returns nil, nil if the code does not exist.
func (r *RedisAuthRepo) ConsumeAuthCode(ctx context.Context, authCode string) (*domain.TokenPair, error) {
	res, err := consumeAuthCodeScript.Run(ctx, r.rdb, []string{keyAuthCode(authCode)}).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("infra: auth code: consume: %w", err)
	}
	if res == nil {
		return nil, nil
	}
	raw, ok := res.(string)
	if !ok || raw == "" {
		return nil, nil
	}
	var pair domain.TokenPair
	if err = json.Unmarshal([]byte(raw), &pair); err != nil {
		return nil, fmt.Errorf("infra: auth code: unmarshal: %w", err)
	}
	return &pair, nil
}

// Ensure RedisAuthRepo satisfies the domain repository interfaces.
var (
	_ domain.SMSCodeRepo    = (*RedisAuthRepo)(nil)
	_ domain.TokenRepo      = (*RedisAuthRepo)(nil)
	_ domain.ScanRepo       = (*RedisAuthRepo)(nil)
	_ domain.BindTicketRepo = (*bindTicketRepoAdapter)(nil)
	_ domain.AuthCodeRepo   = (*RedisAuthRepo)(nil)
)

// bindTicketRepoAdapter wraps RedisAuthRepo to satisfy domain.BindTicketRepo.
type bindTicketRepoAdapter struct {
	r *RedisAuthRepo
}

// NewBindTicketRepo wraps a RedisAuthRepo as a domain.BindTicketRepo.
func NewBindTicketRepo(r *RedisAuthRepo) domain.BindTicketRepo {
	return &bindTicketRepoAdapter{r: r}
}

func (a *bindTicketRepoAdapter) Set(ctx context.Context, ticket, openID, unionID, anonID string, ttlSeconds int) error {
	return a.r.SetBindTicket(ctx, ticket, openID, unionID, anonID, ttlSeconds)
}

func (a *bindTicketRepoAdapter) Get(ctx context.Context, ticket string) (openID, unionID, anonID string, err error) {
	return a.r.GetBindTicket(ctx, ticket)
}

func (a *bindTicketRepoAdapter) Delete(ctx context.Context, ticket string) error {
	return a.r.DeleteBindTicket(ctx, ticket)
}