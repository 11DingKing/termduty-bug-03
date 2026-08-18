package shard_test

import (
	"context"
	"testing"

	"termduty/internal/store/shard"
)

func TestRollbackCleansUpWhenRequestContextIsCanceled(t *testing.T) {
	s := newShardStore(t)
	sid := shard.ShardID("readings", "rollback-cancel")
	lease, err := s.Begin(context.Background(), sid)
	if err != nil { t.Fatalf("begin: %v", err) }
	if err := lease.AppendLine(context.Background(), []byte(`{"queue_count":7}`)); err != nil {
		t.Fatalf("append: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := lease.Rollback(ctx); err != nil {
		t.Fatalf("rollback should clean up despite cancellation: %v", err)
	}
	stat, err := s.Recompute(context.Background(), sid)
	if err != nil { t.Fatalf("recompute: %v", err) }
	if stat.Count != 0 { t.Fatalf("count=%d want 0 after rollback", stat.Count) }
}
