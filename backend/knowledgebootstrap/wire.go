package knowledgebootstrap

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/cloudwego/eino/components/embedding"
	arkembed "github.com/cloudwego/eino-ext/components/embedding/ark"
	redisCli "github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/huodaoshi/harness/backend/conf"
	"github.com/huodaoshi/harness/backend/infra"
	knowledgeapp "github.com/huodaoshi/harness/backend/modules/knowledge/application"
	"github.com/huodaoshi/harness/backend/modules/knowledge/domain"
	knowledgeinf "github.com/huodaoshi/harness/backend/modules/knowledge/infra"
	"github.com/huodaoshi/harness/backend/modules/knowledgeindexing"
	"github.com/huodaoshi/harness/backend/pkg/idgen"
)

// Bundle holds wired knowledge services and optional embedded worker cancel.
type Bundle struct {
	SpaceAdmin knowledgeapp.SpaceAdminService
	Ingest     knowledgeapp.IngestService
	Knowledge  knowledgeapp.KnowledgeService
	ObjStorage infra.ObjectStorage
	EmbedDim   int
	cancel     context.CancelFunc
}

// Close stops background workers (outbox flusher, embedded MQ consumer).
func (b *Bundle) Close() {
	if b.cancel != nil {
		b.cancel()
	}
}

// Wire initializes knowledge repos, default space, MQ/outbox, and optional embedded worker.
func Wire(ctx context.Context, cfg *conf.Config, mongoClient *mongo.Client, redisClient redisCli.UniversalClient, embedWorker bool) (*Bundle, error) {
	if mongoClient == nil || redisClient == nil {
		return nil, fmt.Errorf("knowledgebootstrap: mongo and redis required")
	}

	mongoDB := mongoClient.Database(cfg.MongoDB.Database)
	if err := knowledgeinf.EnsureKnowledgeIndexes(ctx, mongoDB); err != nil {
		return nil, fmt.Errorf("knowledge indexes: %w", err)
	}
	if err := knowledgeinf.EnsureIngestJobDocKeyIndex(ctx, mongoDB); err != nil {
		return nil, fmt.Errorf("knowledge doc_key index: %w", err)
	}

	if _, err := knowledgeinf.EnsureDefaultSpace(ctx, mongoDB); err != nil {
		return nil, fmt.Errorf("default space: %w", err)
	}

	counter := idgen.NewMongoCounter(mongoDB)
	idGenerator, err := idgen.NewIDGenerator(2)
	if err != nil {
		return nil, err
	}

	embedder, dim, err := newEmbedder(ctx, cfg)
	if err != nil {
		return nil, err
	}
	if err := knowledgeinf.EnsureKnowledgeVectorIndex(ctx, redisClient, knowledgeinf.DefaultSpaceID, dim); err != nil {
		return nil, fmt.Errorf("vector index: %w", err)
	}

	mq, err := infra.NewMessageQueue(cfg.MQ)
	if err != nil {
		return nil, fmt.Errorf("mq: %w", err)
	}
	objStorage, err := infra.NewObjectStorage(cfg.COS)
	if err != nil {
		return nil, fmt.Errorf("object storage: %w", err)
	}

	spaceRepo := knowledgeinf.NewMongoSpaceRepo(mongoDB)
	vectorRepo := knowledgeinf.NewRedisVectorRepo(redisClient)
	knowledgeSvc := knowledgeapp.NewKnowledgeService(spaceRepo, vectorRepo)

	spaceAdminRepo := knowledgeinf.NewMongoSpaceAdminRepo(mongoDB, counter, idGenerator)
	ingestJobRepo := knowledgeinf.NewMongoIngestJobRepo(mongoDB, idGenerator)
	mqOutboxRepo := knowledgeinf.NewMongoMQOutboxRepo(mongoDB, idGenerator)
	spaceCleanupRepo := knowledgeinf.NewSpaceCleanupRepo(redisClient, mongoDB)
	ingestChunkRepo := knowledgeinf.NewRedisIngestChunkRepo(redisClient)

	spaceAdminSvc := knowledgeapp.NewSpaceAdminService(spaceAdminRepo, spaceCleanupRepo, spaceAdminRepo)
	ingestSvc := knowledgeapp.NewIngestService(spaceRepo, ingestJobRepo, mq, mqOutboxRepo, ingestChunkRepo, objStorage)

	workerCtx, cancel := context.WithCancel(ctx)
	go knowledgeapp.StartMQOutboxFlusher(workerCtx, mqOutboxRepo, ingestJobRepo, mq)

	b := &Bundle{
		SpaceAdmin: spaceAdminSvc,
		Ingest:     ingestSvc,
		Knowledge:  knowledgeSvc,
		ObjStorage: objStorage,
		EmbedDim:   dim,
		cancel:     cancel,
	}

	if embedWorker && cfg.MQ.Provider == "local" {
		if err := startEmbeddedWorker(workerCtx, cfg, redisClient, ingestJobRepo, spaceRepo, objStorage, embedder, mq); err != nil {
			cancel()
			return nil, err
		}
		log.Printf("knowledge: embedded worker subscribed (mq.provider=local)")
	} else if cfg.MQ.Provider == "local" && !embedWorker {
		slog.Warn("knowledge: mq.provider=local without embedded worker; run cmd/knowledgeindexing in same machine or switch mq.provider")
	}

	return b, nil
}

// StartStandaloneWorker subscribes to knowledge-ingest (for cmd/knowledgeindexing).
// Use with mq.provider=rocketmq in production; local MQ only works in-process with Wire(embedWorker=true).
func StartStandaloneWorker(ctx context.Context, cfg *conf.Config, mongoClient *mongo.Client, redisClient redisCli.UniversalClient) (context.CancelFunc, error) {
	mongoDB := mongoClient.Database(cfg.MongoDB.Database)
	idGenerator, err := idgen.NewIDGenerator(3)
	if err != nil {
		return nil, err
	}
	embedder, _, err := newEmbedder(ctx, cfg)
	if err != nil {
		return nil, err
	}
	mq, err := infra.NewMessageQueue(cfg.MQ)
	if err != nil {
		return nil, err
	}
	objStorage, err := infra.NewObjectStorage(cfg.COS)
	if err != nil {
		return nil, err
	}
	ingestJobRepo := knowledgeinf.NewMongoIngestJobRepo(mongoDB, idGenerator)
	spaceRepo := knowledgeinf.NewMongoSpaceRepo(mongoDB)
	workerCtx, cancel := context.WithCancel(ctx)
	if err := startEmbeddedWorker(workerCtx, cfg, redisClient, ingestJobRepo, spaceRepo, objStorage, embedder, mq); err != nil {
		cancel()
		return nil, err
	}
	return cancel, nil
}

func newEmbedder(ctx context.Context, cfg *conf.Config) (embedding.Embedder, int, error) {
	switch cfg.Embedding.Provider {
	case "fake", "":
		return infra.NewFakeEmbedder(), infra.FakeEmbedDim(), nil
	case "ark":
		emb, err := arkembed.NewEmbedder(ctx, &arkembed.EmbeddingConfig{
			APIKey: cfg.Embedding.APIKey,
			Model:  cfg.Embedding.Model,
		})
		if err != nil {
			return nil, 0, err
		}
		dim := cfg.Embedding.Dim
		if dim <= 0 {
			dim = 1024
		}
		return emb, dim, nil
	default:
		return nil, 0, fmt.Errorf("unknown embedding provider %q", cfg.Embedding.Provider)
	}
}

func startEmbeddedWorker(
	ctx context.Context,
	cfg *conf.Config,
	redisClient redisCli.UniversalClient,
	ingestJobRepo domain.IngestJobRepo,
	spaceRepo domain.SpaceRepo,
	objStorage infra.ObjectStorage,
	embedder embedding.Embedder,
	mq infra.MessageQueue,
) error {
	addr := ""
	if len(cfg.Redis.Addrs) > 0 {
		addr = cfg.Redis.Addrs[0]
	}
	indexerRedis := redisCli.NewClient(&redisCli.Options{
		Addr:     addr,
		Password: cfg.Redis.Password,
		Protocol: 2,
	})

	indexer, err := knowledgeindexing.BuildRedisIndexer(ctx, indexerRedis, embedder)
	if err != nil {
		return err
	}
	pipeline, err := knowledgeindexing.BuildIndexingPipeline(ctx, indexer)
	if err != nil {
		return err
	}
	parsers, err := knowledgeindexing.NewParsers(ctx)
	if err != nil {
		return err
	}
	fetchOpts := knowledgeindexing.FetchOptions{
		AllowHosts:                 cfg.KnowledgeIndexing.IngestFetchAllowHosts,
		AllowedContentTypePrefixes: cfg.KnowledgeIndexing.IngestFetchAllowedContentTypePrefixes,
	}
	proc := knowledgeindexing.NewProcessor(redisClient, ingestJobRepo, spaceRepo, fetchOpts, objStorage, pipeline, parsers)
	return mq.Subscribe(ctx, "knowledge-ingest", func(payload []byte) {
		pctx, pcancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer pcancel()
		if procErr := proc.Process(pctx, payload); procErr != nil {
			slog.Error("knowledge: embedded ingest failed", "error", procErr)
		}
	})
}
