package main

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/huodaoshi/harness/backend/conf"
	"github.com/huodaoshi/harness/backend/infra"
	authapp "github.com/huodaoshi/harness/backend/modules/auth/application"
	authinf "github.com/huodaoshi/harness/backend/modules/auth/infra"
	authpkg "github.com/huodaoshi/harness/backend/modules/auth/pkg"
	"github.com/huodaoshi/harness/backend/pkg/idgen"
)

type authBundle struct {
	Service authapp.AuthService
	Signer  authpkg.JWTSigner
}

func wireAuth(cfg *conf.Config, mongoClient *mongo.Client, redisClient redis.UniversalClient) (*authBundle, error) {
	if mongoClient == nil {
		return nil, fmt.Errorf("auth: mongodb client required")
	}
	if redisClient == nil {
		return nil, fmt.Errorf("auth: redis client required")
	}

	mongoDB := mongoClient.Database(cfg.MongoDB.Database)
	counter := idgen.NewMongoCounter(mongoDB)
	idGenerator, err := idgen.NewIDGenerator(1)
	if err != nil {
		return nil, err
	}

	redisAuthRepo := authinf.NewRedisAuthRepo(redisClient)
	userRepo := authinf.NewMongoUserRepo(mongoDB, counter, idGenerator)
	bindRepo := authinf.NewBindTicketRepo(redisAuthRepo)
	wxClient := authinf.NewHTTPWeChatClient(cfg.WeChat.MiniProgram.AppID, cfg.WeChat.MiniProgram.AppSecret)
	jwtSigner := authpkg.NewHS256Signer(cfg.JWT.Secret, cfg.JWT.AccessTokenTTL)

	smsService, err := infra.NewSMSService(cfg.SMS)
	if err != nil {
		return nil, err
	}

	svc := authapp.NewAuthService(
		userRepo,
		redisAuthRepo,
		redisAuthRepo,
		redisAuthRepo,
		bindRepo,
		wxClient,
		jwtSigner,
		smsService,
		cfg.JWT.RefreshTokenTTL,
		redisAuthRepo,
		authapp.AdminStaticLoginOpts{
			Enabled:  cfg.AdminStaticLogin.Enabled,
			Username: cfg.AdminStaticLogin.Username,
			Password: cfg.AdminStaticLogin.Password,
			UserID:   cfg.AdminStaticLogin.UserID,
			UID:      cfg.AdminStaticLogin.UID,
		},
	)

	return &authBundle{Service: svc, Signer: jwtSigner}, nil
}
