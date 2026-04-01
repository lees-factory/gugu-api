package pricesnapshot

import "context"

type SKUSnapshotWriter interface {
	Upsert(ctx context.Context, snapshot SKUPriceSnapshot) error
}

type skuSnapshotWriter struct {
	repository SKUSnapshotRepository
}

func NewSKUSnapshotWriter(repository SKUSnapshotRepository) SKUSnapshotWriter {
	return &skuSnapshotWriter{repository: repository}
}

func (w *skuSnapshotWriter) Upsert(ctx context.Context, snapshot SKUPriceSnapshot) error {
	return w.repository.Upsert(ctx, snapshot)
}
