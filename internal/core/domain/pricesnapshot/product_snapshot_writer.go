package pricesnapshot

import "context"

type ProductSnapshotWriter interface {
	Upsert(ctx context.Context, snapshot ProductPriceSnapshot) error
}

type productSnapshotWriter struct {
	repository ProductSnapshotRepository
}

func NewProductSnapshotWriter(repository ProductSnapshotRepository) ProductSnapshotWriter {
	return &productSnapshotWriter{repository: repository}
}

func (w *productSnapshotWriter) Upsert(ctx context.Context, snapshot ProductPriceSnapshot) error {
	return w.repository.Upsert(ctx, snapshot)
}
