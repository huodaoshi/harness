package application

import (
	"context"
	"log/slog"
	"time"

	"github.com/huodaoshi/harness/backend/modules/knowledge/domain"
)

const (
	outboxFlushInterval = 15 * time.Second
	outboxBatchLimit    = int64(32)
	maxOutboxAttempts   = 36
	outboxOpTimeout     = 2 * time.Minute
)

// StartMQOutboxFlusher periodically retries MQ publish for rows in mq_outbox (lazy outbox).
// Cancel ctx to stop (call before closing the MessageQueue on shutdown).
func StartMQOutboxFlusher(ctx context.Context, repo domain.MQOutboxRepo, jobs domain.IngestJobRepo, mq domain.MQPublisher) {
	t := time.NewTicker(outboxFlushInterval)
	defer t.Stop()

	flush := func() {
		opCtx, cancel := context.WithTimeout(ctx, outboxOpTimeout)
		flushOutboxBatch(opCtx, repo, jobs, mq)
		cancel()
	}
	flush()

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			flush()
		}
	}
}

func flushOutboxBatch(ctx context.Context, repo domain.MQOutboxRepo, jobs domain.IngestJobRepo, mq domain.MQPublisher) {
	rows, err := repo.ListPending(ctx, outboxBatchLimit)
	if err != nil {
		slog.WarnContext(ctx, "application: outbox: list", "error", err)
		return
	}
	for _, row := range rows {
		if row == nil || row.OutboxID == "" {
			continue
		}
		if err := mq.Publish(ctx, row.Topic, row.Payload); err != nil {
			n, incErr := repo.IncrementFailure(ctx, row.OutboxID, err.Error())
			if incErr != nil {
				slog.WarnContext(ctx, "application: outbox: increment", "outbox_id", row.OutboxID, "error", incErr)
				continue
			}
			if n >= maxOutboxAttempts {
				msg := "mq publish exhausted after background retries"
				_ = jobs.UpdateStatus(ctx, row.JobID, 4, truncateIngestJobErrMsg(msg+": "+err.Error()))
				if delErr := repo.Delete(ctx, row.OutboxID); delErr != nil {
					slog.WarnContext(ctx, "application: outbox: delete exhausted", "outbox_id", row.OutboxID, "error", delErr)
				}
			}
			continue
		}
		if delErr := repo.Delete(ctx, row.OutboxID); delErr != nil {
			slog.WarnContext(ctx, "application: outbox: delete after publish", "outbox_id", row.OutboxID, "error", delErr)
		}
	}
}
