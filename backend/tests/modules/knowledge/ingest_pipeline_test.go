package knowledge_test

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	redisCli "github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/huodaoshi/harness/backend/conf"
	"github.com/huodaoshi/harness/backend/infra"
	knowledgeapp "github.com/huodaoshi/harness/backend/modules/knowledge/application"
	"github.com/huodaoshi/harness/backend/modules/knowledge/domain"
	knowledgeinf "github.com/huodaoshi/harness/backend/modules/knowledge/infra"
	"github.com/huodaoshi/harness/backend/modules/knowledgeindexing"
	"github.com/huodaoshi/harness/backend/pkg/idgen"
)

func TestIngestInlineMarkdownE2E(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer mr.Close()

	ctx := context.Background()
	redisClient := redisCli.NewClient(&redisCli.Options{Addr: mr.Addr(), Protocol: 2})

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		t.Skip("mongodb not available:", err)
	}
	t.Cleanup(func() { _ = mongoClient.Disconnect(context.Background()) })
	if err := mongoClient.Ping(ctx, nil); err != nil {
		t.Skip("mongodb not available:", err)
	}

	db := mongoClient.Database("harness_knowledge_ingest_test")
	_ = db.Drop(ctx)

	if err := knowledgeinf.EnsureKnowledgeIndexes(ctx, db); err != nil {
		t.Fatal(err)
	}
	if _, err := knowledgeinf.EnsureDefaultSpace(ctx, db); err != nil {
		t.Fatal(err)
	}
	// RediSearch (FT.*) is required for vector index/search; miniredis does not implement it.
	supportsFT := redisClient.Do(ctx, "FT._LIST").Err() == nil
	if supportsFT {
		if err := knowledgeinf.EnsureKnowledgeVectorIndex(ctx, redisClient, knowledgeinf.DefaultSpaceID, infra.FakeEmbedDim()); err != nil {
			t.Fatal(err)
		}
	}

	idg, err := idgen.NewIDGenerator(9)
	if err != nil {
		t.Fatal(err)
	}
	mq, err := infra.NewMessageQueue(conf.MQConfig{Provider: "local"})
	if err != nil {
		t.Fatal(err)
	}

	spaceRepo := knowledgeinf.NewMongoSpaceRepo(db)
	ingestJobRepo := knowledgeinf.NewMongoIngestJobRepo(db, idg)
	mqOutbox := knowledgeinf.NewMongoMQOutboxRepo(db, idg)
	chunkRepo := knowledgeinf.NewRedisIngestChunkRepo(redisClient)
	objStorage, err := infra.NewObjectStorage(conf.COSConfig{Provider: "local"})
	if err != nil {
		t.Fatal(err)
	}

	ingestSvc := knowledgeapp.NewIngestService(spaceRepo, ingestJobRepo, mq, mqOutbox, chunkRepo, objStorage)
	knowledgeSvc := knowledgeapp.NewKnowledgeService(spaceRepo, knowledgeinf.NewRedisVectorRepo(redisClient))

	indexerRedis := redisCli.NewClient(&redisCli.Options{Addr: mr.Addr(), Protocol: 2})
	embedder := infra.NewFakeEmbedder()
	indexer, err := knowledgeindexing.BuildRedisIndexer(ctx, indexerRedis, embedder)
	if err != nil {
		t.Fatal(err)
	}
	pipeline, err := knowledgeindexing.BuildIndexingPipeline(ctx, indexer)
	if err != nil {
		t.Fatal(err)
	}
	parsers, err := knowledgeindexing.NewParsers(ctx)
	if err != nil {
		t.Fatal(err)
	}
	proc := knowledgeindexing.NewProcessor(redisClient, ingestJobRepo, spaceRepo, knowledgeindexing.FetchOptions{}, objStorage, pipeline, parsers)

	workerCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	if err := mq.Subscribe(workerCtx, "knowledge-ingest", func(payload []byte) {
		pctx, pcancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer pcancel()
		if procErr := proc.Process(pctx, payload); procErr != nil {
			t.Error(procErr)
		}
	}); err != nil {
		t.Fatal(err)
	}

	content := "# FAQ\n\n洪峰陪伴时请先确认用户安全。"
	job, err := ingestSvc.SubmitIngestJob(ctx, knowledgeinf.DefaultSpaceID, 1, "", content, 2, "")
	if err != nil {
		t.Fatal(err)
	}

	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		got, err := ingestSvc.GetIngestJob(ctx, job.JobID)
		if err != nil {
			t.Fatal(err)
		}
		if got.Status == domain.IngestStatusDone {
			if got.ChunkCount < 1 {
				t.Fatalf("expected chunk_count >= 1, got %d", got.ChunkCount)
			}
			if supportsFT {
				chunks, err := knowledgeSvc.RetrieveChunks(ctx, ptrInt64(knowledgeinf.DefaultSpaceID), "洪峰陪伴", 3, infra.FakeEmbedVector("洪峰陪伴"))
				if err != nil {
					t.Fatal(err)
				}
				if len(chunks) == 0 {
					t.Fatal("expected FT.SEARCH hits, got 0")
				}
			}
			return
		}
		if got.Status == domain.IngestStatusFailed {
			t.Fatalf("ingest failed: %s", got.Error)
		}
		time.Sleep(200 * time.Millisecond)
	}
	t.Fatal("timeout waiting for ingest job done")
}

func ptrInt64(v int64) *int64 { return &v }
