package conf

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

func applyConnectionEnv(cfg *Config) error {
	cfg.MongoDB.Database = strings.TrimSpace(os.ExpandEnv(cfg.MongoDB.Database))
	if d := strings.TrimSpace(os.Getenv("MONGODB_DB")); d != "" {
		cfg.MongoDB.Database = d
	}
	cfg.MongoDB.URI = strings.TrimSpace(os.ExpandEnv(cfg.MongoDB.URI))

	if u := strings.TrimSpace(os.Getenv("MONGODB_URI")); u != "" {
		cfg.MongoDB.URI = u
	} else if user := strings.TrimSpace(os.Getenv("MONGODB_USER")); user != "" {
		pass := os.Getenv("MONGODB_PASSWORD")
		host := getenvDefault("MONGODB_HOST", "127.0.0.1:27017")
		db := cfg.MongoDB.Database
		if db == "" {
			db = "family_wellness"
		}
		authSource := getenvDefault("MONGODB_AUTH_SOURCE", db)
		cfg.MongoDB.URI = buildMongoURI(user, pass, host, db, authSource)
	} else if cfg.MongoDB.URI == "" {
		cfg.MongoDB.URI = "mongodb://localhost:27017"
	}

	if addr := strings.TrimSpace(os.Getenv("REDIS_ADDR")); addr != "" {
		cfg.Redis.Addrs = []string{addr}
	}
	for i := range cfg.Redis.Addrs {
		cfg.Redis.Addrs[i] = strings.TrimSpace(os.ExpandEnv(cfg.Redis.Addrs[i]))
	}
	if _, ok := os.LookupEnv("REDIS_PASSWORD"); ok {
		cfg.Redis.Password = os.Getenv("REDIS_PASSWORD")
	} else {
		cfg.Redis.Password = os.ExpandEnv(cfg.Redis.Password)
	}

	if len(cfg.Redis.Addrs) == 0 && cfg.Redis.Required {
		return fmt.Errorf("conf: redis addrs empty (set redis.addrs or REDIS_ADDR)")
	}

	cfg.LLM.APIKey = strings.TrimSpace(os.ExpandEnv(cfg.LLM.APIKey))
	if v := strings.TrimSpace(os.Getenv("ARK_API_KEY")); v != "" {
		cfg.LLM.APIKey = v
	}
	if v := strings.TrimSpace(os.Getenv("LLM_PROVIDER")); v != "" {
		cfg.LLM.Provider = v
	}
	if v := strings.TrimSpace(os.Getenv("ARK_MODEL_ID")); v != "" {
		cfg.LLM.Model = v
	} else 	if v := strings.TrimSpace(os.Getenv("LLM_MODEL_ID")); v != "" {
		cfg.LLM.Model = v
	}
	cfg.JWT.Secret = strings.TrimSpace(os.ExpandEnv(cfg.JWT.Secret))
	if v := strings.TrimSpace(os.Getenv("JWT_SECRET")); v != "" {
		cfg.JWT.Secret = v
	}
	if cfg.JWT.Secret == "" {
		return fmt.Errorf("conf: jwt secret is empty (set jwt.secret or JWT_SECRET)")
	}
	cfg.AdminStaticLogin.Password = strings.TrimSpace(os.ExpandEnv(cfg.AdminStaticLogin.Password))
	if v := strings.TrimSpace(os.Getenv("ADMIN_STATIC_PASSWORD")); v != "" {
		cfg.AdminStaticLogin.Password = v
	}
	return nil
}

func getenvDefault(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}

func buildMongoURI(user, pass, host, db, authSource string) string {
	u := &url.URL{
		Scheme: "mongodb",
		Host:   host,
		Path:   "/" + db,
	}
	if user != "" || pass != "" {
		u.User = url.UserPassword(user, pass)
	}
	q := url.Values{}
	if authSource != "" {
		q.Set("authSource", authSource)
	}
	u.RawQuery = q.Encode()
	return u.String()
}
