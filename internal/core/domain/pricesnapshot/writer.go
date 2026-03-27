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
